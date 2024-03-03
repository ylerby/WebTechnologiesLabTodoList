package cache

import (
	"backend/internal/model"
	"backend/internal/schemas"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	todoListCollection = "todoListCollection"
	todoListDatabase   = "todoListDatabase"
)

type Repository interface {
	SetValue(value model.TodoListModel) error
	GetValueByTitle(title string) ([]model.TodoListModel, error)
	GetAllValues() ([]model.TodoListModel, error)
	UpdateValue(values schemas.UpdateTodoList) error
	DeleteValue(value model.TodoListModel) error
}

type Cache struct {
	collection *mongo.Collection
}

func New(mongoURI string) (*Cache, error) {
	ctx := context.Background()

	mongoOptions := options.Client().ApplyURI(mongoURI)

	mongoClient, err := mongo.Connect(ctx, mongoOptions)
	if err != nil {
		return nil, err
	}

	collection := mongoClient.Database(todoListDatabase).Collection(todoListCollection)
	return &Cache{
		collection: collection,
	}, nil
}
