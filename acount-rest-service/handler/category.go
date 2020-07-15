package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/mux"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"

	"github.com/garyburd/redigo/redis"
)

type DeleteCustomCategoryMsg struct {
	Message string `json:"message"`
}

func validateCustomCategory(r *http.Request, customCategory *model.CustomCategory) error {
	if strings.HasPrefix(customCategory.Name, " ") || strings.HasPrefix(customCategory.Name, "　") {
		if r.Method == http.MethodPost {
			return &ValidationErrorMsg{"中カテゴリーの登録に失敗しました。 文字列先頭に空白がないか確認してください。"}
		}

		return &ValidationErrorMsg{"中カテゴリーの更新に失敗しました。 文字列先頭に空白がないか確認してください。"}
	}

	if strings.HasSuffix(customCategory.Name, " ") || strings.HasSuffix(customCategory.Name, "　") {
		if r.Method == http.MethodPost {
			return &ValidationErrorMsg{"中カテゴリーの登録に失敗しました。 文字列末尾に空白がないか確認してください。"}
		}

		return &ValidationErrorMsg{"中カテゴリーの更新に失敗しました。 文字列末尾に空白がないか確認してください。"}
	}

	if utf8.RuneCountInString(customCategory.Name) > 9 {
		return &ValidationErrorMsg{"カテゴリー名は9文字以下で入力してください。"}
	}

	return nil
}

func (h *DBHandler) GetCategoriesList(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			responseByJSON(w, r, nil, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	bigCategoriesList, err := h.DBRepo.GetBigCategoriesList()
	if err != nil {
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	mediumCategoriesList, err := h.DBRepo.GetMediumCategoriesList()
	if err != nil {
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	customCategoriesList, err := h.DBRepo.GetCustomCategoriesList(userID)
	if err != nil {
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
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
	responseByJSON(w, r, &categoriesList, nil)
}

func (h *DBHandler) PostCustomCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			responseByJSON(w, r, nil, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	customCategory := model.NewCustomCategory()
	if err := json.NewDecoder(r.Body).Decode(&customCategory); err != nil {
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	if err := validateCustomCategory(r, &customCategory); err != nil {
		responseByJSON(w, r, nil, NewHTTPError(http.StatusBadRequest, err))
		return
	}
	if err := h.DBRepo.FindCustomCategory(&customCategory, userID); err != sql.ErrNoRows {
		if err == nil {
			responseByJSON(w, r, nil, NewHTTPError(http.StatusConflict, &ConflictErrorMsg{"中カテゴリーの登録に失敗しました。 同じカテゴリー名が既に存在していないか確認してください。"}))
			return
		}
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	result, err := h.DBRepo.PostCustomCategory(&customCategory, userID)
	if err != nil {
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	customCategory.ID = int(lastInsertId)
	responseByJSON(w, r, &customCategory, nil)
}

func (h *DBHandler) PutCustomCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			responseByJSON(w, r, nil, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	customCategory := model.NewCustomCategory()
	if err := json.NewDecoder(r.Body).Decode(&customCategory); err != nil {
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	customCategory.ID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	if err := validateCustomCategory(r, &customCategory); err != nil {
		responseByJSON(w, r, nil, NewHTTPError(http.StatusBadRequest, err))
		return
	}
	if err := h.DBRepo.FindCustomCategory(&customCategory, userID); err != sql.ErrNoRows {
		if err == nil {
			responseByJSON(w, r, nil, NewHTTPError(http.StatusConflict, &ConflictErrorMsg{"中カテゴリーの更新に失敗しました。 同じカテゴリー名が既に存在していないか確認してください。"}))
			return
		}
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	if err := h.DBRepo.PutCustomCategory(&customCategory, userID); err != nil {
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	responseByJSON(w, r, &customCategory, nil)
}

func (h *DBHandler) DeleteCustomCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			responseByJSON(w, r, nil, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	customCategoryID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	if err := h.DBRepo.DeleteCustomCategory(customCategoryID, userID); err != nil {
		responseByJSON(w, r, nil, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}
	responseByJSON(w, r, &DeleteCustomCategoryMsg{"カスタムカテゴリーを削除しました。"}, nil)
}
