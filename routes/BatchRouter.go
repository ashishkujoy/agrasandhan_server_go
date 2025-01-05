package routes

import (
	"ashishkujoy/agrasandhan/controllers"
	"ashishkujoy/agrasandhan/services"
	"github.com/gin-gonic/gin"
)

func addBatchRoutes(batchService *services.BatchService, router *gin.Engine) *gin.Engine {
	r := router.Group("/batches")
	r.POST("", controllers.CreateBatch(batchService))
	r.GET("", controllers.GetAllBatches(batchService))
	r.GET("/:id", controllers.GetBatchById(batchService))
	return router
}
