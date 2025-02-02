package middlewares

import (
	"ashishkujoy/agrasandhan/configs"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func NewSession(env *configs.Env) gin.HandlerFunc {
	store, err := redis.NewStore(10, "tcp", env.RedisURI, "", []byte(env.SessionSecret))
	if err != nil {
		panic(fmt.Sprintf("Failed to create redis session store: %v", err))
	}

	return sessions.Sessions("sessions", store)
}
