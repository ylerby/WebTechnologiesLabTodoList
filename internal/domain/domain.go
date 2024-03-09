package domain

type TodoListModel struct {
	Id          int      `json:"id" bson:"id"`
	AuthorName  string   `json:"author_name" bson:"author_name"`
	Title       string   `json:"title" bson:"title"`
	Description string   `json:"description" bson:"description"`
	Comments    []string `json:"comments" bson:"comments"`
}

type UserModel struct {
	Login    string `gorm:"primaryKey;column:login;type:varchar(255)"`
	Password string `gorm:"column:password;type:varchar(255)"`
}

const (
	updateFieldsNumber = 2
)

type CorrectResponse struct {
	Data interface{} `json:"data"`
}

type GetTodoListByTitle struct {
	Title string `json:"title" bson:"title"`
}

type UpdateTodoList [updateFieldsNumber]TodoListModel

type RegisterUserRequestBody struct {
	Login         string `json:"login"`
	Password      string `json:"password"`
	AgainPassword string `json:"again_password"`
}

type TodoListComment struct {
	Id          int    `json:"id" bson:"id"`
	AuthorName  string `json:"author_name" bson:"author_name"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Comment     string `json:"comment" bson:"comment"`
}
