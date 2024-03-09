package cache

import (
	"context"

	"backend/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	todoListCollection = "todoListCollection"
	todoListDatabase   = "todoListDatabase"
)

type Repository interface {
	SetValue(value domain.TodoListModel) error
	GetValueByTitle(title string) ([]domain.TodoListModel, error)
	GetAllValues() ([]domain.TodoListModel, error)
	UpdateValue(values domain.UpdateTodoList) error
	DeleteValue(value domain.TodoListModel) error
	SetComment(value domain.TodoListComment) error
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
