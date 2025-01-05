package services

import (
	"ashishkujoy/agrasandhan/repositories"
	"ashishkujoy/agrasandhan/repositories/models"
	"time"
)

type BatchService struct {
	batchRepository repositories.BatchRepository
	idGenerator     IdGenerator
}

// NewBatchService creates a new instance of BatchService.
func NewBatchService(batchRepository repositories.BatchRepository, idGenerator IdGenerator) *BatchService {
	return &BatchService{batchRepository: batchRepository, idGenerator: idGenerator}
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

func (s *BatchService) GetBatchById(id int) (*models.Batch, error) {
	return s.batchRepository.FindById(id)
}
