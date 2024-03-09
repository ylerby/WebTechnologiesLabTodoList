package cache

import (
	"context"
	"fmt"
	"time"

	"backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	valueForUpdateNumber = iota
	filterValueNumber
	bsonIdFilterKey          = "id"
	bsonAuthorIdFilterKey    = "author_id"
	bsonTitleFilterKey       = "title"
	bsonDescriptionFilterKey = "description"
	bsonCommentFilterKey     = "comments"
	contextTimeout           = 10
)

func (c *Cache) SetValue(value domain.TodoListModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*contextTimeout)
	defer cancel()

	_, err := c.collection.InsertOne(ctx, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetValueByTitle(title string) ([]domain.TodoListModel, error) {
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

	var result []domain.TodoListModel
	for cursor.Next(ctx) {
		var todoList domain.TodoListModel
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

func (c *Cache) GetAllValues() ([]domain.TodoListModel, error) {
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

	var result []domain.TodoListModel
	for cursor.Next(ctx) {
		var todoList domain.TodoListModel
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

func (c *Cache) UpdateValue(values domain.UpdateTodoList) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*contextTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			bsonIdFilterKey:          values[valueForUpdateNumber].Id,
			bsonAuthorIdFilterKey:    values[valueForUpdateNumber].AuthorName,
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

func (c *Cache) DeleteValue(value domain.TodoListModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*contextTimeout)
	defer cancel()

	_, err := c.collection.DeleteOne(ctx, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) SetComment(value domain.TodoListComment) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*contextTimeout)
	defer cancel()

	filter := bson.M{
		"id":          value.Id,
		"author_name": value.AuthorName,
		"title":       value.Title,
		"description": value.Description,
	}

	update := bson.M{
		"$push": bson.M{
			bsonCommentFilterKey: value.Comment,
		},
	}

	_, err := c.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
