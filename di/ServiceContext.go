package di

import "ashishkujoy/agrasandhan/services"

type ServiceContext struct {
	BatchService *services.BatchService
	UserService  *services.UserService
}

func NewServiceContext(repositoryCtx *RepositoryContext) *ServiceContext {
	userService := services.NewUserService(
		repositoryCtx.UserRepository,
		services.NewIdGeneratorImpl("user", repositoryCtx.Counters),
	)
	batchService := services.NewBatchService(
		repositoryCtx.BatchRepository,
		services.NewIdGeneratorImpl("batches", repositoryCtx.Counters),
		userService,
	)

	return &ServiceContext{
		BatchService: batchService,
		UserService:  userService,
	}
}
