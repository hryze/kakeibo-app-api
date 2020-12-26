package repository

import (
	"database/sql"
	"time"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type HealthRepository interface {
	PingMySQL() error
	PingRedis() error
}

type AuthRepository interface {
	GetUserID(sessionID string) (string, error)
}

type TodoRepository interface {
	GetDailyImplementationTodoList(date time.Time, userID string) ([]model.Todo, error)
	GetDailyDueTodoList(date time.Time, userID string) ([]model.Todo, error)
	GetMonthlyImplementationTodoList(firstDay time.Time, lastDay time.Time, userID string) ([]model.Todo, error)
	GetMonthlyDueTodoList(firstDay time.Time, lastDay time.Time, userID string) ([]model.Todo, error)
	GetExpiredTodoList(dueDate time.Time, userID string) (*model.ExpiredTodoList, error)
	GetTodo(todoId int) (*model.Todo, error)
	PostTodo(todo *model.Todo, userID string) (sql.Result, error)
	PutTodo(todo *model.Todo, todoID int) error
	DeleteTodo(todoID int) error
	SearchTodoList(todoSqlQuery string) ([]model.Todo, error)
}

type ShoppingListRepository interface {
	GetRegularShoppingList(userID string) (model.RegularShoppingList, error)
	GetRegularShoppingItem(regularShoppingItemID int) (model.RegularShoppingItem, error)
	GetShoppingListRelatedToRegularShoppingItem(todayShoppingItemID int, laterThanTodayShoppingItemID int) (model.ShoppingList, error)
	PostRegularShoppingItem(regularShoppingItem *model.RegularShoppingItem, userID string, today time.Time) (sql.Result, sql.Result, sql.Result, error)
	PutRegularShoppingItem(regularShoppingItem *model.RegularShoppingItem, regularShoppingItemID int, userID string, today time.Time) (sql.Result, sql.Result, error)
	PutRegularShoppingList(regularShoppingList model.RegularShoppingList, userID string, today time.Time) error
	DeleteRegularShoppingItem(regularShoppingItemID int) error
	GetDailyShoppingListByDay(date time.Time, userID string) (model.ShoppingList, error)
	GetDailyShoppingListByCategory(date time.Time, userID string) (model.ShoppingList, error)
	GetMonthlyShoppingListByDay(firstDay time.Time, lastDay time.Time, userID string) (model.ShoppingList, error)
	GetMonthlyShoppingListByCategory(firstDay time.Time, lastDay time.Time, userID string) (model.ShoppingList, error)
	GetExpiredShoppingList(dueDate time.Time, userID string) (model.ExpiredShoppingList, error)
	GetShoppingItem(shoppingItemID int) (model.ShoppingItem, error)
	PostShoppingItem(shoppingItem *model.ShoppingItem, userID string) (sql.Result, error)
	PutShoppingItem(shoppingItem *model.ShoppingItem) (sql.Result, error)
	DeleteShoppingItem(shoppingItemID int) error
}

type GroupTodoRepository interface {
	GetDailyImplementationGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error)
	GetDailyDueGroupTodoList(date time.Time, groupID int) ([]model.GroupTodo, error)
	GetMonthlyImplementationGroupTodoList(firstDay time.Time, lastDay time.Time, groupID int) ([]model.GroupTodo, error)
	GetMonthlyDueGroupTodoList(firstDay time.Time, lastDay time.Time, groupID int) ([]model.GroupTodo, error)
	GetExpiredGroupTodoList(dueDate time.Time, groupID int) (*model.ExpiredGroupTodoList, error)
	PostGroupTodo(groupTodo *model.GroupTodo, userID string, groupID int) (sql.Result, error)
	GetGroupTodo(groupTodoId int) (*model.GroupTodo, error)
	PutGroupTodo(groupTodo *model.GroupTodo, groupTodoID int) error
	DeleteGroupTodo(groupTodoID int) error
	SearchGroupTodoList(groupTodoSqlQuery string) ([]model.GroupTodo, error)
}

type GroupShoppingListRepository interface {
	GetGroupRegularShoppingList(groupID int) (model.GroupRegularShoppingList, error)
	GetGroupRegularShoppingItem(groupRegularShoppingItemID int) (model.GroupRegularShoppingItem, error)
	GetGroupShoppingListRelatedToGroupRegularShoppingItem(todayGroupShoppingItemID int, laterThanTodayGroupShoppingItemID int) (model.GroupShoppingList, error)
	PostGroupRegularShoppingItem(groupRegularShoppingItem *model.GroupRegularShoppingItem, groupID int, today time.Time) (sql.Result, sql.Result, sql.Result, error)
	PutGroupRegularShoppingItem(groupRegularShoppingItem *model.GroupRegularShoppingItem, groupRegularShoppingItemID int, groupID int, today time.Time) (sql.Result, sql.Result, error)
	PutGroupRegularShoppingList(groupRegularShoppingList model.GroupRegularShoppingList, groupID int, today time.Time) error
	DeleteGroupRegularShoppingItem(groupRegularShoppingItemID int) error
	GetDailyGroupShoppingListByDay(date time.Time, groupID int) (model.GroupShoppingList, error)
	GetDailyGroupShoppingListByCategory(date time.Time, groupID int) (model.GroupShoppingList, error)
	GetMonthlyGroupShoppingListByDay(firstDay time.Time, lastDay time.Time, groupID int) (model.GroupShoppingList, error)
	GetGroupShoppingItem(groupShoppingItemID int) (model.GroupShoppingItem, error)
	PostGroupShoppingItem(groupShoppingItem *model.GroupShoppingItem, groupID int) (sql.Result, error)
	PutGroupShoppingItem(groupShoppingItem *model.GroupShoppingItem) (sql.Result, error)
	DeleteGroupShoppingItem(groupShoppingItemID int) error
}

type GroupTasksRepository interface {
	GetGroupTasksUsersList(groupID int) ([]model.GroupTasksUser, error)
	GetGroupTasksListAssignedToUser(groupID int) ([]model.GroupTask, error)
	PutGroupTasksListAssignedToUser(groupTasksList []model.GroupTask, updateTaskIndexList []int) error
	PostGroupTasksUsersList(groupTasksUsersList model.GroupTasksUsersListReceiver, groupID int) error
	GetGroupTasksIDListAssignedToUser(groupTasksUsersIdList []int, groupID int) ([]int, error)
	DeleteGroupTasksUsersList(groupTasksUsersListReceiver model.GroupTasksUsersListReceiver, groupTasksIDList []int, groupID int) error
	GetGroupTasksList(groupID int) ([]model.GroupTask, error)
	GetGroupTask(groupTasksID int) (*model.GroupTask, error)
	PostGroupTask(groupTask model.GroupTask, groupID int) (sql.Result, error)
	PutGroupTask(groupTask *model.GroupTask, groupTasksID int) (sql.Result, error)
	DeleteGroupTask(groupTasksID int) error
}
