package model

type TodoListModel struct {
	Id          int    `json:"id"`
	AuthorId    int    `json:"author_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
