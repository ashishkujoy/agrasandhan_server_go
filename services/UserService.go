package services

import (
	"ashishkujoy/agrasandhan/repositories"
	"ashishkujoy/agrasandhan/repositories/models"
)

// UserService is the service that provides methods to deal with user domain.
type UserService struct {
	repository  repositories.UserRepository
	idGenerator IdGenerator
}

// NewUserService creates a new instance of UserService.
func NewUserService(repository repositories.UserRepository, idGenerator IdGenerator) *UserService {
	return &UserService{repository: repository, idGenerator: idGenerator}
}

// CreateUser creates a new user with the given name, email and role.
func (s *UserService) CreateUser(name, email string, role int) (*models.User, error) {
	user := &models.User{
		ID:    s.idGenerator.GenerateStr(),
		Name:  name,
		Email: email,
		Role:  models.UserRole(role),
	}

	err := s.repository.Save(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUsers retrieves all the users from the database.
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.repository.GetAll()
}

func (s *UserService) GetUserById(id string) (*models.User, error) {
	return s.repository.FindById(id)
}
