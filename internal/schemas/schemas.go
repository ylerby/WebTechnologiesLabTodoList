package schemas

import "backend/internal/model"

const (
	updateFieldsNumber = 2
)

type CorrectResponse struct {
	Data interface{} `json:"data"`
}

type GetTodoListByTitle struct {
	Title string `json:"title" bson:"title"`
}

type UpdateTodoList [updateFieldsNumber]model.TodoListModel

type RegisterUserRequestBody struct {
	Login         string `json:"login"`
	Password      string `json:"password"`
	AgainPassword string `json:"again_password"`
}
