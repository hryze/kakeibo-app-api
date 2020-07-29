package handler

import (
	"database/sql/driver"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"

	"github.com/gorilla/mux"

	"github.com/garyburd/redigo/redis"
)

type NoContentMsg struct {
	Message string `json:"message"`
}

type TodoValidationErrorMsg struct {
	Message []string `json:"message"`
}

func (e *TodoValidationErrorMsg) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}
	return string(b)
}

func validateTodo(todo *model.Todo) error {
	var todoValidationErrorMsg TodoValidationErrorMsg

	validate := validator.New()
	validate.RegisterCustomTypeFunc(validateValuer, model.Date{})
	validate.RegisterValidation("date", dateValidation)
	validate.RegisterValidation("blank", blankValidation)
	err := validate.Struct(todo)
	if err == nil {
		return nil
	}

	for _, err := range err.(validator.ValidationErrors) {
		var errorMessage string

		fieldName := err.Field()
		switch fieldName {
		case "ImplementationDate":
			errorMessage = "todo実施日を正しく選択してください。"
		case "DueDate":
			errorMessage = "todo期限日を正しく選択してください。"
		case "TodoContent":
			tagName := err.Tag()
			switch tagName {
			case "required":
				errorMessage = "内容が入力されていません。"
			case "max":
				errorMessage = "内容は100文字以内で入力してください"
			case "blank":
				errorMessage = "内容の文字列先頭か末尾に空白がないか確認してください。"
			}
		}
		todoValidationErrorMsg.Message = append(todoValidationErrorMsg.Message, errorMessage)
	}

	return &todoValidationErrorMsg
}

func validateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}
	return nil
}

func dateValidation(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	stringDate := date.String()
	trimDate := strings.Trim(stringDate, "\"")[:10]

	minDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)

	dateTime, err := time.Parse("2006-01-02", trimDate)
	if err != nil {
		return false
	}
	if dateTime.Before(minDate) || dateTime.After(maxDate) {
		return false
	}

	return true
}

func blankValidation(fl validator.FieldLevel) bool {
	text := fl.Field().String()

	if strings.HasPrefix(text, " ") || strings.HasPrefix(text, "　") || strings.HasSuffix(text, " ") || strings.HasSuffix(text, "　") {
		return false
	}

	return true
}

func (h *DBHandler) GetDailyTodoList(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	date, err := time.Parse("2006-01-02", mux.Vars(r)["date"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}

	implementationTodoList, err := h.DBRepo.GetDailyImplementationTodoList(date, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dueTodoList, err := h.DBRepo.GetDailyDueTodoList(date, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(implementationTodoList) == 0 && len(dueTodoList) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&NoContentMsg{"今日実施予定todo、締切予定todoは登録されていません。"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	todoList := model.NewTodoList(implementationTodoList, dueTodoList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&todoList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetMonthlyTodoList(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	firstDay, err := time.Parse("2006-01", mux.Vars(r)["year_month"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"年月を正しく指定してください。"}))
		return
	}
	lastDay := time.Date(firstDay.Year(), firstDay.Month()+1, 1, 0, 0, 0, 0, firstDay.Location()).Add(-1 * time.Second)

	implementationTodoList, err := h.DBRepo.GetMonthlyImplementationTodoList(firstDay, lastDay, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dueTodoList, err := h.DBRepo.GetMonthlyDueTodoList(firstDay, lastDay, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(implementationTodoList) == 0 && len(dueTodoList) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&NoContentMsg{"当月実施予定todoは登録されていません。"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	todoList := model.NewTodoList(implementationTodoList, dueTodoList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&todoList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PostTodo(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var todo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := validateTodo(&todo); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	result, err := h.DBRepo.PostTodo(&todo, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbTodo, err := h.DBRepo.GetTodo(int(lastInsertId))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(dbTodo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
