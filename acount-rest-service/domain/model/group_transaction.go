package model

import "encoding/json"

type GroupTransactionsList struct {
	GroupTransactionsList []GroupTransactionSender `json:"transactions_list"`
}

type GroupTransactionSender struct {
	ID                 int        `json:"id"                   db:"id"`
	TransactionType    string     `json:"transaction_type"     db:"transaction_type"`
	UpdatedDate        DateTime   `json:"updated_date"         db:"updated_date"`
	TransactionDate    Date       `json:"transaction_date"     db:"transaction_date"`
	Shop               NullString `json:"shop"                 db:"shop"`
	Memo               NullString `json:"memo"                 db:"memo"`
	Amount             int        `json:"amount"               db:"amount"`
	UserID             string     `json:"user_id"              db:"user_id"`
	BigCategoryName    string     `json:"big_category_name"    db:"big_category_name"`
	MediumCategoryName NullString `json:"medium_category_name" db:"medium_category_name"`
	CustomCategoryName NullString `json:"custom_category_name" db:"custom_category_name"`
}

type GroupTransactionReceiver struct {
	TransactionType  string     `json:"transaction_type"   db:"transaction_type"   validate:"required,oneof=expense income"`
	TransactionDate  Date       `json:"transaction_date"   db:"transaction_date"   validate:"required,date"`
	Shop             NullString `json:"shop"               db:"shop"               validate:"omitempty,max=20,blank"`
	Memo             NullString `json:"memo"               db:"memo"               validate:"omitempty,max=50,blank"`
	Amount           int        `json:"amount"             db:"amount"             validate:"required,min=1"`
	BigCategoryID    int        `json:"big_category_id"    db:"big_category_id"    validate:"required,min=1,max=17,either_id"`
	MediumCategoryID NullInt64  `json:"medium_category_id" db:"medium_category_id" validate:"omitempty,min=1,max=99"`
	CustomCategoryID NullInt64  `json:"custom_category_id" db:"custom_category_id" validate:"omitempty,min=1"`
}

func NewGroupTransactionsList(groupTransactionsList []GroupTransactionSender) GroupTransactionsList {
	return GroupTransactionsList{GroupTransactionsList: groupTransactionsList}
}

func (t GroupTransactionReceiver) ShowTransactionReceiver() (string, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return string(b), err
	}
	return string(b), nil
}
