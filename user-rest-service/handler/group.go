package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/garyburd/redigo/redis"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
)

type NoContentMsg struct {
	Message string `json:"message"`
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
			if groupList[i].ID == groupUser.GroupID {
				groupList[i].GroupUsersList = append(groupList[i].GroupUsersList, groupUser.UserName)
			}
		}

		for _, groupUnapprovedUser := range groupUnapprovedUsersList {
			if groupList[i].ID == groupUnapprovedUser.GroupID {
				groupList[i].GroupUnapprovedUsersList = append(groupList[i].GroupUnapprovedUsersList, groupUnapprovedUser.UserName)
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

	groupID, err := strconv.Atoi(mux.Vars(r)["id"])
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
