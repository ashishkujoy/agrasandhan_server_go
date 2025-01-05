package repositories

import (
	"ashishkujoy/agrasandhan/repositories/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

// BatchRepository interface defines the methods to save, update, retrieve and list batches.
type BatchRepository interface {
	Save(batch *models.Batch) error
	DeleteAll() error
	FindById(id int) (*models.Batch, error)
	GetAll() ([]*models.Batch, error)
	Update(batch *models.Batch) error
}

// BatchRepositoryImpl is the implementation of BatchRepository interface. It uses mongoDb underneath to store the data.
type BatchRepositoryImpl struct {
	collection *mongo.Collection
	mutex      *sync.RWMutex
}

// NewBatchRepository creates a new instance of BatchRepositoryImpl.
func NewBatchRepository(collection *mongo.Collection) *BatchRepositoryImpl {
	var mutex sync.RWMutex
	return &BatchRepositoryImpl{collection: collection, mutex: &mutex}
}

// Save method saves the batch to the database. It applies a timeout of 2 seconds.
func (r *BatchRepositoryImpl) Save(batch *models.Batch) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, batch)
	return err
}

// DeleteAll method deletes all the batches from the database. It applies a timeout of 2 seconds. It is used for testing purposes.
func (r *BatchRepositoryImpl) DeleteAll() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := r.collection.DeleteMany(ctx, bson.D{})
	return err
}

// FindById method retrieves the batch from the database by its id. It applies a timeout of 5 seconds.
func (r *BatchRepositoryImpl) FindById(id int) (*models.Batch, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var batch models.Batch
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&batch)
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetAll method retrieves all the batches from the database. It applies a timeout of 5 seconds.
func (r *BatchRepositoryImpl) GetAll() ([]*models.Batch, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	var batches []*models.Batch
	for cursor.Next(ctx) {
		var batch models.Batch
		err := cursor.Decode(&batch)
		if err != nil {
			return nil, err
		}
		batches = append(batches, &batch)
	}
	return batches, nil
}

// Update method updates the batch in the database. It applies a timeout of 2 seconds.
func (r *BatchRepositoryImpl) Update(batch *models.Batch) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"id": batch.ID}, batch)
	return err
}
