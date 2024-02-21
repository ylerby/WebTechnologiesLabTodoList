package cache

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	SetValue(key, value string) error
	GetValue(key string) (string, error)
	UpdateValue(key, newValue string) error
	DeleteValue(key string) error
}

type Cache struct {
	client *mongo.Client
}

func New(mongoURI string) (*Cache, error) {
	ctx := context.Background()

	mongoOptions := options.Client().ApplyURI(mongoURI)

	mongoClient, err := mongo.Connect(ctx, mongoOptions)
	if err != nil {
		return nil, err
	}

	return &Cache{
		client: mongoClient,
	}, nil
}
