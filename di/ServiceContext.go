package di

import "ashishkujoy/agrasandhan/services"

type ServiceContext struct {
	UserService *services.UserService
}

func NewServiceContext(repositoryCtx *RepositoryContext) *ServiceContext {
	return &ServiceContext{
		UserService: services.NewUserService(
			repositoryCtx.UserRepository,
			services.NewIdGeneratorImpl("user", repositoryCtx.Counters),
		),
	}
}
