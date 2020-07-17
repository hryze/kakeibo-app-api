package repository

import (
	"database/sql"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"
)

type DBRepository interface {
	AuthRepository
	CategoriesRepository
	TransactionsRepository
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
	GetMonthlyTransactionsList(userID string) ([]model.TransactionSender, error)
	PostTransaction(transaction *model.TransactionReceiver, userID string) (sql.Result, error)
	PutTransaction(transaction *model.TransactionReceiver, transactionID int) error
	DeleteTransaction(transactionID int) error
}
