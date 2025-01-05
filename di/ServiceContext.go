package di

import "ashishkujoy/agrasandhan/services"

type ServiceContext struct {
	BatchService *services.BatchService
	UserService  *services.UserService
}

func NewServiceContext(repositoryCtx *RepositoryContext) *ServiceContext {
	return &ServiceContext{
		BatchService: services.NewBatchService(
			repositoryCtx.BatchRepository,
			services.NewIdGeneratorImpl("batches", repositoryCtx.Counters),
		),
		UserService: services.NewUserService(
			repositoryCtx.UserRepository,
			services.NewIdGeneratorImpl("user", repositoryCtx.Counters),
		),
	}
}
