package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type ShoppingItemCategoryID struct {
	MediumCategoryID model.NullInt64 `json:"medium_category_id"`
	CustomCategoryID model.NullInt64 `json:"custom_category_id"`
}

func getShoppingItemCategoryName(shoppingItem model.ShoppingItem) (model.ShoppingItem, error) {
	accountHost := os.Getenv("ACCOUNT_HOST")
	requestURL := fmt.Sprintf("http://%s:8081/categories/names", accountHost)

	shoppingItemCategoryID := ShoppingItemCategoryID{
		MediumCategoryID: shoppingItem.MediumCategoryID,
		CustomCategoryID: shoppingItem.CustomCategoryID,
	}

	requestBody, err := json.Marshal(&shoppingItemCategoryID)
	if err != nil {
		return shoppingItem, err
	}

	request, err := http.NewRequest(
		"GET",
		requestURL,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return shoppingItem, err
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
		return shoppingItem, err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	if err := json.NewDecoder(response.Body).Decode(&shoppingItem); err != nil {
		return shoppingItem, err
	}

	if response.StatusCode == http.StatusInternalServerError {
		return shoppingItem, &InternalServerErrorMsg{"500 Internal Server Error"}
	}

	return shoppingItem, nil
}

func (h *DBHandler) PostShoppingItem(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var shoppingItem model.ShoppingItem
	if err := json.NewDecoder(r.Body).Decode(&shoppingItem); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	result, err := h.ShoppingListRepo.PostShoppingItem(&shoppingItem, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	shoppingItem, err = h.ShoppingListRepo.GetShoppingItem(int(lastInsertId))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	shoppingItem, err = getShoppingItemCategoryName(shoppingItem)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&shoppingItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
