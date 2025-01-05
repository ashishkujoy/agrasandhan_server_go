package services

import (
	"ashishkujoy/agrasandhan/repositories"
	"ashishkujoy/agrasandhan/repositories/models"
	"ashishkujoy/agrasandhan/requests"
	"time"
)

type BatchService struct {
	batchRepository repositories.BatchRepository
	userService     *UserService
	idGenerator     IdGenerator
}

// NewBatchService creates a new instance of BatchService.
func NewBatchService(
	batchRepository repositories.BatchRepository,
	idGenerator IdGenerator,
	userService *UserService,
) *BatchService {
	return &BatchService{
		batchRepository: batchRepository,
		idGenerator:     idGenerator,
		userService:     userService,
	}
}

// CreateBatch creates a new batch with the given start date.
func (s *BatchService) CreateBatch(name string, startDate time.Time) (*models.Batch, error) {
	id := s.idGenerator.GenerateNum()
	batch := &models.Batch{
		ID:        id,
		Name:      name,
		StartDate: startDate,
		Interns:   []string{},
		Mentors:   []models.Mentor{},
	}
	err := s.batchRepository.Save(batch)
	if err != nil {
		return nil, err
	}

	return batch, nil
}

// GetAllBatches retrieves all the batches from the database.
func (s *BatchService) GetAllBatches() ([]*models.Batch, error) {
	return s.batchRepository.GetAll()
}

// AssignMentor assigns a mentor to a batch.
func (s *BatchService) AssignMentor(id int, req requests.AssignMentorRequest) (*models.Mentor, error) {
	batch, err := s.batchRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	_, err = s.userService.GetUserById(req.Id)
	if err != nil {
		return nil, err
	}
	mentor := models.Mentor{
		ID:          req.Id,
		Permissions: req.Permissions,
	}
	batch.Mentors = append(batch.Mentors, mentor)
	err = s.batchRepository.Update(batch)
	if err != nil {
		return nil, err
	}
	return &mentor, nil
}

func (s *BatchService) GetBatchById(id int) (*models.Batch, error) {
	return s.batchRepository.FindById(id)
}
