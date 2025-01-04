package services

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

type TaggedId struct {
	Tag   string `bson:"tag"`
	Value int    `bson:"value"`
}

type IdGenerator interface {
	GenerateStr() string
	GenerateNum() int
}

type IdGeneratorImpl struct {
	tag        string
	collection *mongo.Collection
	locker     *sync.Mutex
}

func NewIdGeneratorImpl(tag string, collection *mongo.Collection) *IdGeneratorImpl {
	var mutex sync.Mutex
	return &IdGeneratorImpl{tag: tag, collection: collection, locker: &mutex}
}

func (g *IdGeneratorImpl) GenerateStr() string {
	return fmt.Sprintf("%d", g.GenerateNum())
}

func (g *IdGeneratorImpl) GenerateNum() int {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	g.locker.Lock()
	defer g.locker.Unlock()
	var taggedId TaggedId
	searchQuery := bson.D{{"tag", g.tag}}
	update := bson.D{{"$inc", bson.D{{"value", 1}}}}
	g.collection.FindOneAndUpdate(ctx, searchQuery, update)
	err := g.collection.FindOne(ctx, searchQuery).Decode(&taggedId)
	if err != nil {
		_, _ = g.collection.InsertOne(ctx, TaggedId{Tag: g.tag, Value: 1})
	}
	return taggedId.Value
}
