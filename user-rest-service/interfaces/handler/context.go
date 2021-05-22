package handler

import (
	"net/http"

	"github.com/hryze/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/hryze/kakeibo-app-api/user-rest-service/appcontext"
	"github.com/hryze/kakeibo-app-api/user-rest-service/usecase/input"
)

func getUserIDForContext(r *http.Request) (*input.AuthenticatedUser, error) {
	userID, ok := appcontext.GetUserID(r.Context())
	if !ok {
		return nil, apierrors.NewInternalServerError(apierrors.NewErrorString("Internal Server Error"))
	}

	return &input.AuthenticatedUser{
		UserID: userID,
	}, nil
}

func setAppErrorToContext(r *http.Request, appErr error) { //nolint // To be used in the next implementation
	ctx := appcontext.SetAppError(r.Context(), appErr)
	*r = *r.WithContext(ctx)
}
