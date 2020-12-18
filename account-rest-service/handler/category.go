package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/mux"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"

	"github.com/garyburd/redigo/redis"
)

type CustomCategoryValidationErrorMsg struct {
	Message string `json:"message"`
}

func (e *CustomCategoryValidationErrorMsg) Error() string {
	return e.Message
}

func validateCustomCategory(r *http.Request, customCategory *model.CustomCategory) error {
	if strings.HasPrefix(customCategory.Name, " ") || strings.HasPrefix(customCategory.Name, "　") {
		if r.Method == http.MethodPost {
			return &CustomCategoryValidationErrorMsg{"中カテゴリーの登録に失敗しました。 文字列先頭に空白がないか確認してください。"}
		}

		return &CustomCategoryValidationErrorMsg{"中カテゴリーの更新に失敗しました。 文字列先頭に空白がないか確認してください。"}
	}

	if strings.HasSuffix(customCategory.Name, " ") || strings.HasSuffix(customCategory.Name, "　") {
		if r.Method == http.MethodPost {
			return &CustomCategoryValidationErrorMsg{"中カテゴリーの登録に失敗しました。 文字列末尾に空白がないか確認してください。"}
		}

		return &CustomCategoryValidationErrorMsg{"中カテゴリーの更新に失敗しました。 文字列末尾に空白がないか確認してください。"}
	}

	if utf8.RuneCountInString(customCategory.Name) > 9 {
		return &CustomCategoryValidationErrorMsg{"カテゴリー名は9文字以下で入力してください。"}
	}

	return nil
}

func (h *DBHandler) GetCategoriesList(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	bigCategoriesList, err := h.CategoriesRepo.GetBigCategoriesList()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	mediumCategoriesList, err := h.CategoriesRepo.GetMediumCategoriesList()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	customCategoriesList, err := h.CategoriesRepo.GetCustomCategoriesList(userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for i, bigCategory := range bigCategoriesList {
		for _, customCategory := range customCategoriesList {
			if bigCategory.TransactionType == "income" && bigCategory.ID == customCategory.BigCategoryID {
				bigCategoriesList[i].IncomeAssociatedCategoriesList = append(bigCategoriesList[i].IncomeAssociatedCategoriesList, customCategory)
			} else if bigCategory.TransactionType == "expense" && bigCategory.ID == customCategory.BigCategoryID {
				bigCategoriesList[i].ExpenseAssociatedCategoriesList = append(bigCategoriesList[i].ExpenseAssociatedCategoriesList, customCategory)
			}
		}
	}

	for i, bigCategory := range bigCategoriesList {
		for _, mediumCategory := range mediumCategoriesList {
			if bigCategory.TransactionType == "income" && bigCategory.ID == mediumCategory.BigCategoryID {
				bigCategoriesList[i].IncomeAssociatedCategoriesList = append(bigCategoriesList[i].IncomeAssociatedCategoriesList, mediumCategory)
			} else if bigCategory.TransactionType == "expense" && bigCategory.ID == mediumCategory.BigCategoryID {
				bigCategoriesList[i].ExpenseAssociatedCategoriesList = append(bigCategoriesList[i].ExpenseAssociatedCategoriesList, mediumCategory)
			}
		}
	}

	var categoriesList model.CategoriesList
	for _, bigCategory := range bigCategoriesList {
		if bigCategory.TransactionType == "income" {
			categoriesList.IncomeBigCategoriesList = append(categoriesList.IncomeBigCategoriesList, model.NewIncomeBigCategory(&bigCategory))
		} else if bigCategory.TransactionType == "expense" {
			categoriesList.ExpenseBigCategoriesList = append(categoriesList.ExpenseBigCategoriesList, model.NewExpenseBigCategory(&bigCategory))
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&categoriesList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PostCustomCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	customCategory := model.NewCustomCategory()
	if err := json.NewDecoder(r.Body).Decode(&customCategory); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := validateCustomCategory(r, &customCategory); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	if err := h.CategoriesRepo.FindCustomCategory(&customCategory, userID); err != sql.ErrNoRows {
		if err == nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusConflict, &ConflictErrorMsg{"中カテゴリーの登録に失敗しました。 同じカテゴリー名が既に存在していないか確認してください。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	result, err := h.CategoriesRepo.PostCustomCategory(&customCategory, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	customCategory.ID = int(lastInsertId)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&customCategory); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PutCustomCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	customCategory := model.NewCustomCategory()
	if err := json.NewDecoder(r.Body).Decode(&customCategory); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	customCategory.ID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"custom category ID を正しく指定してください。"}))
		return
	}

	if err := validateCustomCategory(r, &customCategory); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	if err := h.CategoriesRepo.FindCustomCategory(&customCategory, userID); err != sql.ErrNoRows {
		if err == nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusConflict, &ConflictErrorMsg{"中カテゴリーの更新に失敗しました。 同じカテゴリー名が既に存在していないか確認してください。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := h.CategoriesRepo.PutCustomCategory(&customCategory); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&customCategory); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) DeleteCustomCategory(w http.ResponseWriter, r *http.Request) {
	_, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	customCategoryID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"カスタムカテゴリーID を正しく指定してください。"}))
		return
	}

	bigCategoryID, err := h.CategoriesRepo.GetBigCategoryID(customCategoryID)
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

	if err := h.CategoriesRepo.DeleteCustomCategory(customCategoryID, mediumCategoryID); err != nil {
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

func (h *DBHandler) GetCategoriesName(w http.ResponseWriter, r *http.Request) {
	var categoriesID model.CategoriesID
	if err := json.NewDecoder(r.Body).Decode(&categoriesID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	categoriesName, err := h.CategoriesRepo.GetCategoriesName(categoriesID)
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

func (h *DBHandler) GetCategoriesNameList(w http.ResponseWriter, r *http.Request) {
	var categoriesIDList []model.CategoriesID
	if err := json.NewDecoder(r.Body).Decode(&categoriesIDList); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	categoriesNameList, err := h.CategoriesRepo.GetCategoriesNameList(categoriesIDList)
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
