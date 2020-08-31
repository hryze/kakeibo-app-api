package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/google/uuid"
	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/repository"
	"github.com/paypay3/kakeibo-app-api/acount-rest-service/testutil"
)

type MockTransactionsRepository struct {
	repository.TransactionsRepository
}

type MockSqlResult struct {
	sql.Result
}

func (r MockSqlResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (t MockTransactionsRepository) GetMonthlyTransactionsList(userID string, firstDay time.Time, lastDay time.Time) ([]model.TransactionSender, error) {
	return []model.TransactionSender{
		{
			ID:                 3,
			TransactionType:    "expense",
			UpdatedDate:        model.DateTime{Time: time.Date(2020, 7, 15, 16, 0, 0, 0, time.UTC)},
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 15, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Memo:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Amount:             1300,
			BigCategoryName:    "食費",
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
		},
		{
			ID:                 2,
			TransactionType:    "income",
			UpdatedDate:        model.DateTime{Time: time.Date(2020, 7, 10, 16, 0, 0, 0, time.UTC)},
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Memo:               model.NullString{NullString: sql.NullString{String: "賞与", Valid: true}},
			Amount:             200000,
			BigCategoryName:    "収入",
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "賞与", Valid: true}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		},
		{
			ID:                 1,
			TransactionType:    "expense",
			UpdatedDate:        model.DateTime{Time: time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC)},
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "ニトリ", Valid: true}},
			Memo:               model.NullString{NullString: sql.NullString{String: "ベッド購入", Valid: true}},
			Amount:             15000,
			BigCategoryName:    "日用品",
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		},
	}, nil
}

func (t MockTransactionsRepository) PostTransaction(transaction *model.TransactionReceiver, userID string) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (t MockTransactionsRepository) GetTransaction(transactionSender *model.TransactionSender, transactionID int) (*model.TransactionSender, error) {
	return &model.TransactionSender{
		ID:                 1,
		TransactionType:    "expense",
		UpdatedDate:        model.DateTime{Time: time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC)},
		TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
		Shop:               model.NullString{NullString: sql.NullString{String: "ニトリ", Valid: true}},
		Memo:               model.NullString{NullString: sql.NullString{String: "ベッド購入", Valid: true}},
		Amount:             15000,
		BigCategoryName:    "日用品",
		MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
		CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
	}, nil
}

func TestDBHandler_GetMonthlyTransactionsList(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		TransactionsRepo: MockTransactionsRepository{},
	}

	r := httptest.NewRequest("GET", "/transactions/2020-07", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetMonthlyTransactionsList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.TransactionsList{}, &model.TransactionsList{})
}

func TestDBHandler_PostTransaction(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		TransactionsRepo: MockTransactionsRepository{},
	}

	r := httptest.NewRequest("POST", "/transactions", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"year_month": "2020-07",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostTransaction(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.TransactionSender{}, &model.TransactionSender{})
}
