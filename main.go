package main

import (
	"ashishkujoy/agrasandhan/configs"
	"ashishkujoy/agrasandhan/di"
	"ashishkujoy/agrasandhan/middlewares"
	"ashishkujoy/agrasandhan/routes"
	"fmt"
)

func main() {
	env := configs.NewEnv()
	repositoriesCtx := di.NewRepositoryContext(env)
	serviceCtx := di.NewServiceContext(repositoriesCtx)

	r := routes.NewRootRouter(serviceCtx, middlewares.NewMongoSession(env), env)

	err := r.Run(fmt.Sprintf(":%s", env.Port))

	if err != nil {
		panic(fmt.Errorf("error starting server: %v", err))
	}
}
