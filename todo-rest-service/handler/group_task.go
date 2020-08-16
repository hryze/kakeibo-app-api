package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

type GroupTasksUserConflictErrorMsg struct {
	Message string `json:"message"`
}

type GroupTasksUserBadRequestErrorMsg struct {
	Message string `json:"message"`
}

type GroupTaskNameBadRequestErrorMsg struct {
	Message string `json:"message"`
}

func (e *GroupTasksUserConflictErrorMsg) Error() string {
	return e.Message
}

func (e *GroupTasksUserBadRequestErrorMsg) Error() string {
	return e.Message
}

func (e *GroupTaskNameBadRequestErrorMsg) Error() string {
	return e.Message
}

func validateGroupTasksUser(groupTasksUser *model.GroupTasksUser) error {
	if strings.HasPrefix(groupTasksUser.UserID, " ") || strings.HasPrefix(groupTasksUser.UserID, "　") {
		return &GroupTasksUserBadRequestErrorMsg{"ユーザーIDの文字列先頭に空白がないか確認してください。"}
	}

	if strings.HasSuffix(groupTasksUser.UserID, " ") || strings.HasSuffix(groupTasksUser.UserID, "　") {
		return &GroupTasksUserBadRequestErrorMsg{"ユーザーIDの文字列末尾に空白がないか確認してください。"}
	}

	if len(groupTasksUser.UserID) == 0 || len(groupTasksUser.UserID) > 10 {
		return &GroupTasksUserBadRequestErrorMsg{"ユーザーIDは1文字以上10文字以内で入力してください。"}
	}

	return nil
}

func validateGroupTaskName(groupTask *model.GroupTask) error {
	if strings.HasPrefix(groupTask.TaskName, " ") || strings.HasPrefix(groupTask.TaskName, "　") {
		return &GroupTaskNameBadRequestErrorMsg{"タスク名の文字列先頭に空白がないか確認してください。"}
	}

	if strings.HasSuffix(groupTask.TaskName, " ") || strings.HasSuffix(groupTask.TaskName, "　") {
		return &GroupTaskNameBadRequestErrorMsg{"タスク名の文字列末尾に空白がないか確認してください。"}
	}

	if len(groupTask.TaskName) == 0 || len(groupTask.TaskName) > 20 {
		return &GroupTaskNameBadRequestErrorMsg{"タスク名は1文字以上20文字以内で入力してください。"}
	}

	return nil
}

func (h *DBHandler) PostGroupTasksUser(w http.ResponseWriter, r *http.Request) {
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

	var groupTasksUser model.GroupTasksUser
	if err := json.NewDecoder(r.Body).Decode(&groupTasksUser); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := validateGroupTasksUser(&groupTasksUser); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	if err := verifyGroupAffiliation(groupID, groupTasksUser.UserID); err != nil {
		_, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &GroupTasksUserBadRequestErrorMsg{"こちらのグループには、指定されたユーザーは所属していません。"}))
		return
	}

	if _, err := h.DBRepo.GetGroupTasksUser(groupTasksUser, groupID); err != sql.ErrNoRows {
		if err == nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusConflict, &GroupTasksUserConflictErrorMsg{"こちらのユーザーは、既にタスクメンバーに追加されています。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if _, err := h.DBRepo.PostGroupTasksUser(groupTasksUser, groupID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroupTasksUser, err := h.DBRepo.GetGroupTasksUser(groupTasksUser, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(dbGroupTasksUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PostGroupTask(w http.ResponseWriter, r *http.Request) {
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

	var groupTask model.GroupTask
	if err := json.NewDecoder(r.Body).Decode(&groupTask); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := validateGroupTaskName(&groupTask); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	result, err := h.DBRepo.PostGroupTask(groupTask, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroupTask, err := h.DBRepo.GetGroupTask(int(lastInsertId))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(dbGroupTask); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
