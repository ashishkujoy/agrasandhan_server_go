package controllers

import (
	"ashishkujoy/agrasandhan/repositories/models"
	"ashishkujoy/agrasandhan/services"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	codeRedirect = 302
	keyToken     = "oauth2_token"
	keyNextPage  = "next"
	LoginPath    = "/login"
	LogoutPath   = "/logout"
	CallbackPath = "/auth/google/callback"
	ErrorPath    = "/unauthorized"
)

type UserGmailProfile struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

// Tokens represents a container that contains user's OAuth 2.0 access and refresh tokens.
type Tokens interface {
	Access() string
	Refresh() string
	Expired() bool
	ExpiryTime() time.Time
}

type token struct {
	oauth2.Token
}

// Access returns the access token.
func (t *token) Access() string {
	return t.AccessToken
}

// Refresh returns the refresh token.
func (t *token) Refresh() string {
	return t.RefreshToken
}

// Expired returns whether the access token is expired or not.
func (t *token) Expired() bool {
	if t == nil {
		return true
	}
	return !t.Token.Valid()
}

// ExpiryTime returns the expiry time of the user's access token.
func (t *token) ExpiryTime() time.Time {
	return t.Expiry
}

// String returns the string representation of the token.
func (t *token) String() string {
	return fmt.Sprintf("tokens: %s expire at: %s", t.Access(), t.ExpiryTime())
}

// generateToken generates a new JWT token for the given user.
func generateToken(singingKey []byte, user *models.User) ([]byte, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(time.Hour * 24 * 366).Unix(),
		"roles": user.Roles,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key, err := token.SignedString(singingKey)
	if err != nil {
		return nil, err
	}
	return []byte(key), nil
}

// decodeToken decodes the given JWT token and returns the user.
func decodeToken(signingKey []byte, tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	user := &models.User{
		ID:    claims["id"].(string),
		Name:  claims["name"].(string),
		Email: claims["email"].(string),
		Roles: claims["roles"].([]string),
	}
	return user, nil
}

// Google returns a new Google OAuth 2.0 backend endpoint.
func Google(
	conf *oauth2.Config,
	userService *services.UserService,
	signingKey []byte,
) gin.HandlerFunc {
	return NewOAuth2Provider(conf, userService, signingKey)
}

// NewOAuth2Provider returns a generic OAuth 2.0 backend endpoint.
func NewOAuth2Provider(
	conf *oauth2.Config,
	userService *services.UserService,
	signingKey []byte,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			switch c.Request.URL.Path {
			case LoginPath:
				login(conf, c)
			case LogoutPath:
				logout(c)
			case CallbackPath:
				handleOAuth2Callback(conf, c, userService, signingKey)
			}
		}
		s := sessions.Default(c)
		fmt.Printf("Session %v\n", s)
		user, err := unmarshallToken(signingKey, s)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func login(f *oauth2.Config, c *gin.Context) {
	s := sessions.Default(c)
	next := extractPath(c.Request.URL.Query().Get(keyNextPage))
	if s.Get(keyToken) == nil {
		// User is not logged in.
		if next == "" {
			next = "/"
		}
		http.Redirect(c.Writer, c.Request, f.AuthCodeURL(next), codeRedirect)
		return
	}
	// No need to log in, redirect to the next page.
	http.Redirect(c.Writer, c.Request, next, codeRedirect)
}

func logout(c *gin.Context) {
	s := sessions.Default(c)
	next := extractPath(c.Request.URL.Query().Get(keyNextPage))
	s.Delete(keyToken)
	_ = s.Save()
	http.Redirect(c.Writer, c.Request, next, codeRedirect)
}

func extractPath(next string) string {
	n, err := url.Parse(next)
	if err != nil {
		return "/"
	}
	return n.Path
}

func handleOAuth2Callback(
	f *oauth2.Config,
	c *gin.Context,
	userService *services.UserService,
	signingKey []byte,
) {
	s := sessions.Default(c)
	code := c.Request.URL.Query().Get("code")
	t, err := f.Exchange(context.Background(), code)
	if err != nil {
		log.Println("exchange oauth token failed:", err)
		http.Redirect(c.Writer, c.Request, ErrorPath, codeRedirect)
		return
	}
	setupSession(t, c, s, userService, signingKey)
}

func setupSession(
	t *oauth2.Token,
	c *gin.Context,
	s sessions.Session,
	userService *services.UserService,
	signingKey []byte,
) {
	userInfo, err := fetchUserInfo(t.AccessToken)
	if err != nil {
		log.Println("fetch user info failed:", err)
		http.Redirect(c.Writer, c.Request, ErrorPath, codeRedirect)
		return
	}
	user, err := userService.GetUserByEmailId(userInfo.Email)
	if err != nil {
		log.Printf("User %s not register", userInfo.Email)
		http.Redirect(c.Writer, c.Request, ErrorPath, codeRedirect)
		return
	}

	val, _ := generateToken(signingKey, user)
	s.Set(keyToken, val)
	err = s.Save()
	if err != nil {
		fmt.Printf("Error saving session %v\n", err)
	}

	next := extractPath(c.Request.URL.Query().Get("state"))
	http.Redirect(c.Writer, c.Request, next, codeRedirect)
}

func fetchUserInfo(accessToken string) (UserGmailProfile, error) {
	var userInfo UserGmailProfile
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return userInfo, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return userInfo, fmt.Errorf("fetch user info failed: %s", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

func unmarshallToken(signingKey []byte, s sessions.Session) (*models.User, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error in unmarshall token %v", r)
			// Handle the panic and log it if necessary
			//err := fmt.Errorf("panic occurred in unmarshallToken: %v", r)
		}
	}()
	if s.Get(keyToken) == nil {
		fmt.Printf("No value present for %s\n", keyToken)
		return nil, errors.New("token not present")
	}
	data := s.Get(keyToken).([]byte)
	fmt.Printf("Token as bytes %s\n", data)
	return decodeToken(signingKey, string(data))
}

var config = oauth2.Config{
	ClientID:     os.Getenv("AGRASANDHAN_GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("AGRASANDHAN_GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8000/auth/google/callback",
	Endpoint:     google.Endpoint,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
	},
}

func NewAuth(
	userService *services.UserService,
	signingKey []byte,
) gin.HandlerFunc {
	return Google(&config, userService, signingKey)
}
