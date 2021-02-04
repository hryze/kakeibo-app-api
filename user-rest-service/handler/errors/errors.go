package errors

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/xerrors"

	uerrors "github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/errors"
)

type httpError struct {
	Status       int   `json:"status"`
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
	var status int
	var errorMessage error

	switch e := err.(type) {
	case *uerrors.BadRequestError:
		status = http.StatusBadRequest
		errorMessage = e.RawError
	case *uerrors.AuthenticationError:
		status = http.StatusUnauthorized
		errorMessage = e.RawError
	case *uerrors.NotFoundError:
		status = http.StatusNotFound
		errorMessage = e.RawError
	case *uerrors.ConflictError:
		status = http.StatusConflict
		errorMessage = e.RawError
	case *uerrors.InternalServerError:
		status = http.StatusInternalServerError
		errorMessage = e.RawError
	default:
		status = http.StatusInternalServerError
		errorMessage = xerrors.New("Internal Server Error")
	}

	return &httpError{
		Status:       status,
		ErrorMessage: errorMessage,
	}
}

func ErrorResponseByJSON(w http.ResponseWriter, err error) {
	httpError := newHTTPError(err)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpError.Status)
	if err := json.NewEncoder(w).Encode(httpError); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
