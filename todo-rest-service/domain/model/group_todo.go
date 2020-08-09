package model

import (
	"encoding/json"
	"time"
)

type GroupTodoList struct {
	ImplementationGroupTodoList []GroupTodo `json:"implementation_todo_list"`
	DueGroupTodoList            []GroupTodo `json:"due_todo_list"`
}

type SearchGroupTodoList struct {
	SearchGroupTodoList []GroupTodo `json:"search_todo_list"`
}

type GroupTodo struct {
	ID                 int       `json:"id"                  db:"id"`
	PostedDate         time.Time `json:"posted_date"         db:"posted_date"`
	ImplementationDate Date      `json:"implementation_date" db:"implementation_date" validate:"required,date"`
	DueDate            Date      `json:"due_date"            db:"due_date"            validate:"required,date"`
	TodoContent        string    `json:"todo_content"        db:"todo_content"        validate:"required,max=100,blank"`
	CompleteFlag       BitBool   `json:"complete_flag"       db:"complete_flag"`
	UserID             string    `json:"user_id"             db:"user_id"`
}

func NewGroupTodoList(implementationGroupTodoList []GroupTodo, dueGroupTodoList []GroupTodo) GroupTodoList {
	return GroupTodoList{
		ImplementationGroupTodoList: implementationGroupTodoList,
		DueGroupTodoList:            dueGroupTodoList,
	}
}

func NewSearchGroupTodoList(searchGroupTodoList []GroupTodo) SearchGroupTodoList {
	return SearchGroupTodoList{
		SearchGroupTodoList: searchGroupTodoList,
	}
}

func (t GroupTodo) ShowTodo() (string, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return string(b), err
	}
	return string(b), nil
}
