package routes

import (
	"ashishkujoy/agrasandhan/controllers"
	"ashishkujoy/agrasandhan/services"
	"github.com/gin-gonic/gin"
)

func addUserRoutes(userService *services.UserService, router *gin.Engine) *gin.Engine {
	r := router.Group("/users")
	r.POST("", controllers.AddUser(userService))
	return router
}
