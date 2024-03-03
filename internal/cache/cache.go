package cache

import (
	"backend/internal/model"
	"backend/internal/schemas"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	valueForUpdateNumber = iota
	filterValueNumber
	bsonIdFilterKey          = "id"
	bsonAuthorIdFilterKey    = "author_id"
	bsonTitleFilterKey       = "title"
	bsonDescriptionFilterKey = "description"
	contextTimeout           = 10
)

func (c *Cache) SetValue(value model.TodoListModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*contextTimeout)
	defer cancel()

	_, err := c.collection.InsertOne(ctx, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetValueByTitle(title string) ([]model.TodoListModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*contextTimeout)
	defer cancel()

	filter := bson.M{bsonTitleFilterKey: title}

	cursor, err := c.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			return
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

	if len(result) == 0 {
		return nil, fmt.Errorf("не найдено ни одного значения")
	}

	return result, nil
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
			return
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

func (c *Cache) UpdateValue(values schemas.UpdateTodoList) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*contextTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			bsonIdFilterKey:          values[valueForUpdateNumber].Id,
			bsonAuthorIdFilterKey:    values[valueForUpdateNumber].AuthorId,
			bsonTitleFilterKey:       values[valueForUpdateNumber].Title,
			bsonDescriptionFilterKey: values[valueForUpdateNumber].Description,
		},
	}

	filter := values[filterValueNumber]

	_, err := c.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) DeleteValue(value model.TodoListModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*contextTimeout)
	defer cancel()

	_, err := c.collection.DeleteOne(ctx, value)
	if err != nil {
		return err
	}
	return nil
}
