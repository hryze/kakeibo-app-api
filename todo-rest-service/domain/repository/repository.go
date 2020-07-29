package repository

import (
	"time"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type DBRepository interface {
	AuthRepository
	TodoRepository
}

type AuthRepository interface {
	GetUserID(sessionID string) (string, error)
}

type TodoRepository interface {
	GetDailyImplementationTodoList(date time.Time, userID string) ([]model.Todo, error)
	GetDailyDueTodoList(date time.Time, userID string) ([]model.Todo, error)
}
