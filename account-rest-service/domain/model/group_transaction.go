package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

type GroupTransactionsList struct {
	GroupTransactionsList []GroupTransactionSender `json:"transactions_list"`
}

type GroupTransactionSender struct {
	ID                 int        `json:"id"                   db:"id"`
	TransactionType    string     `json:"transaction_type"     db:"transaction_type"`
	PostedDate         time.Time  `json:"posted_date"          db:"posted_date"`
	UpdatedDate        time.Time  `json:"updated_date"         db:"updated_date"`
	TransactionDate    SenderDate `json:"transaction_date"     db:"transaction_date"`
	Shop               NullString `json:"shop"                 db:"shop"`
	Memo               NullString `json:"memo"                 db:"memo"`
	Amount             int        `json:"amount"               db:"amount"`
	PostedUserID       string     `json:"posted_user_id"       db:"posted_user_id"`
	UpdatedUserID      NullString `json:"updated_user_id"      db:"updated_user_id"`
	PaymentUserID      string     `json:"payment_user_id"      db:"payment_user_id"`
	BigCategoryID      int        `json:"big_category_id"      db:"big_category_id"`
	BigCategoryName    string     `json:"big_category_name"    db:"big_category_name"`
	MediumCategoryID   NullInt64  `json:"medium_category_id"   db:"medium_category_id"`
	MediumCategoryName NullString `json:"medium_category_name" db:"medium_category_name"`
	CustomCategoryID   NullInt64  `json:"custom_category_id"   db:"custom_category_id"`
	CustomCategoryName NullString `json:"custom_category_name" db:"custom_category_name"`
}

type GroupTransactionReceiver struct {
	TransactionType  string       `json:"transaction_type"   db:"transaction_type"   validate:"required,oneof=expense income"`
	TransactionDate  ReceiverDate `json:"transaction_date"   db:"transaction_date"   validate:"required,date"`
	Shop             NullString   `json:"shop"               db:"shop"               validate:"omitempty,max=20,blank"`
	Memo             NullString   `json:"memo"               db:"memo"               validate:"omitempty,max=50,blank"`
	Amount           int          `json:"amount"             db:"amount"             validate:"required,min=1"`
	PaymentUserID    string       `json:"payment_user_id"    db:"payment_user_id"`
	BigCategoryID    int          `json:"big_category_id"    db:"big_category_id"    validate:"required,min=1,max=17,either_id"`
	MediumCategoryID NullInt64    `json:"medium_category_id" db:"medium_category_id" validate:"omitempty,min=1,max=99"`
	CustomCategoryID NullInt64    `json:"custom_category_id" db:"custom_category_id" validate:"omitempty,min=1"`
}

type GroupTransactionTotalAmountByBigCategory struct {
	BigCategoryID int `db:"big_category_id"`
	TotalAmount   int `db:"total_amount"`
}

type YearlyAccountingStatus struct {
	Year                   string
	YearlyAccountingStatus [12]MonthlyAccountingStatus `json:"yearly_accounting_status"`
}

type MonthlyAccountingStatus struct {
	Month             string        `json:"month"`
	CalculationStatus string        `json:"calculation_status"`
	PaymentStatus     PaymentStatus `json:"payment_status"`
	ReceiptStatus     ReceiptStatus `json:"receipt_status"`
}

type PaymentStatus struct {
	IncompleteCount     int
	WaitingReceiptCount int
	CompleteCount       int
}

type ReceiptStatus struct {
	WaitingPaymentCount int
	IncompleteCount     int
	CompleteCount       int
}

type GroupAccountsList struct {
	GroupID                       int                        `json:"group_id"`
	Month                         time.Time                  `json:"month"`
	GroupTotalPaymentAmount       int                        `json:"group_total_payment_amount"`
	GroupAveragePaymentAmount     int                        `json:"group_average_payment_amount"`
	GroupRemainingAmount          int                        `json:"group_remaining_amount"`
	GroupAccountsListByPayersList []GroupAccountsListByPayer `json:"group_accounts_list_by_payer"`
	GroupAccountsList             []GroupAccount             `json:"-"`
}

type GroupAccountsListByPayer struct {
	Payer             NullString     `json:"payer_user_id"`
	GroupAccountsList []GroupAccount `json:"group_accounts_list"`
}

type GroupAccount struct {
	ID                  int        `json:"id"                   db:"id"`
	GroupID             int        `json:"group_id"             db:"group_id"`
	Month               time.Time  `json:"month"                db:"years_months"`
	Payer               NullString `json:"payer_user_id"        db:"payer_user_id"`
	Recipient           NullString `json:"recipient_user_id"    db:"recipient_user_id"`
	PaymentAmount       NullInt    `json:"payment_amount"       db:"payment_amount"`
	PaymentConfirmation BitBool    `json:"payment_confirmation" db:"payment_confirmation"`
	ReceiptConfirmation BitBool    `json:"receipt_confirmation" db:"receipt_confirmation"`
}

type PayerList struct {
	PayerList []UserPaymentAmount
}

type RecipientList struct {
	RecipientList []UserPaymentAmount
}

type UserPaymentAmount struct {
	UserID              string `db:"user_id"`
	TotalPaymentAmount  int    `db:"total_payment_amount"`
	PaymentAmountToUser int
}

type BitBool bool

type NullInt struct {
	Int   int
	Valid bool
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

func NewPayerList(userPaymentAmountList []UserPaymentAmount) PayerList {
	var payerList PayerList
	for _, userPaymentAmount := range userPaymentAmountList {
		if userPaymentAmount.PaymentAmountToUser < 0 {
			payerList.PayerList = append(payerList.PayerList, userPaymentAmount)
		}
	}

	sort.Slice(payerList.PayerList, func(i, j int) bool {
		return payerList.PayerList[i].PaymentAmountToUser < payerList.PayerList[j].PaymentAmountToUser
	})

	return payerList
}

func NewRecipientList(userPaymentAmountList []UserPaymentAmount) RecipientList {
	var recipientList RecipientList
	for _, userPaymentAmount := range userPaymentAmountList {
		if userPaymentAmount.PaymentAmountToUser > 0 {
			recipientList.RecipientList = append(recipientList.RecipientList, userPaymentAmount)
		}
	}

	sort.Slice(recipientList.RecipientList, func(i, j int) bool {
		return recipientList.RecipientList[i].PaymentAmountToUser > recipientList.RecipientList[j].PaymentAmountToUser
	})

	return recipientList
}

func NewGroupAccountsList(userPaymentAmountList []UserPaymentAmount, groupID int, month time.Time) GroupAccountsList {
	var totalPaymentAmount int
	for _, userPaymentAmount := range userPaymentAmountList {
		totalPaymentAmount += userPaymentAmount.TotalPaymentAmount
	}

	averagePaymentAmount := int(math.Round((float64(totalPaymentAmount)) / float64(len(userPaymentAmountList))))
	remainingAmount := totalPaymentAmount - averagePaymentAmount*len(userPaymentAmountList)

	return GroupAccountsList{
		GroupID:                       groupID,
		Month:                         month,
		GroupTotalPaymentAmount:       totalPaymentAmount,
		GroupAveragePaymentAmount:     averagePaymentAmount,
		GroupRemainingAmount:          remainingAmount,
		GroupAccountsListByPayersList: make([]GroupAccountsListByPayer, 0),
	}
}

func (b BitBool) Value() (driver.Value, error) {
	if b {
		return []byte{1}, nil
	}

	return []byte{0}, nil
}

func (b *BitBool) Scan(src interface{}) error {
	bitBool, ok := src.([]byte)
	if !ok {
		return errors.New("bad []byte type assertion")
	}

	*b = bitBool[0] == 1

	return nil
}

func (ni *NullInt) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(ni.Int)
}

func (ni *NullInt) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}

	if err := json.Unmarshal(b, &ni.Int); err != nil {
		return err
	}

	ni.Valid = true

	return nil
}

func (ni *NullInt) Scan(value interface{}) error {
	if value == nil {
		ni.Int, ni.Valid = 0, false
		return nil
	}

	intValue, ok := value.(int64)
	if !ok {
		return errors.New("type assertion error")
	}

	ni.Int, ni.Valid = int(intValue), true

	return nil
}

func (ni NullInt) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}

	return int64(ni.Int), nil
}

func NewYearlyAccountingStatus(year time.Time, userID string, transactionExistenceByMonths []time.Time, yearlyGroupAccountsList []GroupAccount) YearlyAccountingStatus {
	yearlyAccountingStatus := YearlyAccountingStatus{
		Year: fmt.Sprintf("%d年", year.Year()),
		YearlyAccountingStatus: [12]MonthlyAccountingStatus{
			{Month: "1月", CalculationStatus: "-"},
			{Month: "2月", CalculationStatus: "-"},
			{Month: "3月", CalculationStatus: "-"},
			{Month: "4月", CalculationStatus: "-"},
			{Month: "5月", CalculationStatus: "-"},
			{Month: "6月", CalculationStatus: "-"},
			{Month: "7月", CalculationStatus: "-"},
			{Month: "8月", CalculationStatus: "-"},
			{Month: "9月", CalculationStatus: "-"},
			{Month: "10月", CalculationStatus: "-"},
			{Month: "11月", CalculationStatus: "-"},
			{Month: "12月", CalculationStatus: "-"},
		},
	}

	for _, existenceByMonth := range transactionExistenceByMonths {
		idx := existenceByMonth.Month() - 1
		yearlyAccountingStatus.YearlyAccountingStatus[idx].CalculationStatus = "未精算"
	}

	for _, groupAccount := range yearlyGroupAccountsList {
		idx := groupAccount.Month.Month() - 1
		yearlyAccountingStatus.YearlyAccountingStatus[idx].CalculationStatus = "精算済"

		if userID == groupAccount.Payer.String {
			if !groupAccount.PaymentConfirmation {
				yearlyAccountingStatus.YearlyAccountingStatus[idx].PaymentStatus.IncompleteCount++
			} else if groupAccount.PaymentConfirmation && !groupAccount.ReceiptConfirmation {
				yearlyAccountingStatus.YearlyAccountingStatus[idx].PaymentStatus.WaitingReceiptCount++
			} else if groupAccount.PaymentConfirmation && groupAccount.ReceiptConfirmation {
				yearlyAccountingStatus.YearlyAccountingStatus[idx].PaymentStatus.CompleteCount++
			}
		} else if userID == groupAccount.Recipient.String {
			if !groupAccount.PaymentConfirmation {
				yearlyAccountingStatus.YearlyAccountingStatus[idx].ReceiptStatus.WaitingPaymentCount++
			} else if groupAccount.PaymentConfirmation && !groupAccount.ReceiptConfirmation {
				yearlyAccountingStatus.YearlyAccountingStatus[idx].ReceiptStatus.IncompleteCount++
			} else if groupAccount.PaymentConfirmation && groupAccount.ReceiptConfirmation {
				yearlyAccountingStatus.YearlyAccountingStatus[idx].ReceiptStatus.CompleteCount++
			}
		}
	}

	return yearlyAccountingStatus
}

func (s *PaymentStatus) MarshalJSON() ([]byte, error) {
	var messages []string

	if s.IncompleteCount != 0 {
		incompleteMessage := fmt.Sprintf("未払い: %d件", s.IncompleteCount)
		messages = append(messages, incompleteMessage)
	}

	if s.WaitingReceiptCount != 0 {
		waitingReceiptMessage := fmt.Sprintf("受領待ち: %d件", s.WaitingReceiptCount)
		messages = append(messages, waitingReceiptMessage)
	}

	if s.CompleteCount != 0 {
		completeMessage := fmt.Sprintf("完了: %d件", s.CompleteCount)
		messages = append(messages, completeMessage)
	}

	if len(messages) != 0 {
		message := strings.Join(messages, " / ")
		return json.Marshal(message)
	}

	return json.Marshal("-")
}

func (s *ReceiptStatus) MarshalJSON() ([]byte, error) {
	var messages []string

	if s.WaitingPaymentCount != 0 {
		waitingPaymentMessage := fmt.Sprintf("支払待ち: %d件", s.WaitingPaymentCount)
		messages = append(messages, waitingPaymentMessage)
	}

	if s.IncompleteCount != 0 {
		incompleteMessage := fmt.Sprintf("未受領: %d件", s.IncompleteCount)
		messages = append(messages, incompleteMessage)
	}

	if s.CompleteCount != 0 {
		completeMessage := fmt.Sprintf("完了: %d件", s.CompleteCount)
		messages = append(messages, completeMessage)
	}

	if len(messages) != 0 {
		message := strings.Join(messages, " / ")
		return json.Marshal(message)
	}

	return json.Marshal("-")
}
