package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/garyburd/redigo/redis"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type NoContentMsg struct {
	Message string `json:"message"`
}

type GroupUserConflictErrorMsg struct {
	Message string `json:"message"`
}

type UserIDValidationErrorMsg struct {
	Message string `json:"message"`
}

func (e *GroupUserConflictErrorMsg) Error() string {
	return e.Message
}

func (e *UserIDValidationErrorMsg) Error() string {
	return e.Message
}

func postInitGroupStandardBudgets(groupID int) error {
	request, err := http.NewRequest(
		"POST",
		"http://localhost:8081/groups/standard-budgets",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"group_id":%d}`, groupID))),
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusCreated {
		return nil
	}

	return errors.New("couldn't create a group standard budget")
}

func checkForUniqueGroupUser(h *DBHandler, groupID int, userID string) error {
	if err := h.DBRepo.FindGroupUser(groupID, userID); err != sql.ErrNoRows {
		if err == nil {
			return &GroupUserConflictErrorMsg{"こちらのユーザーは既にグループに参加しています。"}
		}

		return err
	}

	if err := h.DBRepo.FindGroupUnapprovedUser(groupID, userID); err != sql.ErrNoRows {
		if err == nil {
			return &GroupUserConflictErrorMsg{"こちらのユーザーは既にグループに招待しています。"}
		}

		return err
	}

	return nil
}

func validateUserID(userID string) error {
	if strings.HasPrefix(userID, " ") || strings.HasPrefix(userID, "　") {
		return &UserIDValidationErrorMsg{"文字列先頭に空白がないか確認してください。"}
	}

	if strings.HasSuffix(userID, " ") || strings.HasSuffix(userID, "　") {
		return &UserIDValidationErrorMsg{"文字列末尾に空白がないか確認してください。"}
	}

	if len(userID) == 0 || len(userID) > 10 {
		return &UserIDValidationErrorMsg{"ユーザーIDは1文字以上、10文字以内で入力してください。"}
	}

	return nil
}

func (h *DBHandler) GetGroupList(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupList, err := h.DBRepo.GetGroupList(userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if len(groupList) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&NoContentMsg{"参加しているグループはありません。"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	groupUsersList, err := h.DBRepo.GetGroupUsersList(groupList)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupUnapprovedUsersList, err := h.DBRepo.GetGroupUnapprovedUsersList(groupList)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for i := 0; i < len(groupList); i++ {
		for _, groupUser := range groupUsersList {
			if groupList[i].GroupID == groupUser.GroupID {
				groupList[i].GroupUsersList = append(groupList[i].GroupUsersList, groupUser)
			}
		}

		for _, groupUnapprovedUser := range groupUnapprovedUsersList {
			if groupList[i].GroupID == groupUnapprovedUser.GroupID {
				groupList[i].GroupUnapprovedUsersList = append(groupList[i].GroupUnapprovedUsersList, groupUnapprovedUser)
			}
		}
	}

	groupListSender := model.NewGroupList(groupList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupListSender); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PostGroup(w http.ResponseWriter, r *http.Request) {
	userID, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var group model.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	result, err := h.DBRepo.PostGroupAndGroupUser(&group, userID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupLastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := postInitGroupStandardBudgets(int(groupLastInsertId)); err != nil {
		if err := h.DBRepo.DeleteGroupAndGroupUser(int(groupLastInsertId), userID); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroup, err := h.DBRepo.GetGroup(int(groupLastInsertId))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&dbGroup); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PutGroup(w http.ResponseWriter, r *http.Request) {
	_, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var group model.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"group ID を正しく指定してください。"}))
		return
	}

	if err := h.DBRepo.PutGroup(&group, groupID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroup, err := h.DBRepo.GetGroup(groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"グループ名を取得できませんでした"}))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&dbGroup); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PostGroupUnapprovedUsers(w http.ResponseWriter, r *http.Request) {
	_, err := verifySessionID(h, w, r)
	if err != nil {
		if err == http.ErrNoCookie || err == redis.ErrNil {
			errorResponseByJSON(w, NewHTTPError(http.StatusUnauthorized, &AuthenticationErrorMsg{"このページを表示するにはログインが必要です。"}))
			return
		}
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var groupUnapprovedUser model.GroupUnapprovedUser
	if err := json.NewDecoder(r.Body).Decode(&groupUnapprovedUser); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := validateUserID(groupUnapprovedUser.UserID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, err))
		return
	}

	if err := h.DBRepo.FindUserID(groupUnapprovedUser.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"該当するユーザーが見つかりませんでした。"}))
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

	if err := checkForUniqueGroupUser(h, groupID, groupUnapprovedUser.UserID); err != nil {
		groupUserConflictErrorMsg, ok := err.(*GroupUserConflictErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusConflict, groupUserConflictErrorMsg))
		return
	}

	result, err := h.DBRepo.PostGroupUnapprovedUser(&groupUnapprovedUser, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroupUnapprovedUser, err := h.DBRepo.GetGroupUnapprovedUser(int(lastInsertId))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&dbGroupUnapprovedUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
