package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/interfaces/presenter"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
)

func assignColorCodeToUser(groupUserIDList []string) string {
	const (
		red                  = "#FF0000"
		cyan                 = "#00FFFF"
		chartreuseGreen      = "#80FF00"
		violet               = "#8000FF"
		orange               = "#FF8000"
		azure                = "#0080FF"
		emeraldGreen         = "#00FF80"
		rubyRed              = "#FF0080"
		yellow               = "#FFFF00"
		blue                 = "#0000FF"
		green                = "#00FF00"
		magenta              = "#FF00FF"
		goldenYellow         = "#FFBF00"
		cobaltBlue           = "#0040FF"
		brightYellowishGreen = "#BFFF00"
		hyacinth             = "#4000FF"
		cobaltGreen          = "#00FF40"
		reddishPurple        = "#FF00BF"
		leafGreen            = "#40FF00"
		purple               = "#BF00FF"
		vermilion            = "#FF4000"
		ceruleanBlue         = "#00BFFF"
		turquoiseGreen       = "#00FFBF"
		carmine              = "#FF0040"
	)

	colorCodeList := [24]string{
		red,
		cyan,
		chartreuseGreen,
		violet,
		orange,
		azure,
		emeraldGreen,
		rubyRed,
		yellow,
		blue,
		green,
		magenta,
		goldenYellow,
		cobaltBlue,
		brightYellowishGreen,
		hyacinth,
		cobaltGreen,
		reddishPurple,
		leafGreen,
		purple,
		vermilion,
		ceruleanBlue,
		turquoiseGreen,
		carmine,
	}

	idx := len(groupUserIDList) % len(colorCodeList)

	colorCode := colorCodeList[idx]

	return colorCode
}

type groupHandler struct {
	groupUsecase usecase.GroupUsecase
}

func NewGroupHandler(groupUsecase usecase.GroupUsecase) *groupHandler {
	return &groupHandler{
		groupUsecase: groupUsecase,
	}
}

func (h *groupHandler) FetchGroupList(w http.ResponseWriter, r *http.Request) {
	in, err := getUserIDOfContext(r)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	out, err := h.groupUsecase.FetchGroupList(in)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusOK, out)
}

func (h *groupHandler) StoreGroup(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := getUserIDOfContext(r)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	var group input.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("正しいデータを入力してください")))
		return
	}

	out, err := h.groupUsecase.StoreGroup(authenticatedUser, &group)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusCreated, out)
}

func (h *groupHandler) UpdateGroupName(w http.ResponseWriter, r *http.Request) {
	var group input.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("正しいデータを入力してください")))
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("グループIDを正しく指定してください")))
		return
	}

	group.GroupID = groupID

	out, err := h.groupUsecase.UpdateGroupName(&group)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusOK, out)
}

func (h *groupHandler) StoreGroupUnapprovedUser(w http.ResponseWriter, r *http.Request) {
	var unapprovedUser input.UnapprovedUser
	if err := json.NewDecoder(r.Body).Decode(&unapprovedUser); err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("正しいデータを入力してください")))
		return
	}

	var group input.Group
	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("グループIDを正しく指定してください")))
		return
	}

	group.GroupID = groupID

	out, err := h.groupUsecase.StoreGroupUnapprovedUser(&unapprovedUser, &group)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusCreated, out)
}

func (h *groupHandler) DeleteGroupApprovedUser(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, err := getUserIDOfContext(r)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		presenter.ErrorJSON(w, apierrors.NewBadRequestError(apierrors.NewErrorString("グループIDを正しく指定してください")))
		return
	}

	group := input.Group{GroupID: groupID}

	if err := h.groupUsecase.DeleteGroupApprovedUser(authenticatedUser, &group); err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusOK, presenter.NewSuccessString("グループを退会しました"))
}

func (h *DBHandler) PostGroupApprovedUser(w http.ResponseWriter, r *http.Request) {
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

	if err := h.GroupRepo.FindUnapprovedUser(groupID, userID); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"グループに招待されていないため、参加できませんでした。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	groupUserIDList, err := h.GroupRepo.GetGroupUsersList(groupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	colorCode := assignColorCodeToUser(groupUserIDList)

	result, err := h.GroupRepo.PostGroupApprovedUserAndDeleteGroupUnapprovedUser(groupID, userID, colorCode)
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	approvedUser, err := h.GroupRepo.GetApprovedUser(int(lastInsertId))
	if err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&approvedUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) DeleteGroupUnapprovedUser(w http.ResponseWriter, r *http.Request) {
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

	if err := h.GroupRepo.FindUnapprovedUser(groupID, userID); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			errorResponseByJSON(w, NewHTTPError(http.StatusBadRequest, &BadRequestErrorMsg{"こちらのグループには招待されていません。"}))
			return
		}

		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	if err := h.GroupRepo.DeleteGroupUnapprovedUser(groupID, userID); err != nil {
		errorResponseByJSON(w, NewHTTPError(http.StatusInternalServerError, nil))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&DeleteContentMsg{"グループ招待を拒否しました。"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *DBHandler) VerifyGroupAffiliation(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := mux.Vars(r)["user_id"]

	if err := h.GroupRepo.FindApprovedUser(groupID, userID); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *DBHandler) VerifyGroupAffiliationOfUsersList(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var groupUsersList model.GroupTasksUsersListReceiver
	if err := json.NewDecoder(r.Body).Decode(&groupUsersList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbGroupUsersList, err := h.GroupRepo.FindApprovedUsersList(groupID, groupUsersList.GroupUsersList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(groupUsersList.GroupUsersList) != len(dbGroupUsersList.GroupUsersList) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *DBHandler) GetGroupUserIDList(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(mux.Vars(r)["group_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	groupUserIDList, err := h.GroupRepo.GetGroupUsersList(groupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(groupUserIDList) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&groupUserIDList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
