package routes

import (
	"ashishkujoy/agrasandhan/repositories/models"
	"cmp"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var router = addUserRoutes(ServiceContext.UserService, gin.New())

func TestAddUserRoutes(t *testing.T) {
	_ = RepositoryContext.UserRepository.DeleteAll()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/users",
		strings.NewReader(`{"name": "Ashish", "email": "akj@test.com", "roles": ["admin"]}`),
	)

	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var actualBody models.User
	err := json.Unmarshal(w.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Equal(t, "Ashish", actualBody.Name)
	assert.Equal(t, "akj@test.com", actualBody.Email)
	assert.Equal(t, []string{"admin"}, actualBody.Roles)
	assert.NotNil(t, actualBody.ID)
}

func TestGetAllUsersRoutes(t *testing.T) {
	repository := RepositoryContext.UserRepository
	_ = repository.DeleteAll()

	_ = repository.Save(&models.User{Name: "Martha", Email: "", Roles: []string{"admin"}})
	_ = repository.Save(&models.User{Name: "James", Email: "", Roles: []string{"admin"}})
	_ = repository.Save(&models.User{Name: "Jordan", Email: "", Roles: []string{"admin"}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var actualBody []*models.User
	err := json.Unmarshal(w.Body.Bytes(), &actualBody)
	assert.NoError(t, err)

	// sort the users by name
	actualBody = Sort(actualBody, func(user *models.User) string { return user.Name })

	assert.Len(t, actualBody, 3)
	assert.Equal(t, "James", actualBody[0].Name)
	assert.Equal(t, "Jordan", actualBody[1].Name)
	assert.Equal(t, "Martha", actualBody[2].Name)
}

// Sort sorts the items using the given function.
func Sort[T any, U cmp.Ordered](items []*T, f func(*T) U) []*T {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if f(items[i]) > f(items[j]) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	return items
}
