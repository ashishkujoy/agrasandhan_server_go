package di

import (
	"ashishkujoy/agrasandhan/configs"
	"ashishkujoy/agrasandhan/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryContext struct {
	UserRepository repositories.UserRepository
	Counters       *mongo.Collection
}

func NewRepositoryContext(env *configs.Env) *RepositoryContext {
	dbClient := configs.ConnectDB(*env)
	db := dbClient.Database(env.DBName)

	userRepository := repositories.NewUserRepository(db.Collection("users"))

	return &RepositoryContext{
		UserRepository: userRepository,
		Counters:       db.Collection("counters"),
	}
}
