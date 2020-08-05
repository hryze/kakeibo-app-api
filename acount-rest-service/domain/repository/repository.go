package repository

import (
	"database/sql"
	"time"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"
)

type DBRepository interface {
	AuthRepository
	CategoriesRepository
	TransactionsRepository
	BudgetsRepository
	GroupCategoriesRepository
	GroupBudgetsRepository
}

type AuthRepository interface {
	GetUserID(sessionID string) (string, error)
}

type CategoriesRepository interface {
	GetBigCategoriesList() ([]model.BigCategory, error)
	GetMediumCategoriesList() ([]model.MediumCategory, error)
	GetCustomCategoriesList(userID string) ([]model.CustomCategory, error)
	FindCustomCategory(customCategory *model.CustomCategory, userID string) error
	PostCustomCategory(customCategory *model.CustomCategory, userID string) (sql.Result, error)
	PutCustomCategory(customCategory *model.CustomCategory) error
	DeleteCustomCategory(customCategoryID int) error
}

type TransactionsRepository interface {
	GetTransaction(transactionSender *model.TransactionSender, transactionID int) (*model.TransactionSender, error)
	GetMonthlyTransactionsList(userID string, firstDay time.Time, lastDay time.Time) ([]model.TransactionSender, error)
	PostTransaction(transaction *model.TransactionReceiver, userID string) (sql.Result, error)
	PutTransaction(transaction *model.TransactionReceiver, transactionID int) error
	DeleteTransaction(transactionID int) error
	SearchTransactionsList(query string) ([]model.TransactionSender, error)
}

type BudgetsRepository interface {
	PostInitStandardBudgets(userID string) error
	GetStandardBudgets(userID string) (*model.StandardBudgets, error)
	PutStandardBudgets(standardBudgets *model.StandardBudgets, userID string) error
	GetCustomBudgets(yearMonth time.Time, userID string) (*model.CustomBudgets, error)
	PostCustomBudgets(customBudgets *model.CustomBudgets, yearMonth time.Time, userID string) error
	PutCustomBudgets(customBudgets *model.CustomBudgets, yearMonth time.Time, userID string) error
	DeleteCustomBudgets(yearMonth time.Time, userID string) error
	GetMonthlyStandardBudget(userID string) (model.MonthlyBudget, error)
	GetMonthlyCustomBudgets(year time.Time, userID string) ([]model.MonthlyBudget, error)
}

type GroupCategoriesRepository interface {
	GetGroupBigCategoriesList() ([]model.GroupBigCategory, error)
	GetGroupMediumCategoriesList() ([]model.GroupMediumCategory, error)
	GetGroupCustomCategoriesList(groupID int) ([]model.GroupCustomCategory, error)
}

type GroupBudgetsRepository interface {
	PostInitGroupStandardBudgets(groupID int) error
}
