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
	PutCustomCategory(customCategory *model.CustomCategory, userID string) error
	DeleteCustomCategory(customCategoryID int, userID string) error
}

type TransactionsRepository interface {
	GetTransaction(transactionSender *model.TransactionSender, transactionID int) (*model.TransactionSender, error)
	PostTransaction(transaction *model.TransactionReceiver, userID string) (sql.Result, error)
	GetMonthlyTransactionsList(userID string) ([]model.TransactionSender, error)
}
