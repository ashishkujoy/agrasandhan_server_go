package routes

import (
	"ashishkujoy/agrasandhan/repositories/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var r = addBatchRoutes(ServiceContext.BatchService, gin.New())

func TestBatchRoutes(t *testing.T) {
	t.Run("Create Batch", func(t *testing.T) {
		_ = RepositoryContext.BatchRepository.DeleteAll()
		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			"POST",
			"/batches",
			strings.NewReader(`{"name": "Batch 10", "startDate": "2025-01-03T15:04:05Z"}`),
		)
		r.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		var batch models.Batch
		err := json.Unmarshal(w.Body.Bytes(), &batch)
		assert.NoError(t, err)
		assert.Equal(t, "Batch 10", batch.Name)
		ti, err := time.Parse("2006-01-02T15:04:05Z", "2025-01-03T15:04:05Z")
		assert.Equal(t, ti, batch.StartDate)
	})

	t.Run("Fetch All Batch", func(t *testing.T) {
		repository := RepositoryContext.BatchRepository
		_ = repository.DeleteAll()

		_ = repository.Save(&models.Batch{Name: "Batch 1", StartDate: time.Now()})
		_ = repository.Save(&models.Batch{Name: "Batch 2", StartDate: time.Now()})
		_ = repository.Save(&models.Batch{Name: "Batch 3", StartDate: time.Now()})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/batches", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		var batches []*models.Batch
		err := json.Unmarshal(w.Body.Bytes(), &batches)
		assert.NoError(t, err)

		// sort the batches by name
		batches = Sort(batches, func(t *models.Batch) string {
			return t.Name
		})

		assert.Len(t, batches, 3)
		assert.Equal(t, "Batch 1", batches[0].Name)
		assert.Equal(t, "Batch 2", batches[1].Name)
		assert.Equal(t, "Batch 3", batches[2].Name)
	})

	t.Run("Fetch Batch By Id", func(t *testing.T) {
		repository := RepositoryContext.BatchRepository
		batch12 := &models.Batch{Name: "Batch 2", StartDate: time.Now().UTC(), ID: 12}
		_ = repository.DeleteAll()

		_ = repository.Save(&models.Batch{Name: "Batch 1", StartDate: time.Now()})
		_ = repository.Save(batch12)
		_ = repository.Save(&models.Batch{Name: "Batch 3", StartDate: time.Now()})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/batches/12", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		var batch *models.Batch
		err := json.Unmarshal(w.Body.Bytes(), &batch)
		assert.NoError(t, err)

		assert.Equal(t, batch12.Name, batch.Name)
	})

	t.Run("Assign Mentor to Batch", func(t *testing.T) {
		repository := RepositoryContext.BatchRepository
		_ = repository.DeleteAll()
		_ = RepositoryContext.UserRepository.Save(&models.User{
			ID:    "11",
			Name:  "Mentor 1",
			Email: "",
			Roles: []string{"mentor"},
		})

		batch := &models.Batch{ID: 10, Name: "Batch 1", StartDate: time.Now()}
		_ = repository.Save(batch)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			"POST",
			"/batches/10/mentors",
			strings.NewReader(`{
				"id": "11", 
				"permissions": {
					"allowProvideObservations": true, 
					"allowReleaseIntern": false,
					"allowProvideFeedback": false,
					"allowDeliverFeedback": true
				}
			}`))

		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		var mentor models.Mentor
		err := json.Unmarshal(w.Body.Bytes(), &mentor)
		assert.NoError(t, err)
		assert.Equal(t, models.Mentor{
			ID: "11",
			Permissions: models.MentorPermission{
				AllowProvideObservations: true,
				AllowReleaseIntern:       false,
				AllowProvideFeedback:     false,
				AllowDeliverFeedback:     true,
			},
		}, mentor)
	})
}
