package repository

import (
	"database/sql"
	"time"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type DBRepository interface {
	AuthRepository
	TodoRepository
	GroupTodoRepository
}

type AuthRepository interface {
	GetUserID(sessionID string) (string, error)
}

type TodoRepository interface {
	GetDailyImplementationTodoList(date time.Time, userID string) ([]model.Todo, error)
	GetDailyDueTodoList(date time.Time, userID string) ([]model.Todo, error)
	GetMonthlyImplementationTodoList(firstDay time.Time, lastDay time.Time, userID string) ([]model.Todo, error)
	GetMonthlyDueTodoList(firstDay time.Time, lastDay time.Time, userID string) ([]model.Todo, error)
	GetTodo(todoId int) (*model.Todo, error)
	PostTodo(todo *model.Todo, userID string) (sql.Result, error)
	PutTodo(todo *model.Todo, todoID int) error
	DeleteTodo(todoID int) error
	SearchTodoList(query string) ([]model.Todo, error)
}

type GroupTodoRepository interface {
	GetDailyImplementationGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error)
	GetDailyDueGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error)
	GetMonthlyImplementationGroupTodoList(firstDay time.Time, lastDay time.Time, groupID int) ([]model.GroupTodo, error)
	GetMonthlyDueGroupTodoList(firstDay time.Time, lastDay time.Time, groupID int) ([]model.GroupTodo, error)
}
