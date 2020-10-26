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
	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/account-rest-service/testutil"
)

type MockTransactionsRepository struct{}

func (t MockTransactionsRepository) GetMonthlyTransactionsList(userID string, firstDay time.Time, lastDay time.Time) ([]model.TransactionSender, error) {
	return []model.TransactionSender{
		{
			ID:                 3,
			TransactionType:    "expense",
			PostedDate:         time.Date(2020, 7, 15, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 7, 15, 16, 0, 0, 0, time.UTC),
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
			PostedDate:         time.Date(2020, 7, 10, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 7, 10, 16, 0, 0, 0, time.UTC),
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
			PostedDate:         time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
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

func (t MockTransactionsRepository) Get10LatestTransactionsList(userID string) (*model.TransactionsList, error) {
	return &model.TransactionsList{
		TransactionsList: []model.TransactionSender{
			{
				ID:                 1,
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 7, 10, 16, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 7, 10, 16, 0, 0, 0, time.UTC),
				TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
				Shop:               model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Memo:               model.NullString{NullString: sql.NullString{String: "セールで牛肉購入", Valid: true}},
				Amount:             4500,
				BigCategoryName:    "食費",
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "食料品", Valid: true}},
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
				BigCategoryName:    "日用品",
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
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
				BigCategoryName:    "食費",
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
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
				BigCategoryName:    "交通費",
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "電車", Valid: true}},
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
				BigCategoryName:    "住宅",
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家賃", Valid: true}},
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
				BigCategoryName:    "食費",
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "カフェ", Valid: true}},
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
				BigCategoryName:    "健康・医療",
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "フィットネス", Valid: true}},
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
				BigCategoryName:    "教養・教育",
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "参考書", Valid: true}},
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
				BigCategoryName:    "食費",
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
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
				BigCategoryName:    "日用品",
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "歯磨き粉", Valid: true}},
			},
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
		PostedDate:         time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
		UpdatedDate:        time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
		TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
		Shop:               model.NullString{NullString: sql.NullString{String: "ニトリ", Valid: true}},
		Memo:               model.NullString{NullString: sql.NullString{String: "ベッド購入", Valid: true}},
		Amount:             15000,
		BigCategoryName:    "日用品",
		MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
		CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
	}, nil
}

func (t MockTransactionsRepository) PutTransaction(transaction *model.TransactionReceiver, transactionID int) error {
	return nil
}

func (t MockTransactionsRepository) DeleteTransaction(transactionID int) error {
	return nil
}

func (t MockTransactionsRepository) SearchTransactionsList(query string) ([]model.TransactionSender, error) {
	return []model.TransactionSender{
		{
			ID:                 1,
			TransactionType:    "expense",
			PostedDate:         time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC),
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "ニトリ", Valid: true}},
			Memo:               model.NullString{NullString: sql.NullString{String: "ベッド購入", Valid: true}},
			Amount:             15000,
			BigCategoryName:    "日用品",
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
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
			BigCategoryName:    "食費",
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
		},
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

func TestDBHandler_Get10LatestTransactionsList(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		TransactionsRepo: MockTransactionsRepository{},
	}

	r := httptest.NewRequest("GET", "/transactions/latest", nil)
	w := httptest.NewRecorder()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.Get10LatestTransactionsList(w, r)

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

func TestDBHandler_PutTransaction(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		TransactionsRepo: MockTransactionsRepository{},
	}

	r := httptest.NewRequest("PUT", "/transactions/1", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutTransaction(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.TransactionSender{}, &model.TransactionSender{})
}

func TestDBHandler_DeleteTransaction(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		TransactionsRepo: MockTransactionsRepository{},
	}

	r := httptest.NewRequest("DELETE", "/transactions/1", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.DeleteTransaction(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteContentMsg{}, &DeleteContentMsg{})
}

func TestDBHandler_SearchTransactionsList(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		TransactionsRepo: MockTransactionsRepository{},
	}

	r := httptest.NewRequest("GET", "/transactions/search", nil)
	w := httptest.NewRecorder()

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

	h.SearchTransactionsList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.TransactionsList{}, &model.TransactionsList{})
}
