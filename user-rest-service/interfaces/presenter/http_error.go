package presenter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hryze/kakeibo-app-api/user-rest-service/apperrors"

	"github.com/hryze/kakeibo-app-api/user-rest-service/apierrors"
)

type httpError struct {
	StatusCode   int   `json:"status"`
	ErrorMessage error `json:"error"`
}

func (e *httpError) Error() string {
	b, err := json.Marshal(&e)
	if err != nil {
		log.Println(err)
	}

	return string(b)
}

func newHTTPError(err error) *httpError {
	var statusCode int
	var errorMessage error

	switch err := err.(type) {
	case *apierrors.BadRequestError:
		statusCode = http.StatusBadRequest
		errorMessage = err.RawError
	case *apierrors.AuthenticationError:
		statusCode = http.StatusUnauthorized
		errorMessage = err.RawError
	case *apierrors.NotFoundError:
		statusCode = http.StatusNotFound
		errorMessage = err.RawError
	case *apierrors.ConflictError:
		statusCode = http.StatusConflict
		errorMessage = err.RawError
	case *apierrors.InternalServerError:
		statusCode = http.StatusInternalServerError
		errorMessage = err.RawError
	default:
		statusCode = http.StatusInternalServerError
		errorMessage = apierrors.NewErrorString("Internal Server Error")
	}

	return &httpError{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	}
}

func ErrorJSON(w http.ResponseWriter, err error) {
	httpError := newHTTPError(err)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpError.StatusCode)
	if err := json.NewEncoder(w).Encode(httpError); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ErrorJSONV2(w http.ResponseWriter, err error) {
	appErr := apperrors.AsAppError(err)

	httpErr := &httpError{
		StatusCode:   appErr.Status(),
		ErrorMessage: appErr.InfoMessage(),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpErr.StatusCode)
	if err := json.NewEncoder(w).Encode(httpErr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
