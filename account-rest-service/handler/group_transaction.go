package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

type GroupTransactionsSearchQuery struct {
	TransactionType string
	BigCategoryID   string
	Shop            string
	Memo            string
	LowAmount       string
	HighAmount      string
	StartDate       string
	EndDate         string
	Sort            string
	SortType        string
	Limit           string
	GroupID         string
	UsersID         []string
}

type GroupTransactionProcessLockErrorMsg struct {
	Message string `json:"message"`
}

func NewGroupTransactionsSearchQuery(urlQuery url.Values, groupID string) GroupTransactionsSearchQuery {
	startDate := trimDate(urlQuery.Get("start_date"))
	endDate := trimDate(urlQuery.Get("end_date"))

	return GroupTransactionsSearchQuery{
		TransactionType: urlQuery.Get("transaction_type"),
		BigCategoryID:   urlQuery.Get("big_category_id"),
		Shop:            urlQuery.Get("shop"),
		Memo:            urlQuery.Get("memo"),
		LowAmount:       urlQuery.Get("low_amount"),
		HighAmount:      urlQuery.Get("high_amount"),
		StartDate:       startDate,
		EndDate:         endDate,
		Sort:            urlQuery.Get("sort"),
		SortType:        urlQuery.Get("sort_type"),
		Limit:           urlQuery.Get("limit"),
		GroupID:         groupID,
		UsersID:         urlQuery["user_id"],
	}
}

func (e *GroupTransactionProcessLockErrorMsg) Error() string {
	return e.Message
}

func generateGroupTransactionsSqlQuery(searchQuery GroupTransactionsSearchQuery) (string, error) {
	query := `
        SELECT
            group_transactions.id id,
            group_transactions.transaction_type transaction_type,
            group_transactions.updated_date updated_date,
            group_transactions.transaction_date transaction_date,
            group_transactions.shop shop,
            group_transactions.memo memo,
            group_transactions.amount amount,
            group_transactions.user_id user_id,
            big_categories.category_name big_category_name,
            medium_categories.category_name medium_category_name,
            group_custom_categories.category_name custom_category_name
        FROM
            group_transactions
        LEFT JOIN
            big_categories
        ON
            group_transactions.big_category_id = big_categories.id
        LEFT JOIN
            medium_categories
        ON
            group_transactions.medium_category_id = medium_categories.id
        LEFT JOIN
            group_custom_categories
        ON
            group_transactions.custom_category_id = group_custom_categories.id
        WHERE
            group_transactions.group_id = {{.GroupID}}

        {{ if eq (len .UsersID) 1 }}
        {{ range $i, $UserID := .UsersID }}
        AND
            group_transactions.user_id = "{{ $UserID }}"
        {{ end }}
        {{ end }}

        {{ if gt (len .UsersID) 1 }}
        {{ range $i, $UserID := .UsersID }}
        {{ if eq $i 0}}
        AND
            group_transactions.user_id IN("{{ $UserID }}"
        {{ end }}
        {{ if gt $i 0 }}
        ,"{{ $UserID }}"
        {{ end }}
        {{ end }}
        {{ end }}
        {{ if gt (len .UsersID) 1 }}
        )
        {{ end }}

        {{ with $StartDate := .StartDate }}
        AND
            group_transactions.transaction_date >= "{{ $StartDate }}"
        {{ end }}

        {{ with $EndDate := .EndDate }}
        AND
            group_transactions.transaction_date <= "{{ $EndDate }}"
        {{ end }}

        {{ with $TransactionType := .TransactionType }}
        AND
            group_transactions.transaction_type = "{{ $TransactionType }}"
        {{ end }}

        {{ with $BigCategoryID := .BigCategoryID }}
        AND
            group_transactions.big_category_id = "{{ $BigCategoryID }}"
        {{ end }}

        {{ with $LowAmount := .LowAmount }}
        AND
            group_transactions.amount >= "{{ $LowAmount }}"
        {{ end }}

        {{ with $HighAmount := .HighAmount }}
        AND
            group_transactions.amount <= "{{ $HighAmount }}"
        {{ end }}

        {{ with $Shop := .Shop }}
        AND
            group_transactions.shop
        LIKE
            "%{{ $Shop }}%"
        {{ end }}

        {{ with $Memo := .Memo }}
        AND
            group_transactions.memo
        LIKE
            "%{{ $Memo }}%"
        {{ end }}

        {{ with $Sort := .Sort}}
        ORDER BY
            group_transactions.{{ $Sort }}
        {{ else }}
        ORDER BY
            group_transactions.transaction_date
        {{ end }}

        {{ with $SortType := .SortType}}
        {{ $SortType }}, group_transactions.updated_date DESC
        {{ else }}
        DESC, group_transactions.updated_date DESC
        {{ end }}

        {{ with $Limit := .Limit}}
        LIMIT
        {{ $Limit }}
        {{ end }}`

	var buffer bytes.Buffer
	queryTemplate, err := template.New("GroupTransactionsSqlQueryTemplate").Parse(query)
	if err != nil {
		return "", err
	}

	if err := queryTemplate.Execute(&buffer, searchQuery); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (h *DBHandler) GetMonthlyGroupTransactionsList(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	firstDay, err := time.Parse("2006-01", mux.Vars(r)["year_month"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	lastDay := time.Date(firstDay.Year(), firstDay.Month()+1, 1, 0, 0, 0, 0, firstDay.Location()).Add(-1 * time.Second)

	dbGroupTransactionsList, err := h.GroupTransactionsRepo.GetMonthlyGroupTransactionsList(groupID, firstDay, lastDay)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(dbGroupTransactionsList) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&NoContentMsg{"条件に一致する取引履歴は見つかりませんでした。"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	groupTransactionsList := model.NewGroupTransactionsList(dbGroupTransactionsList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupTransactionsList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PostGroupTransaction(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	var groupTransactionReceiver model.GroupTransactionReceiver
	if err := json.NewDecoder(r.Body).Decode(&groupTransactionReceiver); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	yearMonth := time.Date(groupTransactionReceiver.TransactionDate.Time.Year(), groupTransactionReceiver.TransactionDate.Time.Month(), 1, 0, 0, 0, 0, time.UTC)

	dbGroupAccountsList, err := h.GroupTransactionsRepo.GetGroupAccountsList(yearMonth, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	} else if len(dbGroupAccountsList) != 0 {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &GroupTransactionProcessLockErrorMsg{"当月のグループでの取引は会計済みのため追加できません。"}))
		return
	}

	if err := validateTransaction(&groupTransactionReceiver); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	result, err := h.GroupTransactionsRepo.PostGroupTransaction(&groupTransactionReceiver, groupID, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroupTransactionSender, err := h.GroupTransactionsRepo.GetGroupTransaction(int(lastInsertId))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(dbGroupTransactionSender); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PutGroupTransaction(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	var groupTransactionReceiver model.GroupTransactionReceiver
	if err := json.NewDecoder(r.Body).Decode(&groupTransactionReceiver); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	yearMonth := time.Date(groupTransactionReceiver.TransactionDate.Time.Year(), groupTransactionReceiver.TransactionDate.Time.Month(), 1, 0, 0, 0, 0, time.UTC)

	dbGroupAccountsList, err := h.GroupTransactionsRepo.GetGroupAccountsList(yearMonth, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	} else if len(dbGroupAccountsList) != 0 {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &GroupTransactionProcessLockErrorMsg{"当月のグループでの取引は会計済みのため更新できません。"}))
		return
	}

	if err := validateTransaction(&groupTransactionReceiver); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	groupTransactionID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"transaction ID を正しく指定してください。"}))
		return
	}

	if err := h.GroupTransactionsRepo.PutGroupTransaction(&groupTransactionReceiver, groupTransactionID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupTransactionSender, err := h.GroupTransactionsRepo.GetGroupTransaction(groupTransactionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"トランザクションを取得できませんでした。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(groupTransactionSender); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) DeleteGroupTransaction(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	groupTransactionID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"transaction ID を正しく指定してください。"}))
		return
	}

	groupTransaction, err := h.GroupTransactionsRepo.GetGroupTransaction(groupTransactionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"こちらのトランザクションは既に削除されています。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	yearMonth := time.Date(groupTransaction.TransactionDate.Time.Year(), groupTransaction.TransactionDate.Time.Month(), 1, 0, 0, 0, 0, time.UTC)

	dbGroupAccountsList, err := h.GroupTransactionsRepo.GetGroupAccountsList(yearMonth, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	} else if len(dbGroupAccountsList) != 0 {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &GroupTransactionProcessLockErrorMsg{"当月のグループでの取引は会計済みのため削除できません。"}))
		return
	}

	if err := h.GroupTransactionsRepo.DeleteGroupTransaction(groupTransactionID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&DeleteContentMsg{"トランザクションを削除しました。"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) SearchGroupTransactionsList(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	strGroupID := mux.Vars(r)["group_id"]

	groupID, err := strconv.Atoi(strGroupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	if err := r.ParseForm(); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	searchQuery := NewGroupTransactionsSearchQuery(r.Form, strGroupID)

	query, err := generateGroupTransactionsSqlQuery(searchQuery)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroupTransactionsList, err := h.GroupTransactionsRepo.SearchGroupTransactionsList(query)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(dbGroupTransactionsList) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&NoContentMsg{"条件に一致する取引履歴は見つかりませんでした。"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	groupTransactionsList := model.NewGroupTransactionsList(dbGroupTransactionsList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupTransactionsList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetMonthlyGroupTransactionsAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	firstDay, err := time.Parse("2006-01", mux.Vars(r)["year_month"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	lastDay := time.Date(firstDay.Year(), firstDay.Month()+1, 1, 0, 0, 0, 0, firstDay.Location()).Add(-1 * time.Second)

	userPaymentAmountList, err := h.GroupTransactionsRepo.GetUserPaymentAmountList(groupID, firstDay, lastDay)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(userPaymentAmountList) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&NoContentMsg{"当月の取引履歴は見つかりませんでした。"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	groupAccountsList := model.NewGroupAccountsList(userPaymentAmountList, groupID, firstDay)

	for i := 0; i < len(userPaymentAmountList); i++ {
		userPaymentAmountList[i].PaymentAmountToUser = userPaymentAmountList[i].TotalPaymentAmount - groupAccountsList.GroupAveragePaymentAmount
	}

	dbGroupAccountsList, err := h.GroupTransactionsRepo.GetGroupAccountsList(firstDay, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(dbGroupAccountsList) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&NoContentMsg{"当月の会計データは見つかりませんでした。"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	groupAccountsList.GroupAccountsList = dbGroupAccountsList

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(groupAccountsList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PostMonthlyGroupTransactionsAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	firstDay, err := time.Parse("2006-01", mux.Vars(r)["year_month"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	lastDay := time.Date(firstDay.Year(), firstDay.Month()+1, 1, 0, 0, 0, 0, firstDay.Location()).Add(-1 * time.Second)

	userPaymentAmountList, err := h.GroupTransactionsRepo.GetUserPaymentAmountList(groupID, firstDay, lastDay)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupAccountsList := model.NewGroupAccountsList(userPaymentAmountList, groupID, firstDay)

	for i := 0; i < len(userPaymentAmountList); i++ {
		userPaymentAmountList[i].PaymentAmountToUser = userPaymentAmountList[i].TotalPaymentAmount - groupAccountsList.GroupAveragePaymentAmount
	}

	payerList := model.NewPayerList(userPaymentAmountList)
	recipientList := model.NewRecipientList(userPaymentAmountList)

	for i, payer := range payerList.PayerList {
		for j, recipient := range recipientList.RecipientList {
			if payer.PaymentAmountToUser+recipient.PaymentAmountToUser == 0 && payer.PaymentAmountToUser != 0 && recipient.PaymentAmountToUser != 0 {
				groupAccount := model.GroupAccount{
					Recipient:     recipient.UserID,
					Payer:         payer.UserID,
					PaymentAmount: recipient.PaymentAmountToUser,
				}

				groupAccountsList.GroupAccountsList = append(groupAccountsList.GroupAccountsList, groupAccount)

				payerList.PayerList[i].PaymentAmountToUser = 0
				recipientList.RecipientList[j].PaymentAmountToUser = 0
			}
		}
	}

	for i, j := 0, 0; i < len(recipientList.RecipientList) && j < len(payerList.PayerList); {
		if recipientList.RecipientList[i].PaymentAmountToUser == 0 {
			i++
			j = 0
			continue
		}

		if payerList.PayerList[j].PaymentAmountToUser == 0 {
			j++
			continue
		}

		groupAccount := model.GroupAccount{
			Recipient: recipientList.RecipientList[i].UserID,
			Payer:     payerList.PayerList[j].UserID,
		}

		remainingAmount := recipientList.RecipientList[i].PaymentAmountToUser + payerList.PayerList[j].PaymentAmountToUser

		switch {
		case remainingAmount == 0:
			groupAccount.PaymentAmount = recipientList.RecipientList[i].PaymentAmountToUser
			groupAccountsList.GroupAccountsList = append(groupAccountsList.GroupAccountsList, groupAccount)

			recipientList.RecipientList[i].PaymentAmountToUser = 0
			payerList.PayerList[j].PaymentAmountToUser = 0

			i++
			j++
		case remainingAmount < 0:
			groupAccount.PaymentAmount = recipientList.RecipientList[i].PaymentAmountToUser
			groupAccountsList.GroupAccountsList = append(groupAccountsList.GroupAccountsList, groupAccount)

			recipientList.RecipientList[i].PaymentAmountToUser = 0
			payerList.PayerList[j].PaymentAmountToUser = remainingAmount

			i++
		case remainingAmount > 0:
			groupAccount.PaymentAmount = int(math.Abs(float64(payerList.PayerList[j].PaymentAmountToUser)))
			groupAccountsList.GroupAccountsList = append(groupAccountsList.GroupAccountsList, groupAccount)

			recipientList.RecipientList[i].PaymentAmountToUser = remainingAmount
			payerList.PayerList[j].PaymentAmountToUser = 0

			j++
		}
	}

	if err := h.GroupTransactionsRepo.PostGroupAccountsList(groupAccountsList.GroupAccountsList, firstDay, groupID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroupAccountsList, err := h.GroupTransactionsRepo.GetGroupAccountsList(firstDay, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupAccountsList.GroupAccountsList = dbGroupAccountsList

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(groupAccountsList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PutMonthlyGroupTransactionsAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	var groupAccountsList model.GroupAccountsList
	if err := json.NewDecoder(r.Body).Decode(&groupAccountsList); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := h.GroupTransactionsRepo.PutGroupAccountsList(groupAccountsList.GroupAccountsList); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(groupAccountsList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) DeleteMonthlyGroupTransactionsAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := verifyGroupAffiliation(groupID, userID); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	yearMonth, err := time.Parse("2006-01", mux.Vars(r)["year_month"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	dbGroupAccountsList, err := h.GroupTransactionsRepo.GetGroupAccountsList(yearMonth, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(dbGroupAccountsList) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&NoContentMsg{"当月の会計データは見つかりませんでした。"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	if err := h.GroupTransactionsRepo.DeleteGroupAccountsList(yearMonth, groupID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&DeleteContentMsg{"グループ会計データを削除しました。"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
