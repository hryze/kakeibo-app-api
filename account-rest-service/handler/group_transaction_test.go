package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/account-rest-service/testutil"
)

var (
	counter int64
	mu      sync.Mutex
)

type MockGroupTransactionsRepository struct{}

func (m MockGroupTransactionsRepository) GetMonthlyGroupTransactionsList(groupID int, firstDay time.Time, lastDay time.Time) ([]model.GroupTransactionSender, error) {
	return []model.GroupTransactionSender{
		{
			ID:                 1,
			TransactionType:    "expense",
			PostedDate:         time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "ニトリ", Valid: true}},
			Memo:               model.NullString{NullString: sql.NullString{String: "ベッド購入", Valid: true}},
			Amount:             15000,
			PostedUserID:       "userID1",
			UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			PaymentUserID:      "userID1",
			BigCategoryID:      3,
			BigCategoryName:    "日用品",
			MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 16, Valid: true}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
			CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		},
		{
			ID:                 2,
			TransactionType:    "income",
			PostedDate:         time.Date(2020, 7, 10, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 7, 10, 16, 0, 0, 0, time.UTC),
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Memo:               model.NullString{NullString: sql.NullString{String: "賞与", Valid: true}},
			Amount:             200000,
			PostedUserID:       "userID2",
			UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			PaymentUserID:      "userID2",
			BigCategoryID:      1,
			BigCategoryName:    "収入",
			MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 2, Valid: true}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "賞与", Valid: true}},
			CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		},
		{
			ID:                 3,
			TransactionType:    "expense",
			PostedDate:         time.Date(2020, 7, 15, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 7, 15, 16, 0, 0, 0, time.UTC),
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 15, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Memo:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Amount:             1300,
			PostedUserID:       "userID1",
			UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			PaymentUserID:      "userID1",
			BigCategoryID:      2,
			BigCategoryName:    "食費",
			MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
		},
	}, nil
}

func (m MockGroupTransactionsRepository) Get10LatestGroupTransactionsList(groupID int) (*model.GroupTransactionsList, error) {
	return &model.GroupTransactionsList{
		GroupTransactionsList: []model.GroupTransactionSender{
			{
				ID:                 1,
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 7, 10, 16, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 7, 10, 16, 0, 0, 0, time.UTC),
				TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
				Shop:               model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Memo:               model.NullString{NullString: sql.NullString{String: "セールで牛肉購入", Valid: true}},
				Amount:             4500,
				PostedUserID:       "userID1",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID1",
				BigCategoryID:      2,
				BigCategoryName:    "食費",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 6, Valid: true}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "食料品", Valid: true}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			},
			{
				ID:                 2,
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 7, 9, 16, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 7, 9, 16, 0, 0, 0, time.UTC),
				TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
				Shop:               model.NullString{NullString: sql.NullString{String: "ニトリ", Valid: true}},
				Memo:               model.NullString{NullString: sql.NullString{String: "ベッド購入", Valid: true}},
				Amount:             15000,
				PostedUserID:       "userID1",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID1",
				BigCategoryID:      3,
				BigCategoryName:    "日用品",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 16, Valid: true}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			},
			{
				ID:                 3,
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 7, 8, 16, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 7, 8, 16, 0, 0, 0, time.UTC),
				TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
				Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				Memo:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				Amount:             1300,
				PostedUserID:       "userID2",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID2",
				BigCategoryID:      2,
				BigCategoryName:    "食費",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
			},
			{
				ID:                 4,
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 7, 7, 16, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 7, 7, 16, 0, 0, 0, time.UTC),
				TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
				Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				Memo:               model.NullString{NullString: sql.NullString{String: "電車定期代", Valid: true}},
				Amount:             12000,
				PostedUserID:       "userID2",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID2",
				BigCategoryID:      6,
				BigCategoryName:    "交通費",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 33, Valid: true}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "電車", Valid: true}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			},
			{
				ID:                 5,
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 7, 6, 16, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 7, 6, 16, 0, 0, 0, time.UTC),
				TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
				Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				Memo:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				Amount:             65000,
				PostedUserID:       "userID2",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID2",
				BigCategoryID:      11,
				BigCategoryName:    "住宅",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 66, Valid: true}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家賃", Valid: true}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			},
			{
				ID:                 6,
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 7, 5, 16, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 7, 5, 16, 0, 0, 0, time.UTC),
				TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
				Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				Memo:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				Amount:             500,
				PostedUserID:       "userID3",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID3",
				BigCategoryID:      2,
				BigCategoryName:    "食費",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 11, Valid: true}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "カフェ", Valid: true}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			},
			{
				ID:                 7,
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 7, 4, 16, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 7, 4, 16, 0, 0, 0, time.UTC),
				TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
				Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				Memo:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				Amount:             4800,
				PostedUserID:       "userID3",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID3",
				BigCategoryID:      8,
				BigCategoryName:    "健康・医療",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 49, Valid: true}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "フィットネス", Valid: true}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			},
			{
				ID:                 8,
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 7, 3, 16, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 7, 3, 16, 0, 0, 0, time.UTC),
				TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
				Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				Memo:               model.NullString{NullString: sql.NullString{String: "みんなのGo言語", Valid: true}},
				Amount:             2500,
				PostedUserID:       "userID3",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID3",
				BigCategoryID:      10,
				BigCategoryName:    "教養・教育",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 60, Valid: true}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "参考書", Valid: true}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			},
			{
				ID:                 9,
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 7, 2, 16, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 7, 2, 16, 0, 0, 0, time.UTC),
				TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
				Shop:               model.NullString{NullString: sql.NullString{String: "コンビニ", Valid: true}},
				Memo:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				Amount:             120,
				PostedUserID:       "userID1",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID1",
				BigCategoryID:      2,
				BigCategoryName:    "食費",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 2, Valid: true}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "パン", Valid: true}},
			},
			{
				ID:                 10,
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
				TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
				Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				Memo:               model.NullString{NullString: sql.NullString{String: "歯磨き粉3つ購入", Valid: true}},
				Amount:             300,
				PostedUserID:       "userID1",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID1",
				BigCategoryID:      3,
				BigCategoryName:    "日用品",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 3, Valid: true}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "歯磨き粉", Valid: true}},
			},
		},
	}, nil
}

func (m MockGroupTransactionsRepository) GetGroupTransaction(groupTransactionID int) (*model.GroupTransactionSender, error) {
	if groupTransactionID == 1 {
		return &model.GroupTransactionSender{
			ID:                 1,
			TransactionType:    "expense",
			PostedDate:         time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "ニトリ", Valid: true}},
			Memo:               model.NullString{NullString: sql.NullString{String: "ベッド購入", Valid: true}},
			Amount:             15000,
			PostedUserID:       "userID1",
			UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			PaymentUserID:      "userID1",
			BigCategoryID:      3,
			BigCategoryName:    "日用品",
			MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 16, Valid: true}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
			CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		}, nil
	}

	return &model.GroupTransactionSender{
		ID:                 2,
		TransactionType:    "expense",
		PostedDate:         time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
		UpdatedDate:        time.Date(2020, 7, 2, 16, 0, 0, 0, time.UTC),
		TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
		Shop:               model.NullString{NullString: sql.NullString{String: "ニトリ", Valid: true}},
		Memo:               model.NullString{NullString: sql.NullString{String: "ベッド購入", Valid: true}},
		Amount:             25000,
		PostedUserID:       "userID1",
		UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "userID2", Valid: true}},
		PaymentUserID:      "userID1",
		BigCategoryID:      3,
		BigCategoryName:    "日用品",
		MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 16, Valid: true}},
		MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
		CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
		CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
	}, nil
}

func (m MockGroupTransactionsRepository) PostGroupTransaction(groupTransaction *model.GroupTransactionReceiver, groupID int, postedUserID string) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (m MockGroupTransactionsRepository) PutGroupTransaction(groupTransaction *model.GroupTransactionReceiver, groupTransactionID int, updatedUserID string) error {
	return nil
}

func (m MockGroupTransactionsRepository) DeleteGroupTransaction(groupTransactionID int) error {
	return nil
}

func (m MockGroupTransactionsRepository) SearchGroupTransactionsList(query string) ([]model.GroupTransactionSender, error) {
	return []model.GroupTransactionSender{
		{
			ID:                 1,
			TransactionType:    "expense",
			PostedDate:         time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "ニトリ", Valid: true}},
			Memo:               model.NullString{NullString: sql.NullString{String: "ベッド購入", Valid: true}},
			Amount:             15000,
			PostedUserID:       "userID1",
			UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			PaymentUserID:      "userID1",
			BigCategoryID:      3,
			BigCategoryName:    "日用品",
			MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 16, Valid: true}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
			CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		},
		{
			ID:                 3,
			TransactionType:    "expense",
			PostedDate:         time.Date(2020, 7, 15, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 7, 15, 16, 0, 0, 0, time.UTC),
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 15, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Memo:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Amount:             1300,
			PostedUserID:       "userID1",
			UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			PaymentUserID:      "userID1",
			BigCategoryID:      2,
			BigCategoryName:    "食費",
			MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
		},
	}, nil
}

func (m MockGroupTransactionsRepository) GetGroupShoppingItemRelatedTransactionDataList(transactionIdList []int) ([]model.GroupTransactionSender, error) {
	return []model.GroupTransactionSender{
		{
			ID:                 1,
			TransactionType:    "expense",
			PostedDate:         time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "ニトリ", Valid: true}},
			Memo:               model.NullString{NullString: sql.NullString{String: "ベッド購入", Valid: true}},
			Amount:             15000,
			PostedUserID:       "userID1",
			UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			PaymentUserID:      "userID1",
			BigCategoryID:      3,
			BigCategoryName:    "日用品",
			MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 16, Valid: true}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
			CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		},
		{
			ID:                 2,
			TransactionType:    "expense",
			PostedDate:         time.Date(2020, 7, 15, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 7, 15, 16, 0, 0, 0, time.UTC),
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 15, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Memo:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Amount:             1300,
			PostedUserID:       "userID1",
			UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			PaymentUserID:      "userID1",
			BigCategoryID:      2,
			BigCategoryName:    "食費",
			MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
		},
	}, nil
}

func (m MockGroupTransactionsRepository) GetUserPaymentAmountList(groupID int, groupUserIDList []string, firstDay time.Time, lastDay time.Time) ([]model.UserPaymentAmount, error) {
	return []model.UserPaymentAmount{
		{UserID: "userID1", TotalPaymentAmount: 60000, PaymentAmountToUser: 0},
		{UserID: "userID4", TotalPaymentAmount: 45000, PaymentAmountToUser: 0},
		{UserID: "userID5", TotalPaymentAmount: 30000, PaymentAmountToUser: 0},
		{UserID: "userID3", TotalPaymentAmount: 7000, PaymentAmountToUser: 0},
		{UserID: "userID2", TotalPaymentAmount: 6000, PaymentAmountToUser: 0},
	}, nil
}

func (m MockGroupTransactionsRepository) GetGroupAccountsList(yearMonth time.Time, groupID int) ([]model.GroupAccount, error) {
	if groupID == 2 {
		return []model.GroupAccount{
			{
				ID:                  1,
				GroupID:             2,
				Month:               time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
				Payer:               model.NullString{NullString: sql.NullString{String: "userID2", Valid: true}},
				Recipient:           model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				PaymentAmount:       model.NullInt{Int: 23600, Valid: true},
				PaymentConfirmation: false,
				ReceiptConfirmation: false,
			},
			{
				ID:                  2,
				GroupID:             2,
				Month:               time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
				Payer:               model.NullString{NullString: sql.NullString{String: "userID3", Valid: true}},
				Recipient:           model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				PaymentAmount:       model.NullInt{Int: 6800, Valid: true},
				PaymentConfirmation: false,
				ReceiptConfirmation: false,
			},
			{
				ID:                  3,
				GroupID:             2,
				Month:               time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
				Payer:               model.NullString{NullString: sql.NullString{String: "userID3", Valid: true}},
				Recipient:           model.NullString{NullString: sql.NullString{String: "userID4", Valid: true}},
				PaymentAmount:       model.NullInt{Int: 15400, Valid: true},
				PaymentConfirmation: false,
				ReceiptConfirmation: false,
			},
			{
				ID:                  4,
				GroupID:             2,
				Month:               time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
				Payer:               model.NullString{NullString: sql.NullString{String: "userID3", Valid: true}},
				Recipient:           model.NullString{NullString: sql.NullString{String: "userID5", Valid: true}},
				PaymentAmount:       model.NullInt{Int: 400, Valid: true},
				PaymentConfirmation: false,
				ReceiptConfirmation: false,
			},
		}, nil
	}

	if groupID == 3 {
		if counter == 1 {
			atomic.AddInt64(&counter, -1)

			return []model.GroupAccount{
				{
					ID:                  1,
					GroupID:             3,
					Month:               time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
					Payer:               model.NullString{NullString: sql.NullString{String: "userID2", Valid: true}},
					Recipient:           model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
					PaymentAmount:       model.NullInt{Int: 23600, Valid: true},
					PaymentConfirmation: false,
					ReceiptConfirmation: false,
				},
				{
					ID:                  2,
					GroupID:             3,
					Month:               time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
					Payer:               model.NullString{NullString: sql.NullString{String: "userID3", Valid: true}},
					Recipient:           model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
					PaymentAmount:       model.NullInt{Int: 6800, Valid: true},
					PaymentConfirmation: false,
					ReceiptConfirmation: false,
				},
				{
					ID:                  3,
					GroupID:             3,
					Month:               time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
					Payer:               model.NullString{NullString: sql.NullString{String: "userID3", Valid: true}},
					Recipient:           model.NullString{NullString: sql.NullString{String: "userID4", Valid: true}},
					PaymentAmount:       model.NullInt{Int: 15400, Valid: true},
					PaymentConfirmation: false,
					ReceiptConfirmation: false,
				},
				{
					ID:                  4,
					GroupID:             3,
					Month:               time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
					Payer:               model.NullString{NullString: sql.NullString{String: "userID3", Valid: true}},
					Recipient:           model.NullString{NullString: sql.NullString{String: "userID5", Valid: true}},
					PaymentAmount:       model.NullInt{Int: 400, Valid: true},
					PaymentConfirmation: false,
					ReceiptConfirmation: false,
				},
			}, nil
		}

		atomic.AddInt64(&counter, 1)

		return make([]model.GroupAccount, 0), nil
	}

	return make([]model.GroupAccount, 0), nil
}

func (m MockGroupTransactionsRepository) PostGroupAccountsList(groupAccountsList []model.GroupAccount) error {
	return nil
}

func (m MockGroupTransactionsRepository) PutGroupAccountsList(groupAccountsList []model.GroupAccount) error {
	return nil
}

func (m MockGroupTransactionsRepository) DeleteGroupAccountsList(yearMonth time.Time, groupID int) error {
	return nil
}

func (m MockGroupTransactionsRepository) GetMonthlyGroupTransactionTotalAmountByBigCategory(groupID int, firstDay time.Time, lastDay time.Time) ([]model.GroupTransactionTotalAmountByBigCategory, error) {
	return []model.GroupTransactionTotalAmountByBigCategory{
		{
			BigCategoryID: 2,
			TotalAmount:   55000,
		},
		{
			BigCategoryID: 3,
			TotalAmount:   5000,
		},
		{
			BigCategoryID: 9,
			TotalAmount:   7000,
		},
		{
			BigCategoryID: 12,
			TotalAmount:   13000,
		},
		{
			BigCategoryID: 15,
			TotalAmount:   12000,
		},
	}, nil
}

func (m MockGroupTransactionsRepository) YearlyGroupTransactionExistenceConfirmation(firstDayOfYear time.Time, groupID int) ([]time.Time, error) {
	return []time.Time{
		time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC),
	}, nil
}

func (m MockGroupTransactionsRepository) GetYearlyGroupAccountsList(firstDayOfYear time.Time, groupID int) ([]model.GroupAccount, error) {
	return []model.GroupAccount{
		{
			Month:               time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
			Payer:               model.NullString{NullString: sql.NullString{String: "userID2", Valid: true}},
			Recipient:           model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
			PaymentConfirmation: false,
			ReceiptConfirmation: false,
		},
		{
			Month:               time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
			Payer:               model.NullString{NullString: sql.NullString{String: "userID3", Valid: true}},
			Recipient:           model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
			PaymentConfirmation: true,
			ReceiptConfirmation: false,
		},
		{
			Month:               time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
			Payer:               model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
			Recipient:           model.NullString{NullString: sql.NullString{String: "userID2", Valid: true}},
			PaymentConfirmation: true,
			ReceiptConfirmation: false,
		},
		{
			Month:               time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
			Payer:               model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
			Recipient:           model.NullString{NullString: sql.NullString{String: "userID3", Valid: true}},
			PaymentConfirmation: true,
			ReceiptConfirmation: true,
		},
		{
			Month:               time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC),
			Payer:               model.NullString{NullString: sql.NullString{String: "userID3", Valid: true}},
			Recipient:           model.NullString{NullString: sql.NullString{String: "userID2", Valid: true}},
			PaymentConfirmation: true,
			ReceiptConfirmation: false,
		},
		{
			Month:               time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC),
			Payer:               model.NullString{NullString: sql.NullString{String: "userID3", Valid: true}},
			Recipient:           model.NullString{NullString: sql.NullString{String: "userID4", Valid: true}},
			PaymentConfirmation: false,
			ReceiptConfirmation: false,
		},
	}, nil
}

func TestDBHandler_GetMonthlyGroupTransactionsList(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/1/transactions/2020-07", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id":   "1",
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetMonthlyGroupTransactionsList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupTransactionsList{}, &model.GroupTransactionsList{})
}

func TestDBHandler_Get10LatestGroupTransactionsList(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/1/transactions/latest", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.Get10LatestGroupTransactionsList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupTransactionsList{}, &model.GroupTransactionsList{})
}

func TestDBHandler_PostGroupTransaction(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("POST", "/groups/1/transactions", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostGroupTransaction(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.GroupTransactionSender{}, &model.GroupTransactionSender{})
}

func TestDBHandler_PutGroupTransaction(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("PUT", "/groups/1/transactions/2", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"id":       "2",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutGroupTransaction(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupTransactionSender{}, &model.GroupTransactionSender{})
}

func TestDBHandler_DeleteGroupTransaction(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("DELETE", "/groups/1/transactions/1", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"id":       "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.DeleteGroupTransaction(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteContentMsg{}, &DeleteContentMsg{})
}

func TestDBHandler_SearchGroupTransactionsList(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/1/transactions/search", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	urlQuery := r.URL.Query()

	params := map[string]string{
		"start_date":       "2020-07-01T00:00:00.0000",
		"end_date":         "2020-07-15T00:00:00.0000",
		"transaction_type": "expense",
		"sort":             "amount",
		"sort_type":        "desc",
	}

	for k, v := range params {
		urlQuery.Add(k, v)
	}

	r.URL.RawQuery = urlQuery.Encode()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.SearchGroupTransactionsList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupTransactionsList{}, &model.GroupTransactionsList{})
}

func TestDBHandler_GetGroupShoppingItemRelatedTransactionDataList(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/1/transactions/related-shopping-list", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetGroupShoppingItemRelatedTransactionDataList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &[]model.GroupTransactionSender{}, &[]model.GroupTransactionSender{})
}

func TestDBHandler_GetYearlyAccountingStatus(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/1/transactions/2020/account", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"year":     "2020",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetYearlyAccountingStatus(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)

	goldenFilePath := filepath.Join("testdata", t.Name(), "response.json.golden")

	wantData, err := ioutil.ReadFile(goldenFilePath)
	if err != nil {
		t.Fatalf("unexpected error by ioutil.ReadFile '%#v'", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("unexpected error by ioutil.ReadAll() '%#v'", err)
	}

	var gotData bytes.Buffer
	if err = json.Indent(&gotData, body, "", "  "); err != nil {
		t.Fatalf("unexpected error by json.Indent '%#v'", err)
	}

	if diff := cmp.Diff(string(wantData), gotData.String()); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func TestDBHandler_GetMonthlyGroupTransactionsAccount(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/2/transactions/2020-07/account", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id":   "2",
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetMonthlyGroupTransactionsAccount(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupAccountsList{}, &model.GroupAccountsList{})
}

func TestDBHandler_PostMonthlyGroupTransactionsAccount(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("POST", "/groups/2/transactions/2020-07/account", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id":   "3",
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	mu.Lock()
	defer mu.Unlock()

	h.PostMonthlyGroupTransactionsAccount(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.GroupAccountsList{}, &model.GroupAccountsList{})
}

func TestDBHandler_PutMonthlyGroupTransactionsAccount(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("PUT", "/groups/3/transactions/2020-07/account", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id":   "2",
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutMonthlyGroupTransactionsAccount(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupAccountsList{}, &model.GroupAccountsList{})
}

func TestDBHandler_DeleteMonthlyGroupTransactionsAccount(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("DELETE", "/groups/2/transactions/2020-07/account", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id":   "2",
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.DeleteMonthlyGroupTransactionsAccount(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteContentMsg{}, &DeleteContentMsg{})
}
