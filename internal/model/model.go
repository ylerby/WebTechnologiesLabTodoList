package model

type TodoListModel struct {
	Id          int    `json:"id" bson:"id"`
	AuthorId    int    `json:"author_id" bson:"author_id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

type UserModel struct {
	Login    string `gorm:"primaryKey;type:varchar(255)"`
	Password string `gorm:"type:varchar(255)"`
}
