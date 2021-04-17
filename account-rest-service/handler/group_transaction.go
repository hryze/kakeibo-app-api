package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/config"
	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
)

type GroupTransactionsSearchQuery struct {
	TransactionType   string
	BigCategoryID     string
	Shop              string
	Memo              string
	LowAmount         string
	HighAmount        string
	StartDate         string
	EndDate           string
	Sort              string
	SortType          string
	Limit             string
	GroupID           string
	PaymentUserIDList []string
}

type GroupTransactionProcessLockErrorMsg struct {
	Message string `json:"message"`
}

func NewGroupTransactionsSearchQuery(urlQuery url.Values, groupID string) GroupTransactionsSearchQuery {
	startDate := trimDate(urlQuery.Get("start_date"))
	endDate := trimDate(urlQuery.Get("end_date"))

	return GroupTransactionsSearchQuery{
		TransactionType:   urlQuery.Get("transaction_type"),
		BigCategoryID:     urlQuery.Get("big_category_id"),
		Shop:              urlQuery.Get("shop"),
		Memo:              urlQuery.Get("memo"),
		LowAmount:         urlQuery.Get("low_amount"),
		HighAmount:        urlQuery.Get("high_amount"),
		StartDate:         startDate,
		EndDate:           endDate,
		Sort:              urlQuery.Get("sort"),
		SortType:          urlQuery.Get("sort_type"),
		Limit:             urlQuery.Get("limit"),
		GroupID:           groupID,
		PaymentUserIDList: urlQuery["payment_user_id"],
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
            group_transactions.posted_date posted_date,
            group_transactions.updated_date updated_date,
            group_transactions.transaction_date transaction_date,
            group_transactions.shop shop,
            group_transactions.memo memo,
            group_transactions.amount amount,
            group_transactions.posted_user_id posted_user_id,
            group_transactions.updated_user_id updated_user_id,
            group_transactions.payment_user_id payment_user_id,
            group_transactions.big_category_id big_category_id,
            big_categories.category_name big_category_name,
            group_transactions.medium_category_id medium_category_id,
            medium_categories.category_name medium_category_name,
            group_transactions.custom_category_id custom_category_id,
            group_custom_categories.category_name custom_category_name
        FROM
            group_transactions
        INNER JOIN
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

        {{ if eq (len .PaymentUserIDList) 1 }}
        {{ range $i, $PaymentUserID := .PaymentUserIDList }}
        AND
            group_transactions.payment_user_id = "{{ $PaymentUserID }}"
        {{ end }}
        {{ end }}

        {{ if gt (len .PaymentUserIDList) 1 }}
        {{ range $i, $PaymentUserID := .PaymentUserIDList }}
        {{ if eq $i 0}}
        AND
            group_transactions.payment_user_id IN("{{ $PaymentUserID }}"
        {{ end }}
        {{ if gt $i 0 }}
        ,"{{ $PaymentUserID }}"
        {{ end }}
        {{ end }}
        {{ end }}
        {{ if gt (len .PaymentUserIDList) 1 }}
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

func getGroupUserIDList(groupID int) ([]string, error) {
	requestURL := fmt.Sprintf(
		"http://%s:%d/groups/%d/users",
		config.Env.UserApi.Host, config.Env.UserApi.Port, groupID,
	)

	request, err := http.NewRequest(
		"GET",
		requestURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          500,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 60 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	var groupUserIDList []string
	if err := json.NewDecoder(response.Body).Decode(&groupUserIDList); err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusBadRequest {
		return nil, &BadRequestErrorMsg{"指定されたグループには、ユーザーは所属していません。"}
	}

	if response.StatusCode == http.StatusInternalServerError {
		return nil, &InternalServerErrorMsg{"500 Internal Server Error"}
	}

	return groupUserIDList, nil
}

func paymentAmountSplitBill(groupAccountsList *model.GroupAccountsList, payerList model.PayerList, recipientList model.RecipientList, groupID int, month time.Time) {
	for i, payer := range payerList.PayerList {
		for j, recipient := range recipientList.RecipientList {
			if payer.PaymentAmountToUser+recipient.PaymentAmountToUser == 0 && payer.PaymentAmountToUser != 0 && recipient.PaymentAmountToUser != 0 {
				groupAccount := model.GroupAccount{
					GroupID:       groupID,
					Month:         month,
					Recipient:     model.NullString{NullString: sql.NullString{String: recipient.UserID, Valid: true}},
					Payer:         model.NullString{NullString: sql.NullString{String: payer.UserID, Valid: true}},
					PaymentAmount: model.NullInt{Int: recipient.PaymentAmountToUser, Valid: true},
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
			GroupID:   groupID,
			Month:     month,
			Recipient: model.NullString{NullString: sql.NullString{String: recipientList.RecipientList[i].UserID, Valid: true}},
			Payer:     model.NullString{NullString: sql.NullString{String: payerList.PayerList[j].UserID, Valid: true}},
		}

		remainingAmount := recipientList.RecipientList[i].PaymentAmountToUser + payerList.PayerList[j].PaymentAmountToUser

		switch {
		case remainingAmount == 0:
			groupAccount.PaymentAmount.Int = recipientList.RecipientList[i].PaymentAmountToUser
			groupAccount.PaymentAmount.Valid = true
			groupAccountsList.GroupAccountsList = append(groupAccountsList.GroupAccountsList, groupAccount)

			recipientList.RecipientList[i].PaymentAmountToUser = 0
			payerList.PayerList[j].PaymentAmountToUser = 0

			i++
			j++
		case remainingAmount < 0:
			groupAccount.PaymentAmount.Int = recipientList.RecipientList[i].PaymentAmountToUser
			groupAccount.PaymentAmount.Valid = true
			groupAccountsList.GroupAccountsList = append(groupAccountsList.GroupAccountsList, groupAccount)

			recipientList.RecipientList[i].PaymentAmountToUser = 0
			payerList.PayerList[j].PaymentAmountToUser = remainingAmount

			i++
		case remainingAmount > 0:
			groupAccount.PaymentAmount.Int = int(math.Abs(float64(payerList.PayerList[j].PaymentAmountToUser)))
			groupAccount.PaymentAmount.Valid = true
			groupAccountsList.GroupAccountsList = append(groupAccountsList.GroupAccountsList, groupAccount)

			recipientList.RecipientList[i].PaymentAmountToUser = remainingAmount
			payerList.PayerList[j].PaymentAmountToUser = 0

			j++
		}
	}
}

func generateGroupAccountsListByPayer(groupAccountsList *model.GroupAccountsList) {
L:
	for _, groupAccount := range groupAccountsList.GroupAccountsList {
		for i, groupAccountsListByPayer := range groupAccountsList.GroupAccountsListByPayersList {
			if groupAccount.Payer.String == groupAccountsListByPayer.Payer.String {
				groupAccountsList.GroupAccountsListByPayersList[i].GroupAccountsList = append(groupAccountsList.GroupAccountsListByPayersList[i].GroupAccountsList, groupAccount)

				continue L
			}
		}

		groupAccountsList.GroupAccountsListByPayersList = append(groupAccountsList.GroupAccountsListByPayersList, model.GroupAccountsListByPayer{
			Payer:             groupAccount.Payer,
			GroupAccountsList: []model.GroupAccount{groupAccount},
		})
	}
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

func (h *DBHandler) Get10LatestGroupTransactionsList(w http.ResponseWriter, r *http.Request) {
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

	latestGroupTransactionsList, err := h.GroupTransactionsRepo.Get10LatestGroupTransactionsList(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(latestGroupTransactionsList.GroupTransactionsList) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&NoContentMsg{"取引履歴がありません。"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&latestGroupTransactionsList); err != nil {
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

	// Check if the transaction date of the json request transaction has been settled.
	yearMonth := time.Date(groupTransactionReceiver.TransactionDate.Time.Year(), groupTransactionReceiver.TransactionDate.Time.Month(), 1, 0, 0, 0, 0, time.UTC)

	dbGroupAccountsList, err := h.GroupTransactionsRepo.GetGroupAccountsList(yearMonth, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	} else if len(dbGroupAccountsList) != 0 {
		message := fmt.Sprintf("%d年%d月の取引は精算済みのため追加できません。", yearMonth.Year(), yearMonth.Month())
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &GroupTransactionProcessLockErrorMsg{Message: message}))
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

	groupTransactionID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"transaction ID を正しく指定してください。"}))
		return
	}

	dbGroupTransaction, err := h.GroupTransactionsRepo.GetGroupTransaction(groupTransactionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusNotFound, &NotFoundErrorMsg{"該当する取引が見つかりませんでした。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	// Check if the transaction date of the transaction retrieved from the Database has been settled.
	yearMonth := time.Date(dbGroupTransaction.TransactionDate.Time.Year(), dbGroupTransaction.TransactionDate.Time.Month(), 1, 0, 0, 0, 0, time.UTC)

	dbGroupAccountsList, err := h.GroupTransactionsRepo.GetGroupAccountsList(yearMonth, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	} else if len(dbGroupAccountsList) != 0 {
		message := fmt.Sprintf("%d年%d月の取引は精算済みのため更新できません。", yearMonth.Year(), yearMonth.Month())
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &GroupTransactionProcessLockErrorMsg{Message: message}))
		return
	}

	var groupTransactionReceiver model.GroupTransactionReceiver
	if err := json.NewDecoder(r.Body).Decode(&groupTransactionReceiver); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	// Check if the transaction date of the json request transaction has been settled.
	yearMonth = time.Date(groupTransactionReceiver.TransactionDate.Time.Year(), groupTransactionReceiver.TransactionDate.Time.Month(), 1, 0, 0, 0, 0, time.UTC)

	dbGroupAccountsList, err = h.GroupTransactionsRepo.GetGroupAccountsList(yearMonth, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	} else if len(dbGroupAccountsList) != 0 {
		message := fmt.Sprintf("%d年%d月の取引は精算済みのため更新できません。", yearMonth.Year(), yearMonth.Month())
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &GroupTransactionProcessLockErrorMsg{Message: message}))
		return
	}

	if err := validateTransaction(&groupTransactionReceiver); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	if err := h.GroupTransactionsRepo.PutGroupTransaction(&groupTransactionReceiver, groupTransactionID, userID); err != nil {
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

	dbGroupTransaction, err := h.GroupTransactionsRepo.GetGroupTransaction(groupTransactionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusNotFound, &NotFoundErrorMsg{"該当する取引が見つかりませんでした。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	// Check if the transaction date of the transaction retrieved from the Database has been settled.
	yearMonth := time.Date(dbGroupTransaction.TransactionDate.Time.Year(), dbGroupTransaction.TransactionDate.Time.Month(), 1, 0, 0, 0, 0, time.UTC)

	dbGroupAccountsList, err := h.GroupTransactionsRepo.GetGroupAccountsList(yearMonth, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	} else if len(dbGroupAccountsList) != 0 {
		message := fmt.Sprintf("%d年%d月の取引は精算済みのため削除できません。", yearMonth.Year(), yearMonth.Month())
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &GroupTransactionProcessLockErrorMsg{Message: message}))
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

func (h *DBHandler) GetGroupShoppingItemRelatedTransactionDataList(w http.ResponseWriter, r *http.Request) {
	var transactionIdList []int
	if err := json.NewDecoder(r.Body).Decode(&transactionIdList); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupTransactionsList, err := h.GroupTransactionsRepo.GetGroupShoppingItemRelatedTransactionDataList(transactionIdList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupTransactionsList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetYearlyAccountingStatus(w http.ResponseWriter, r *http.Request) {
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

	firstDayOfYear, err := time.Parse("2006", mux.Vars(r)["year"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年を正しく指定してください。"}))
		return
	}

	transactionExistenceByMonths, err := h.GroupTransactionsRepo.YearlyGroupTransactionExistenceConfirmation(firstDayOfYear, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	yearlyGroupAccountsList, err := h.GroupTransactionsRepo.GetYearlyGroupAccountsList(firstDayOfYear, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	yearlyAccountingStatus := model.NewYearlyAccountingStatus(firstDayOfYear, userID, transactionExistenceByMonths, yearlyGroupAccountsList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&yearlyAccountingStatus); err != nil {
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

	groupUserIDList, err := getGroupUserIDList(groupID)
	if err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	userPaymentAmountList, err := h.GroupTransactionsRepo.GetUserPaymentAmountList(groupID, groupUserIDList, firstDay, lastDay)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var isNotZero bool
	for _, userPaymentAmount := range userPaymentAmountList {
		if userPaymentAmount.TotalPaymentAmount > 0 {
			isNotZero = true
			break
		}
	}

	if !isNotZero {
		errorResponseByJSON(w, NewHTTPError(http.StatusNotFound, &NotFoundErrorMsg{"当月の取引履歴が見つかりませんでした。"}))
		return
	}

	groupAccountsList := model.NewGroupAccountsList(userPaymentAmountList, groupID, firstDay)

	for i := 0; i < len(userPaymentAmountList); i++ {
		userPaymentAmountList[i].PaymentAmountToUser = userPaymentAmountList[i].TotalPaymentAmount - groupAccountsList.GroupAveragePaymentAmount
	}

	groupAccountsList.GroupAccountsList, err = h.GroupTransactionsRepo.GetGroupAccountsList(firstDay, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(groupAccountsList.GroupAccountsList) == 0 {
		errorResponseByJSON(w, NewHTTPError(http.StatusNotFound, &NotFoundErrorMsg{"当月は未会計です。"}))
		return
	}

	generateGroupAccountsListByPayer(&groupAccountsList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupAccountsList); err != nil {
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

	groupUserIDList, err := getGroupUserIDList(groupID)
	if err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	if len(groupUserIDList) == 1 {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"グループ人数が1人のため会計できません。"}))
		return
	}

	dbGroupAccountsList, err := h.GroupTransactionsRepo.GetGroupAccountsList(firstDay, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(dbGroupAccountsList) >= 1 {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"当月は会計済です。"}))
		return
	}

	userPaymentAmountList, err := h.GroupTransactionsRepo.GetUserPaymentAmountList(groupID, groupUserIDList, firstDay, lastDay)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var isNotZero bool
	for _, userPaymentAmount := range userPaymentAmountList {
		if userPaymentAmount.TotalPaymentAmount > 0 {
			isNotZero = true
			break
		}
	}

	if !isNotZero {
		errorResponseByJSON(w, NewHTTPError(http.StatusNotFound, &NotFoundErrorMsg{"当月の取引履歴が見つかりませんでした。"}))
		return
	}

	groupAccountsList := model.NewGroupAccountsList(userPaymentAmountList, groupID, firstDay)

	for i := 0; i < len(userPaymentAmountList); i++ {
		userPaymentAmountList[i].PaymentAmountToUser = userPaymentAmountList[i].TotalPaymentAmount - groupAccountsList.GroupAveragePaymentAmount
	}

	payerList := model.NewPayerList(userPaymentAmountList)
	recipientList := model.NewRecipientList(userPaymentAmountList)

	if len(payerList.PayerList) == 0 && len(recipientList.RecipientList) == 0 {
		groupAccountsList.GroupAccountsList = append(groupAccountsList.GroupAccountsList, model.GroupAccount{
			GroupID:             groupID,
			Month:               firstDay,
			PaymentConfirmation: true,
			ReceiptConfirmation: true,
		})
	} else if len(payerList.PayerList) != 0 && len(recipientList.RecipientList) != 0 {
		paymentAmountSplitBill(&groupAccountsList, payerList, recipientList, groupID, firstDay)
	}

	if err := h.GroupTransactionsRepo.PostGroupAccountsList(groupAccountsList.GroupAccountsList); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupAccountsList.GroupAccountsList, err = h.GroupTransactionsRepo.GetGroupAccountsList(firstDay, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	generateGroupAccountsListByPayer(&groupAccountsList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&groupAccountsList); err != nil {
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
		errorResponseByJSON(w, NewHTTPError(http.StatusNotFound, &NotFoundErrorMsg{"当月は未会計です。"}))
		return
	}

	groupAccountID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"会計ID を正しく指定してください。"}))
		return
	}

	var groupAccount model.GroupAccount
	if err := json.NewDecoder(r.Body).Decode(&groupAccount); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := h.GroupTransactionsRepo.PutGroupAccount(groupAccount, groupAccountID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupAccount); err != nil {
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
		errorResponseByJSON(w, NewHTTPError(http.StatusNotFound, &NotFoundErrorMsg{"当月は未会計です。"}))
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
