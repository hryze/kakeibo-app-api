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
	GroupTasksRepository
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
	SearchTodoList(todoSqlQuery string) ([]model.Todo, error)
}

type GroupTodoRepository interface {
	GetDailyImplementationGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error)
	GetDailyDueGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error)
	GetMonthlyImplementationGroupTodoList(firstDay time.Time, lastDay time.Time, groupID int) ([]model.GroupTodo, error)
	GetMonthlyDueGroupTodoList(firstDay time.Time, lastDay time.Time, groupID int) ([]model.GroupTodo, error)
	PostGroupTodo(groupTodo *model.GroupTodo, userID string, groupID int) (sql.Result, error)
	GetGroupTodo(groupTodoId int) (*model.GroupTodo, error)
	PutGroupTodo(groupTodo *model.GroupTodo, groupTodoID int) error
	DeleteGroupTodo(groupTodoID int) error
	SearchGroupTodoList(groupTodoSqlQuery string) ([]model.GroupTodo, error)
}

type GroupTasksRepository interface {
	GetGroupTasksUsersList(groupID int) ([]model.GroupTasksUser, error)
	GetGroupTasksListAssignedToUser(groupID int) ([]model.GroupTask, error)
	GetGroupTasksUser(groupTasksUser model.GroupTasksUser, groupID int) (*model.GroupTasksUser, error)
	PostGroupTasksUser(groupTasksUser model.GroupTasksUser, groupID int) (sql.Result, error)
	GetGroupTasksList(groupID int) ([]model.GroupTask, error)
	GetGroupTask(groupTasksID int) (*model.GroupTask, error)
	PostGroupTask(groupTask model.GroupTask, groupID int) (sql.Result, error)
	PutGroupTask(groupTask *model.GroupTask, groupTasksID int) error
	DeleteGroupTask(groupTasksID int) error
}
