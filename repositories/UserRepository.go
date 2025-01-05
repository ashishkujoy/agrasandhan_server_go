package repositories

import (
	"ashishkujoy/agrasandhan/repositories/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// UserRepository interface defines the methods to save, update, retrieve and list users.
type UserRepository interface {
	Save(user *models.User) error
	DeleteAll() error
	FindById(id string) (*models.User, error)
	GetAll() ([]*models.User, error)
	UpdateRole(id string, role models.UserRole) error
}

// UserRepositoryImpl is the implementation of UserRepository interface. It uses mongoDb underneath to store the data.
type UserRepositoryImpl struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new instance of UserRepositoryImpl.
func NewUserRepository(collection *mongo.Collection) *UserRepositoryImpl {
	return &UserRepositoryImpl{collection: collection}
}

// Save method saves the user to the database. It applies a timeout of 2 seconds.
func (r *UserRepositoryImpl) Save(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// FindById method retrieves the user from the database by its id. It applies a timeout of 5 seconds.
func (r *UserRepositoryImpl) FindById(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll method retrieves all the users from the database. It applies a timeout of 5 seconds.
func (r *UserRepositoryImpl) GetAll() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			fmt.Printf("Error closing cursor: %v\n", err)
		}
	}(cursor, ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// UpdateRole method updates the role of the user in the database. It applies a timeout of 5 seconds.
func (r *UserRepositoryImpl) UpdateRole(id string, role models.UserRole) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"role": role}})
	return err
}

// DeleteAll method deletes all the users from the database. It applies a timeout of 2 seconds. It is used for testing purposes.
func (r *UserRepositoryImpl) DeleteAll() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := r.collection.DeleteMany(ctx, bson.M{})
	return err
}
