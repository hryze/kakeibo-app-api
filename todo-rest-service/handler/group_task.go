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

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/config"
	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

func verifyGroupAffiliationOfUsersList(groupID int, groupUsersList model.GroupTasksUsersListReceiver) error {
	requestURL := fmt.Sprintf(
		"http://%s:%d/groups/%d/users/verify",
		config.Env.UserApi.Host, config.Env.UserApi.Port, groupID,
	)

	requestBody, err := json.Marshal(groupUsersList)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		"GET",
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

	if response.StatusCode == http.StatusBadRequest {
		return &BadRequestErrorMsg{"こちらのグループには、指定されたユーザーは所属していません。"}
	}

	if response.StatusCode == http.StatusInternalServerError {
		return &InternalServerErrorMsg{"500 Internal Server Error"}
	}

	return nil
}

func validateGroupTaskName(groupTask *model.GroupTask) error {
	if strings.HasPrefix(groupTask.TaskName, " ") || strings.HasPrefix(groupTask.TaskName, "　") {
		return &BadRequestErrorMsg{"タスク名の文字列先頭に空白がないか確認してください。"}
	}

	if strings.HasSuffix(groupTask.TaskName, " ") || strings.HasSuffix(groupTask.TaskName, "　") {
		return &BadRequestErrorMsg{"タスク名の文字列末尾に空白がないか確認してください。"}
	}

	if utf8.RuneCountInString(groupTask.TaskName) == 0 || utf8.RuneCountInString(groupTask.TaskName) > 20 {
		return &BadRequestErrorMsg{"タスク名は1文字以上20文字以内で入力してください。"}
	}

	return nil
}

func generateNextBaseDate(today time.Time, baseDate time.Time, cycleDays int) (time.Time, int) {
	diffDays := int(today.Sub(baseDate).Hours() / 24)
	exceededDays := diffDays % cycleDays

	nextBaseDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -exceededDays)
	cycleCount := diffDays / cycleDays

	return nextBaseDate, cycleCount
}

func generateNextGroupTasksUserID(groupTaskAssignedToUser model.GroupTask, groupTasksUsersList []model.GroupTasksUser, cycleCount int) int {
	var userIndex int
	usersListLength := len(groupTasksUsersList)

	for i, groupTasksUser := range groupTasksUsersList {
		if groupTasksUser.ID == groupTaskAssignedToUser.GroupTasksUserID.Int {
			userIndex = i
			break
		}
	}

	nextGroupTasksUserIdIndex := (cycleCount%usersListLength + userIndex) % usersListLength
	nextGroupTasksUserID := groupTasksUsersList[nextGroupTasksUserIdIndex].ID

	return nextGroupTasksUserID
}

func (h *DBHandler) GetGroupTasksListForEachUser(w http.ResponseWriter, r *http.Request) {
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

	groupTasksUsersList, err := h.GroupTasksRepo.GetGroupTasksUsersList(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupTasksListAssignedToUser, err := h.GroupTasksRepo.GetGroupTasksListAssignedToUser(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var updateTaskIndexList []int
	now := h.TimeManage.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

	for i := 0; i < len(groupTasksListAssignedToUser); i++ {
		baseDate := groupTasksListAssignedToUser[i].BaseDate.Time
		cycleDays := groupTasksListAssignedToUser[i].Cycle.Int
		taskEndDate := baseDate.AddDate(0, 0, cycleDays).Add(-1 * time.Second)

		if !today.After(taskEndDate) {
			continue
		}

		if groupTasksListAssignedToUser[i].CycleType.String == "none" {
			groupTasksListAssignedToUser[i].BaseDate.Valid = false
			groupTasksListAssignedToUser[i].CycleType.Valid = false
			groupTasksListAssignedToUser[i].Cycle.Valid = false
			groupTasksListAssignedToUser[i].GroupTasksUserID.Valid = false

			updateTaskIndexList = append(updateTaskIndexList, i)

			continue
		}

		nextBaseDate, cycleCount := generateNextBaseDate(today, baseDate, cycleDays)
		nextGroupTasksUserID := generateNextGroupTasksUserID(groupTasksListAssignedToUser[i], groupTasksUsersList, cycleCount)

		groupTasksListAssignedToUser[i].BaseDate.Time = nextBaseDate
		groupTasksListAssignedToUser[i].GroupTasksUserID.Int = nextGroupTasksUserID

		updateTaskIndexList = append(updateTaskIndexList, i)
	}

	if len(updateTaskIndexList) != 0 {
		if err := h.GroupTasksRepo.PutGroupTasksListAssignedToUser(groupTasksListAssignedToUser, updateTaskIndexList); err != nil {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}
	}

	for i := 0; i < len(groupTasksUsersList); i++ {
		for j := 0; j < len(groupTasksListAssignedToUser); j++ {
			if groupTasksUsersList[i].ID == groupTasksListAssignedToUser[j].GroupTasksUserID.Int {
				groupTasksUsersList[i].TasksList = append(groupTasksUsersList[i].TasksList, groupTasksListAssignedToUser[j])
			}
		}
	}

	groupTasksListForEachUser := model.NewGroupTasksListForEachUser(groupTasksUsersList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(groupTasksListForEachUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) PostGroupTasksUsersList(w http.ResponseWriter, r *http.Request) {
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

	var groupTasksUsersListReceiver model.GroupTasksUsersListReceiver
	if err := json.NewDecoder(r.Body).Decode(&groupTasksUsersListReceiver); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := verifyGroupAffiliationOfUsersList(groupID, groupTasksUsersListReceiver); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	dbGroupTasksUsersList, err := h.GroupTasksRepo.GetGroupTasksUsersList(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for _, userID := range groupTasksUsersListReceiver.GroupUsersList {
		for _, dbUser := range dbGroupTasksUsersList {
			if userID == dbUser.UserID {
				errorResponseByJSON(w, NewHTTPError(http.StatusConflict, &ConflictErrorMsg{"選択したユーザーは、既にタスクメンバーに追加されています。"}))
				return
			}
		}
	}

	if err := h.GroupTasksRepo.PostGroupTasksUsersList(groupTasksUsersListReceiver, groupID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroupTasksUsersListSender, err := h.GroupTasksRepo.GetGroupTasksUsersList(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	for _, groupTasksUser := range dbGroupTasksUsersList {
		for i, dbGroupTasksUser := range dbGroupTasksUsersListSender {
			if groupTasksUser.ID == dbGroupTasksUser.ID {
				dbGroupTasksUsersListSender = append(dbGroupTasksUsersListSender[:i], dbGroupTasksUsersListSender[i+1:]...)
			}
		}
	}

	groupTasksListForEachUser := model.NewGroupTasksListForEachUser(dbGroupTasksUsersListSender)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&groupTasksListForEachUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) DeleteGroupTasksUsersList(w http.ResponseWriter, r *http.Request) {
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

	var groupTasksUsersListReceiver model.GroupTasksUsersListReceiver
	if err := json.NewDecoder(r.Body).Decode(&groupTasksUsersListReceiver); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := verifyGroupAffiliationOfUsersList(groupID, groupTasksUsersListReceiver); err != nil {
		badRequestErrorMsg, ok := err.(*BadRequestErrorMsg)
		if !ok {
			errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, badRequestErrorMsg))
		return
	}

	dbGroupTasksUsersList, err := h.GroupTasksRepo.GetGroupTasksUsersList(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	var groupTasksUsersIdList []int

	for _, userID := range groupTasksUsersListReceiver.GroupUsersList {
		var isExist bool

		for _, dbUser := range dbGroupTasksUsersList {
			if userID == dbUser.UserID {
				groupTasksUsersIdList = append(groupTasksUsersIdList, dbUser.ID)

				isExist = true
				break
			}
		}

		if isExist {
			continue
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"選択したユーザーは、既にタスクメンバーから削除されています。"}))
		return
	}

	groupTasksIDList, err := h.GroupTasksRepo.GetGroupTasksIDListAssignedToUser(groupTasksUsersIdList, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := h.GroupTasksRepo.DeleteGroupTasksUsersList(groupTasksUsersListReceiver, groupTasksIDList, groupID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&DeleteContentMsg{"タスクメンバーを削除しました。"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) GetGroupTasksList(w http.ResponseWriter, r *http.Request) {
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

	groupTasksList, err := h.GroupTasksRepo.GetGroupTasksList(groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	senderGroupTasksList := model.NewGroupTasksList(groupTasksList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(senderGroupTasksList); err != nil {
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

	result, err := h.GroupTasksRepo.PostGroupTask(groupTask, groupID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	dbGroupTask, err := h.GroupTasksRepo.GetGroupTask(int(lastInsertId))
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

func (h *DBHandler) PutGroupTask(w http.ResponseWriter, r *http.Request) {
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

	groupTasksID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"タスクIDを正しく指定してください。"}))
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

	result, err := h.GroupTasksRepo.PutGroupTask(&groupTask, groupTasksID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	} else if rowsAffected == 0 {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"指定されたタスクは存在しません。"}))
		return
	}

	dbGroupTask, err := h.GroupTasksRepo.GetGroupTask(groupTasksID)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dbGroupTask); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) DeleteGroupTask(w http.ResponseWriter, r *http.Request) {
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

	groupTasksID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"タスクIDを正しく指定してください。"}))
		return
	}

	if _, err := h.GroupTasksRepo.GetGroupTask(groupTasksID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"指定されたタスクは存在しません。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := h.GroupTasksRepo.DeleteGroupTask(groupTasksID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(DeleteContentMsg{"タスクを削除しました。"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
