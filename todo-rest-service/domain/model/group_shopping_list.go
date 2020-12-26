package model

import (
	"errors"
	"time"
)

type GroupShoppingDataByDay struct {
	GroupRegularShoppingList
	GroupShoppingList
}

type GroupRegularShoppingList struct {
	GroupRegularShoppingList []GroupRegularShoppingItem `json:"regular_shopping_list"`
}

type GroupRegularShoppingItem struct {
	ID                   int        `json:"id"                     db:"id"`
	PostedDate           time.Time  `json:"posted_date"            db:"posted_date"`
	UpdatedDate          time.Time  `json:"updated_date"           db:"updated_date"`
	ExpectedPurchaseDate Date       `json:"expected_purchase_date" db:"expected_purchase_date"`
	CycleType            string     `json:"cycle_type"             db:"cycle_type"`
	Cycle                NullInt    `json:"cycle"                  db:"cycle"`
	Purchase             string     `json:"purchase"               db:"purchase"`
	Shop                 NullString `json:"shop"                   db:"shop"`
	Amount               NullInt64  `json:"amount"                 db:"amount"`
	BigCategoryID        int        `json:"big_category_id"        db:"big_category_id"`
	BigCategoryName      string     `json:"big_category_name"      db:"big_category_name"`
	MediumCategoryID     NullInt64  `json:"medium_category_id"     db:"medium_category_id"`
	MediumCategoryName   NullString `json:"medium_category_name"   db:"medium_category_name"`
	CustomCategoryID     NullInt64  `json:"custom_category_id"     db:"custom_category_id"`
	CustomCategoryName   NullString `json:"custom_category_name"   db:"custom_category_name"`
	PaymentUserID        NullString `json:"payment_user_id"        db:"payment_user_id"`
	TransactionAutoAdd   BitBool    `json:"transaction_auto_add"   db:"transaction_auto_add"`
}

type GroupShoppingList struct {
	GroupShoppingList []GroupShoppingItem `json:"shopping_list"`
}

type GroupShoppingItem struct {
	ID                     int                   `json:"id"                       db:"id"`
	PostedDate             time.Time             `json:"posted_date"              db:"posted_date"`
	UpdatedDate            time.Time             `json:"updated_date"             db:"updated_date"`
	ExpectedPurchaseDate   Date                  `json:"expected_purchase_date"   db:"expected_purchase_date"`
	CompleteFlag           BitBool               `json:"complete_flag"            db:"complete_flag"`
	Purchase               string                `json:"purchase"                 db:"purchase"`
	Shop                   NullString            `json:"shop"                     db:"shop"`
	Amount                 NullInt64             `json:"amount"                   db:"amount"`
	BigCategoryID          int                   `json:"big_category_id"          db:"big_category_id"`
	BigCategoryName        string                `json:"big_category_name"        db:"big_category_name"`
	MediumCategoryID       NullInt64             `json:"medium_category_id"       db:"medium_category_id"`
	MediumCategoryName     NullString            `json:"medium_category_name"     db:"medium_category_name"`
	CustomCategoryID       NullInt64             `json:"custom_category_id"       db:"custom_category_id"`
	CustomCategoryName     NullString            `json:"custom_category_name"     db:"custom_category_name"`
	RegularShoppingListID  NullInt64             `json:"regular_shopping_list_id" db:"regular_shopping_list_id"`
	PaymentUserID          NullString            `json:"payment_user_id"          db:"payment_user_id"`
	TransactionAutoAdd     BitBool               `json:"transaction_auto_add"     db:"transaction_auto_add"`
	RelatedTransactionData *GroupTransactionData `json:"related_transaction_data" db:"transaction_id"`
}

type GroupTransactionData struct {
	ID                 NullInt64  `json:"id"`
	TransactionType    string     `json:"transaction_type"`
	PostedDate         time.Time  `json:"posted_date"`
	UpdatedDate        time.Time  `json:"updated_date"`
	TransactionDate    string     `json:"transaction_date,omitempty"`
	Shop               NullString `json:"shop"`
	Memo               NullString `json:"memo"`
	Amount             int        `json:"amount,omitempty"`
	PostedUserID       string     `json:"posted_user_id"`
	UpdatedUserID      NullString `json:"updated_user_id"`
	PaymentUserID      string     `json:"payment_user_id"`
	BigCategoryID      int        `json:"big_category_id"`
	BigCategoryName    string     `json:"big_category_name"`
	MediumCategoryID   NullInt64  `json:"medium_category_id"`
	MediumCategoryName NullString `json:"medium_category_name"`
	CustomCategoryID   NullInt64  `json:"custom_category_id"`
	CustomCategoryName NullString `json:"custom_category_name"`
}

func (t *GroupTransactionData) Scan(value interface{}) error {
	id, ok := value.(int64)
	if !ok {
		return errors.New("bad int64 type assertion")
	}

	t.ID.Int64, t.ID.Valid = id, true

	return nil
}
