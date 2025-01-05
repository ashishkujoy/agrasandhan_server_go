package controllers

import (
	"ashishkujoy/agrasandhan/requests"
	"ashishkujoy/agrasandhan/services"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type batchCreationRequest struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"startDate"`
}

// CreateBatch create a new batch
func CreateBatch(batchService *services.BatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody batchCreationRequest
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		batch, err := batchService.CreateBatch(reqBody.Name, reqBody.StartDate)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, batch)
	}
}

// GetAllBatches retrieves all the batches.
func GetAllBatches(service *services.BatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		batches, err := service.GetAllBatches()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, batches)
	}
}

// GetBatchById retrieves a batch by its id.
func GetBatchById(service *services.BatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		batch, err := service.GetBatchById(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, batch)
	}
}

// AssignMentor assigns a mentor to a batch.
func AssignMentor(service *services.BatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		var mentorReq requests.AssignMentorRequest
		err = c.ShouldBindJSON(&mentorReq)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		mentor, err := service.AssignMentor(id, mentorReq)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, mentor)
	}
}
