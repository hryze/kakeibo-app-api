package model

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type TransactionsList struct {
	TransactionsList []TransactionSender `json:"transactions_list"`
}

type TransactionSender struct {
	ID                 int        `json:"id"                   db:"id"`
	TransactionType    string     `json:"transaction_type"     db:"transaction_type"`
	UpdatedDate        DateTime   `json:"updated_date"         db:"updated_date"`
	TransactionDate    Date       `json:"transaction_date"     db:"transaction_date"`
	Shop               NullString `json:"shop"                 db:"shop"`
	Memo               NullString `json:"memo"                 db:"memo"`
	Amount             int        `json:"amount"               db:"amount"`
	BigCategoryName    string     `json:"big_category_name"    db:"big_category_name"`
	MediumCategoryName NullString `json:"medium_category_name" db:"medium_category_name"`
	CustomCategoryName NullString `json:"custom_category_name" db:"custom_category_name"`
}

type TransactionReceiver struct {
	TransactionType  string     `json:"transaction_type"   db:"transaction_type"`
	TransactionDate  Date       `json:"transaction_date"   db:"transaction_date"`
	Shop             NullString `json:"shop"               db:"shop"`
	Memo             NullString `json:"memo"               db:"memo"`
	Amount           int        `json:"amount"             db:"amount"`
	BigCategoryID    int        `json:"big_category_id"    db:"big_category_id"`
	MediumCategoryID NullInt64  `json:"medium_category_id" db:"medium_category_id"`
	CustomCategoryID NullInt64  `json:"custom_category_id" db:"custom_category_id"`
}

type DateTime struct {
	time.Time
}

type Date struct {
	time.Time
}

type NullString struct {
	sql.NullString
}

type NullInt64 struct {
	sql.NullInt64
}

func NewTransactionsList(transactionsList []TransactionSender) TransactionsList {
	return TransactionsList{TransactionsList: transactionsList}
}

func (dt *DateTime) Scan(value interface{}) error {
	dateTime, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("type assertion error")
	}
	dt.Time = dateTime
	return nil
}

func (dt DateTime) Value() (driver.Value, error) {
	return dt.Time, nil
}

func (d *Date) Scan(value interface{}) error {
	date, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("type assertion error")
	}
	d.Time = date
	return nil
}

func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	date := d.Time.Format("01/02")
	dayOfWeeks := [...]string{"日", "月", "火", "水", "木", "金", "土"}
	dayOfWeek := dayOfWeeks[d.Time.Weekday()]
	return []byte(`"` + date + `(` + dayOfWeek + `)` + `"`), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	trimData := strings.Trim(string(data), "\"")[:10]
	date, err := time.Parse("2006-01-02", trimData)
	if err != nil {
		return err
	}
	d.Time = date
	return nil
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}
	err := json.Unmarshal(b, &ns.String)
	if err == nil {
		ns.Valid = true
		return nil
	}
	return err
}

func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}
	err := json.Unmarshal(b, &ni.Int64)
	if err == nil {
		ni.Valid = true
		return nil
	}
	return err
}
