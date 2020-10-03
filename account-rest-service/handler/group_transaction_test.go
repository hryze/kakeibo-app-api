package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/account-rest-service/testutil"
)

type MockGroupTransactionsRepository struct{}

func (m MockGroupTransactionsRepository) GetMonthlyGroupTransactionsList(groupID int, firstDay time.Time, lastDay time.Time) ([]model.GroupTransactionSender, error) {
	return []model.GroupTransactionSender{
		{
			ID:                 3,
			TransactionType:    "expense",
			UpdatedDate:        model.DateTime{Time: time.Date(2020, 7, 15, 16, 0, 0, 0, time.UTC)},
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 15, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Memo:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Amount:             1300,
			UserID:             "userID1",
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
			UserID:             "userID2",
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
			UserID:             "userID1",
			BigCategoryName:    "日用品",
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		},
	}, nil
}

func (m MockGroupTransactionsRepository) GetGroupTransaction(groupTransactionID int) (*model.GroupTransactionSender, error) {
	return &model.GroupTransactionSender{
		ID:                 1,
		TransactionType:    "expense",
		UpdatedDate:        model.DateTime{Time: time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC)},
		TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
		Shop:               model.NullString{NullString: sql.NullString{String: "ニトリ", Valid: true}},
		Memo:               model.NullString{NullString: sql.NullString{String: "ベッド購入", Valid: true}},
		Amount:             15000,
		UserID:             "userID1",
		BigCategoryName:    "日用品",
		MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
		CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
	}, nil
}

func (m MockGroupTransactionsRepository) PostGroupTransaction(groupTransaction *model.GroupTransactionReceiver, groupID int, userID string) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (m MockGroupTransactionsRepository) PutGroupTransaction(groupTransaction *model.GroupTransactionReceiver, groupTransactionID int) error {
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
			UpdatedDate:        model.DateTime{Time: time.Date(2020, 7, 1, 16, 0, 0, 0, time.UTC)},
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "ニトリ", Valid: true}},
			Memo:               model.NullString{NullString: sql.NullString{String: "ベッド購入", Valid: true}},
			Amount:             15000,
			UserID:             "userID1",
			BigCategoryName:    "日用品",
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "家具", Valid: true}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		},
		{
			ID:                 3,
			TransactionType:    "expense",
			UpdatedDate:        model.DateTime{Time: time.Date(2020, 7, 15, 16, 0, 0, 0, time.UTC)},
			TransactionDate:    model.SenderDate{Time: time.Date(2020, 7, 15, 0, 0, 0, 0, time.UTC)},
			Shop:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Memo:               model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			Amount:             1300,
			UserID:             "userID1",
			BigCategoryName:    "食費",
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
		},
	}, nil
}

func (m MockGroupTransactionsRepository) GetUserPaymentAmountList(groupID int, firstDay time.Time, lastDay time.Time) ([]model.UserPaymentAmount, error) {
	return []model.UserPaymentAmount{
		{UserID: "userID1", TotalPaymentAmount: 60000, PaymentAmountToUser: 0},
		{UserID: "userID4", TotalPaymentAmount: 45000, PaymentAmountToUser: 0},
		{UserID: "userID5", TotalPaymentAmount: 30000, PaymentAmountToUser: 0},
		{UserID: "userID3", TotalPaymentAmount: 7000, PaymentAmountToUser: 0},
		{UserID: "userID2", TotalPaymentAmount: 6000, PaymentAmountToUser: 0},
	}, nil
}

func (m MockGroupTransactionsRepository) GetGroupAccountsList(yearMonth time.Time, groupID int) ([]model.GroupAccount, error) {
	if groupID == 1 {
		return make([]model.GroupAccount, 0), nil
	}

	return []model.GroupAccount{
		{ID: 1, GroupID: 2, Month: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC), Payer: "userID2", Recipient: "userID1", PaymentAmount: 23600, PaymentConfirmation: false, ReceiptConfirmation: false},
		{ID: 2, GroupID: 2, Month: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC), Payer: "userID3", Recipient: "userID1", PaymentAmount: 6800, PaymentConfirmation: false, ReceiptConfirmation: false},
		{ID: 3, GroupID: 2, Month: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC), Payer: "userID3", Recipient: "userID4", PaymentAmount: 15400, PaymentConfirmation: false, ReceiptConfirmation: false},
		{ID: 4, GroupID: 2, Month: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC), Payer: "userID3", Recipient: "userID5", PaymentAmount: 400, PaymentConfirmation: false, ReceiptConfirmation: false},
	}, nil
}

func (m MockGroupTransactionsRepository) PostGroupAccountsList(groupAccountsList []model.GroupAccount, yearMonth time.Time, groupID int) error {
	return nil
}

func (m MockGroupTransactionsRepository) PutGroupAccountsList(groupAccountsList []model.GroupAccount) error {
	return nil
}

func (m MockGroupTransactionsRepository) DeleteGroupAccountsList(yearMonth time.Time, groupID int) error {
	return nil
}

func TestDBHandler_GetMonthlyGroupTransactionsList(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

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

func TestDBHandler_PostGroupTransaction(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

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
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("PUT", "/groups/1/transactions/1", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
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

	h.PutGroupTransaction(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupTransactionSender{}, &model.GroupTransactionSender{})
}

func TestDBHandler_DeleteGroupTransaction(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

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
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

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

func TestDBHandler_GetMonthlyGroupTransactionsAccount(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

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
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("POST", "/groups/2/transactions/2020-07/account", nil)
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

	h.PostMonthlyGroupTransactionsAccount(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.GroupAccountsList{}, &model.GroupAccountsList{})
}

func TestDBHandler_PutMonthlyGroupTransactionsAccount(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupTransactionsRepo: MockGroupTransactionsRepository{},
	}

	r := httptest.NewRequest("PUT", "/groups/3/transactions/2020-07/account", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
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

	h.PutMonthlyGroupTransactionsAccount(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupAccountsList{}, &model.GroupAccountsList{})
}

func TestDBHandler_DeleteMonthlyGroupTransactionsAccount(t *testing.T) {
	tearDown := testutil.SetUpMockServer(t)
	defer tearDown()

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
