package di

import (
	"ashishkujoy/agrasandhan/configs"
	"ashishkujoy/agrasandhan/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryContext struct {
	UserRepository  repositories.UserRepository
	BatchRepository repositories.BatchRepository
	Counters        *mongo.Collection
}

func NewRepositoryContext(env *configs.Env) *RepositoryContext {
	dbClient := configs.ConnectDB(*env)
	db := dbClient.Database(env.DBName)

	userRepository := repositories.NewUserRepository(db.Collection("users"))
	batchRepository := repositories.NewBatchRepository(db.Collection("batches"))

	return &RepositoryContext{
		UserRepository:  userRepository,
		BatchRepository: batchRepository,
		Counters:        db.Collection("counters"),
	}
}
