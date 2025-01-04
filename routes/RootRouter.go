package routes

import (
	"ashishkujoy/agrasandhan/di"
	"github.com/gin-gonic/gin"
)

func NewRootRouter(serviceCtx *di.ServiceContext) *gin.Engine {
	r := gin.Default()
	addUserRoutes(serviceCtx.UserService, r)

	return r
}
