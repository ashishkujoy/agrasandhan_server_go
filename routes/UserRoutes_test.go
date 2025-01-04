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
