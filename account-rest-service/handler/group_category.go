package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/config"
	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
)

type GroupCustomCategoryValidationErrorMsg struct {
	Message string `json:"message"`
}

func (e *GroupCustomCategoryValidationErrorMsg) Error() string {
	return e.Message
}

func validateGroupCustomCategory(r *http.Request, groupCustomCategory *model.GroupCustomCategory) error {
	if strings.HasPrefix(groupCustomCategory.Name, " ") || strings.HasPrefix(groupCustomCategory.Name, "　") {
		if r.Method == http.MethodPost {
			return &GroupCustomCategoryValidationErrorMsg{"中カテゴリーの登録に失敗しました。 文字列先頭に空白がないか確認してください。"}
		}

		return &GroupCustomCategoryValidationErrorMsg{"中カテゴリーの更新に失敗しました。 文字列先頭に空白がないか確認してください。"}
	}

	if strings.HasSuffix(groupCustomCategory.Name, " ") || strings.HasSuffix(groupCustomCategory.Name, "　") {
		if r.Method == http.MethodPost {
			return &GroupCustomCategoryValidationErrorMsg{"中カテゴリーの登録に失敗しました。 文字列末尾に空白がないか確認してください。"}
		}

		return &GroupCustomCategoryValidationErrorMsg{"中カテゴリーの更新に失敗しました。 文字列末尾に空白がないか確認してください。"}
	}

	if utf8.RuneCountInString(groupCustomCategory.Name) > 9 {
		return &GroupCustomCategoryValidationErrorMsg{"カテゴリー名は9文字以下で入力してください。"}
	}

	return nil
}

func putGroupShoppingListCustomCategoryIdToMediumCategoryId(mediumCategoryID int, customCategoryID int, groupID int) error {
	requestURL := fmt.Sprintf(
		"http://%s:%d/groups/%d/shopping-list/categories",
		config.Env.TodoApi.Host, config.Env.TodoApi.Port, groupID,
	)

	categoriesID := struct {
		MediumCategoryID int `json:"medium_category_id"`
		CustomCategoryID int `json:"custom_category_id"`
	}{
		MediumCategoryID: mediumCategoryID,
		CustomCategoryID: customCategoryID,
	}

	requestBody, err := json.Marshal(&categoriesID)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		"PUT",
		requestURL,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return err
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
		return err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	if response.StatusCode == http.StatusOK {
		return nil
	}

	return errors.New("couldn't replace customCategoryID in shopping list with default mediumCategoryID")
}

func (h *DBHandler) GetGroupCategoriesList(w http.ResponseWriter, r *http.Request) {
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

	groupBigCategoriesList, err := h.GroupCategoriesRepo.GetGroupBigCategoriesList()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupMediumCategoriesList, err := h.GroupCategoriesRepo.GetGroupMediumCategoriesList()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupCustomCategoriesList, err := h.GroupCategoriesRepo.GetGroupCustomCategoriesList(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for i, groupBigCategory := range groupBigCategoriesList {
		for _, groupCustomCategory := range groupCustomCategoriesList {
			if groupBigCategory.TransactionType == "income" && groupBigCategory.ID == groupCustomCategory.BigCategoryID {
				groupBigCategoriesList[i].IncomeAssociatedCategoriesList = append(groupBigCategoriesList[i].IncomeAssociatedCategoriesList, groupCustomCategory)
			} else if groupBigCategory.TransactionType == "expense" && groupBigCategory.ID == groupCustomCategory.BigCategoryID {
				groupBigCategoriesList[i].ExpenseAssociatedCategoriesList = append(groupBigCategoriesList[i].ExpenseAssociatedCategoriesList, groupCustomCategory)
			}
		}
	}

	for i, groupBigCategory := range groupBigCategoriesList {
		for _, groupMediumCategory := range groupMediumCategoriesList {
			if groupBigCategory.TransactionType == "income" && groupBigCategory.ID == groupMediumCategory.BigCategoryID {
				groupBigCategoriesList[i].IncomeAssociatedCategoriesList = append(groupBigCategoriesList[i].IncomeAssociatedCategoriesList, groupMediumCategory)
			} else if groupBigCategory.TransactionType == "expense" && groupBigCategory.ID == groupMediumCategory.BigCategoryID {
				groupBigCategoriesList[i].ExpenseAssociatedCategoriesList = append(groupBigCategoriesList[i].ExpenseAssociatedCategoriesList, groupMediumCategory)
			}
		}
	}

	var groupCategoriesList model.GroupCategoriesList
	for _, groupBigCategory := range groupBigCategoriesList {
		if groupBigCategory.TransactionType == "income" {
			groupCategoriesList.GroupIncomeBigCategoriesList = append(groupCategoriesList.GroupIncomeBigCategoriesList, model.NewGroupIncomeBigCategory(&groupBigCategory))
		} else if groupBigCategory.TransactionType == "expense" {
			groupCategoriesList.GroupExpenseBigCategoriesList = append(groupCategoriesList.GroupExpenseBigCategoriesList, model.NewGroupExpenseBigCategory(&groupBigCategory))
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupCategoriesList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PostGroupCustomCategory(w http.ResponseWriter, r *http.Request) {
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

	groupCustomCategory := model.NewGroupCustomCategory()
	if err := json.NewDecoder(r.Body).Decode(&groupCustomCategory); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := validateGroupCustomCategory(r, &groupCustomCategory); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	if err := h.GroupCategoriesRepo.FindGroupCustomCategory(&groupCustomCategory, groupID); err != sql.ErrNoRows {
		if err == nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusConflict, &ConflictErrorMsg{"中カテゴリーの登録に失敗しました。 同じカテゴリー名が既に存在していないか確認してください。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	result, err := h.GroupCategoriesRepo.PostGroupCustomCategory(&groupCustomCategory, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupCustomCategory.ID = int(lastInsertId)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&groupCustomCategory); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PutGroupCustomCategory(w http.ResponseWriter, r *http.Request) {
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

	groupCustomCategory := model.NewGroupCustomCategory()
	if err := json.NewDecoder(r.Body).Decode(&groupCustomCategory); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupCustomCategory.ID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"custom category ID を正しく指定してください。"}))
		return
	}

	if err := validateGroupCustomCategory(r, &groupCustomCategory); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	if err := h.GroupCategoriesRepo.FindGroupCustomCategory(&groupCustomCategory, groupID); err != sql.ErrNoRows {
		if err == nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusConflict, &ConflictErrorMsg{"中カテゴリーの更新に失敗しました。 同じカテゴリー名が既に存在していないか確認してください。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := h.GroupCategoriesRepo.PutGroupCustomCategory(&groupCustomCategory); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupCustomCategory); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) DeleteGroupCustomCategory(w http.ResponseWriter, r *http.Request) {
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

	groupCustomCategoryID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"custom category ID を正しく指定してください。"}))
		return
	}

	if err := h.GroupCategoriesRepo.FindGroupCustomCategoryID(groupCustomCategoryID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"指定されたカスタムカテゴリーは既に削除されています。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	bigCategoryID, err := h.GroupCategoriesRepo.GetBigCategoryID(groupCustomCategoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusNotFound, &NotFoundErrorMsg{"カスタムカテゴリーに関連する大カテゴリーが見つかりませんでした。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	mediumCategoryID, err := replaceMediumCategoryID(bigCategoryID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusNotFound, err))
		return
	}

	if err := h.GroupCategoriesRepo.DeleteGroupCustomCategory(groupCustomCategoryID, mediumCategoryID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := putGroupShoppingListCustomCategoryIdToMediumCategoryId(mediumCategoryID, groupCustomCategoryID, groupID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&DeleteContentMsg{"カスタムカテゴリーを削除しました。"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetGroupCategoriesName(w http.ResponseWriter, r *http.Request) {
	var categoriesID model.CategoriesID
	if err := json.NewDecoder(r.Body).Decode(&categoriesID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	categoriesName, err := h.GroupCategoriesRepo.GetGroupCategoriesName(categoriesID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&categoriesName); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetGroupCategoriesNameList(w http.ResponseWriter, r *http.Request) {
	var categoriesIDList []model.CategoriesID
	if err := json.NewDecoder(r.Body).Decode(&categoriesIDList); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	categoriesNameList, err := h.GroupCategoriesRepo.GetGroupCategoriesNameList(categoriesIDList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&categoriesNameList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
