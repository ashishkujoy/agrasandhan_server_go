package controllers

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type dbTemplate struct {
	dbName string
	db     *mongo.Database
	client *mongo.Client
	locker sync.Locker
}

func newDBTemplate(dbName string) *dbTemplate {
	uniqueDbName := fmt.Sprintf("%s%d", dbName, time.Now().UnixMilli())
	client := createDbClient()
	db := client.Database(uniqueDbName)
	var locker sync.Mutex
	return &dbTemplate{
		locker: &locker,
		db:     db,
		client: client,
		dbName: uniqueDbName,
	}
}

func (template *dbTemplate) close() {
	_ = template.client.Disconnect(context.Background())
}

func createDbClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil
	}
	return client
}
