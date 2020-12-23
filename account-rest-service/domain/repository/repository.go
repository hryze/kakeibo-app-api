package repository

import (
	"database/sql"
	"time"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
)

type HealthRepository interface {
	PingMySQL() error
	PingRedis() error
}

type AuthRepository interface {
	GetUserID(sessionID string) (string, error)
}

type CategoriesRepository interface {
	GetBigCategoriesList() ([]model.BigCategory, error)
	GetMediumCategoriesList() ([]model.AssociatedCategory, error)
	GetCustomCategoriesList(userID string) ([]model.AssociatedCategory, error)
	FindCustomCategory(customCategory *model.CustomCategory, userID string) error
	PostCustomCategory(customCategory *model.CustomCategory, userID string) (sql.Result, error)
	PutCustomCategory(customCategory *model.CustomCategory) error
	GetBigCategoryID(customCategoryID int) (int, error)
	DeleteCustomCategory(previousCustomCategoryID int, replaceMediumCategoryID int) error
	GetCategoriesName(categoriesID model.CategoriesID) (*model.CategoriesName, error)
	GetCategoriesNameList(categoriesIDList []model.CategoriesID) ([]model.CategoriesName, error)
}

type TransactionsRepository interface {
	GetMonthlyTransactionsList(userID string, firstDay time.Time, lastDay time.Time) ([]model.TransactionSender, error)
	Get10LatestTransactionsList(userID string) (*model.TransactionsList, error)
	GetTransaction(transactionSender *model.TransactionSender, transactionID int) (*model.TransactionSender, error)
	PostTransaction(transaction *model.TransactionReceiver, userID string) (sql.Result, error)
	PutTransaction(transaction *model.TransactionReceiver, transactionID int) error
	DeleteTransaction(transactionID int) error
	SearchTransactionsList(query string) ([]model.TransactionSender, error)
	GetMonthlyTransactionTotalAmountByBigCategory(userID string, firstDay time.Time, lastDay time.Time) ([]model.TransactionTotalAmountByBigCategory, error)
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
	GetGroupMediumCategoriesList() ([]model.GroupAssociatedCategory, error)
	GetGroupCustomCategoriesList(groupID int) ([]model.GroupAssociatedCategory, error)
	FindGroupCustomCategory(groupCustomCategory *model.GroupCustomCategory, groupID int) error
	PostGroupCustomCategory(groupCustomCategory *model.GroupCustomCategory, groupID int) (sql.Result, error)
	PutGroupCustomCategory(groupCustomCategory *model.GroupCustomCategory) error
	FindGroupCustomCategoryID(groupCustomCategoryID int) error
	GetBigCategoryID(groupCustomCategoryID int) (int, error)
	DeleteGroupCustomCategory(previousGroupCustomCategoryID int, replaceMediumCategoryID int) error
	GetGroupCategoriesName(categoriesID model.CategoriesID) (*model.CategoriesName, error)
}

type GroupTransactionsRepository interface {
	GetMonthlyGroupTransactionsList(groupID int, firstDay time.Time, lastDay time.Time) ([]model.GroupTransactionSender, error)
	Get10LatestGroupTransactionsList(groupID int) (*model.GroupTransactionsList, error)
	GetGroupTransaction(groupTransactionID int) (*model.GroupTransactionSender, error)
	PostGroupTransaction(groupTransaction *model.GroupTransactionReceiver, groupID int, postedUserID string) (sql.Result, error)
	PutGroupTransaction(groupTransaction *model.GroupTransactionReceiver, groupTransactionID int, updatedUserID string) error
	DeleteGroupTransaction(groupTransactionID int) error
	SearchGroupTransactionsList(query string) ([]model.GroupTransactionSender, error)
	GetUserPaymentAmountList(groupID int, groupUserIDList []string, firstDay time.Time, lastDay time.Time) ([]model.UserPaymentAmount, error)
	GetGroupAccountsList(yearMonth time.Time, groupID int) ([]model.GroupAccount, error)
	PostGroupAccountsList(groupAccountsList []model.GroupAccount) error
	PutGroupAccountsList(groupAccountsList []model.GroupAccount) error
	DeleteGroupAccountsList(yearMonth time.Time, groupID int) error
	GetMonthlyGroupTransactionTotalAmountByBigCategory(groupID int, firstDay time.Time, lastDay time.Time) ([]model.GroupTransactionTotalAmountByBigCategory, error)
	YearlyGroupTransactionExistenceConfirmation(firstDayOfYear time.Time, groupID int) ([]time.Time, error)
	GetYearlyGroupAccountsList(firstDayOfYear time.Time, groupID int) ([]model.GroupAccount, error)
}

type GroupBudgetsRepository interface {
	PostInitGroupStandardBudgets(groupID int) error
	GetGroupStandardBudgets(groupID int) (*model.GroupStandardBudgets, error)
	PutGroupStandardBudgets(groupStandardBudgets *model.GroupStandardBudgets, groupID int) error
	GetGroupCustomBudgets(yearMonth time.Time, groupID int) (*model.GroupCustomBudgets, error)
	PostGroupCustomBudgets(groupCustomBudgets *model.GroupCustomBudgets, yearMonth time.Time, groupID int) error
	PutGroupCustomBudgets(groupCustomBudgets *model.GroupCustomBudgets, yearMonth time.Time, groupID int) error
	DeleteGroupCustomBudgets(yearMonth time.Time, groupID int) error
	GetMonthlyGroupStandardBudget(groupID int) (model.MonthlyGroupBudget, error)
	GetMonthlyGroupCustomBudgets(year time.Time, groupID int) ([]model.MonthlyGroupBudget, error)
}
