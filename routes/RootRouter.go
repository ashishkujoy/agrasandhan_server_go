package routes

import (
	"ashishkujoy/agrasandhan/configs"
	"ashishkujoy/agrasandhan/controllers"
	"ashishkujoy/agrasandhan/di"

	"github.com/gin-gonic/gin"
)

func NewRootRouter(
	serviceCtx *di.ServiceContext,
	session gin.HandlerFunc,
	env *configs.Env,
) *gin.Engine {
	r := gin.Default()
	r.Use(session)
	r.Use(controllers.NewAuth(serviceCtx.UserService, env.JWTKey))
	r.GET("", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
	addUserRoutes(serviceCtx.UserService, r)
	addBatchRoutes(serviceCtx.BatchService, r)
	return r
}
