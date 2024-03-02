package cache

import (
	"backend/internal/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

const (
	contextTimeout = 10
)

func (c *Cache) SetValue(value model.TodoListModel) error {
	insertValue, err := c.collection.InsertOne(context.Background(), value)
	if err != nil {
		return err
	}
	log.Printf("inserted value - %v = %v", insertValue, insertValue.InsertedID)
	return nil
}

func (c *Cache) GetAllValues() ([]model.TodoListModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*contextTimeout)
	defer cancel()

	cursor, err := c.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var result []model.TodoListModel
	for cursor.Next(ctx) {
		var todoList model.TodoListModel
		if err = cursor.Decode(&todoList); err != nil {
			return nil, err
		}
		result = append(result, todoList)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
