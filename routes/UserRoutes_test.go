package routes

import (
	"ashishkujoy/agrasandhan/configs"
	"ashishkujoy/agrasandhan/di"
	"ashishkujoy/agrasandhan/repositories/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var config = &configs.Env{
	MongoURI: "mongodb://localhost:27017",
	DBName:   fmt.Sprintf("test-%d", time.Now().UnixMilli()),
}
var repositoryContext = di.NewRepositoryContext(config)
var serviceContext = di.NewServiceContext(repositoryContext)
var router = addUserRoutes(serviceContext.UserService, gin.New())

func TestMain(m *testing.M) {
	exitCode := m.Run()
	dbClient := configs.ConnectDB(*config)
	_ = dbClient.Database(config.DBName).Drop(context.Background())

	os.Exit(exitCode)
}

func TestAddUserRoutes(t *testing.T) {
	_ = repositoryContext.UserRepository.DeleteAll()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/users",
		strings.NewReader(`{"name": "Ashish", "email": "akj@test.com", "role": 1}`),
	)

	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var actualBody models.User
	err := json.Unmarshal(w.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Equal(t, "Ashish", actualBody.Name)
	assert.Equal(t, "akj@test.com", actualBody.Email)
	assert.Equal(t, models.UserRole(1), actualBody.Role)
	assert.NotNil(t, actualBody.ID)
}

func TestGetAllUsersRoutes(t *testing.T) {
	repository := repositoryContext.UserRepository
	_ = repository.DeleteAll()

	_ = repository.Save(&models.User{Name: "Martha", Email: "", Role: models.UserRole(1)})
	_ = repository.Save(&models.User{Name: "James", Email: "", Role: models.UserRole(1)})
	_ = repository.Save(&models.User{Name: "Jordan", Email: "", Role: models.UserRole(1)})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var actualBody []models.User
	err := json.Unmarshal(w.Body.Bytes(), &actualBody)
	assert.NoError(t, err)

	// sort the users by name
	actualBody = sortUsers(actualBody)

	assert.Len(t, actualBody, 3)
	assert.Equal(t, "James", actualBody[0].Name)
	assert.Equal(t, "Jordan", actualBody[1].Name)
	assert.Equal(t, "Martha", actualBody[2].Name)
}

func sortUsers(users []models.User) []models.User {
	for i := 0; i < len(users); i++ {
		for j := i + 1; j < len(users); j++ {
			if users[i].Name > users[j].Name {
				users[i], users[j] = users[j], users[i]
			}
		}
	}
	return users
}
