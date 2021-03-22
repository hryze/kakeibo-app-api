package handler

import (
	"net/http"

	"github.com/gorilla/context"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/input"
)

func getUserIDOfContext(r *http.Request) (*input.AuthenticatedUser, error) {
	ctx, ok := context.GetOk(r, config.Env.RequestCtx.UserID)
	if !ok {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	ctxUserID, ok := ctx.(string)
	if !ok {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return &input.AuthenticatedUser{
		UserID: ctxUserID,
	}, nil
}
