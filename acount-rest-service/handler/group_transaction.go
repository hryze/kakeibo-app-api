package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"

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

type DeleteGroupTransactionMsg struct {
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

	dbGroupTransactionsList, err := h.DBRepo.GetMonthlyGroupTransactionsList(groupID, firstDay, lastDay)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(dbGroupTransactionsList) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&NoSearchContentMsg{"条件に一致する取引履歴は見つかりませんでした。"}); err != nil {
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

	if err := validateTransaction(&groupTransactionReceiver); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	result, err := h.DBRepo.PostGroupTransaction(&groupTransactionReceiver, groupID, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroupTransactionSender, err := h.DBRepo.GetGroupTransaction(int(lastInsertId))
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

	if err := validateTransaction(&groupTransactionReceiver); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	groupTransactionID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"transaction ID を正しく指定してください。"}))
		return
	}

	if err := h.DBRepo.PutGroupTransaction(&groupTransactionReceiver, groupTransactionID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupTransactionSender, err := h.DBRepo.GetGroupTransaction(groupTransactionID)
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

	if _, err := h.DBRepo.GetGroupTransaction(groupTransactionID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"こちらのトランザクションは既に削除されています。"}))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := h.DBRepo.DeleteGroupTransaction(groupTransactionID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&DeleteGroupTransactionMsg{"トランザクションを削除しました。"}); err != nil {
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

	dbGroupTransactionsList, err := h.DBRepo.SearchGroupTransactionsList(query)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(dbGroupTransactionsList) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&NoSearchContentMsg{"条件に一致する取引履歴は見つかりませんでした。"}); err != nil {
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
