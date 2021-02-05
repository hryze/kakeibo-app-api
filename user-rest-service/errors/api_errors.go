package errors

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/xerrors"
)

type errorString struct {
	Message string `json:"message"`
}

func (e *errorString) Error() string {
	return e.Message
}

func NewErrorString(message string) error {
	return &errorString{
		Message: message,
	}
}

type badRequestError struct {
	RawError error
}

func (e *badRequestError) Error() string {
	return e.RawError.Error()
}

func NewBadRequestError(err error) error {
	return &badRequestError{
		RawError: err,
	}
}

type authenticationError struct {
	RawError error
}

func (e *authenticationError) Error() string {
	return e.RawError.Error()
}

func NewAuthenticationError(err error) error {
	return &authenticationError{
		RawError: err,
	}
}

type notFoundError struct {
	RawError error
}

func (e *notFoundError) Error() string {
	return e.RawError.Error()
}

func NewNotFoundError(err error) error {
	return &notFoundError{
		RawError: err,
	}
}

type conflictError struct {
	RawError error
}

func (e *conflictError) Error() string {
	return e.RawError.Error()
}

func NewConflictError(err error) error {
	return &conflictError{
		RawError: err,
	}
}

type internalServerError struct {
	RawError error
}

func (e *internalServerError) Error() string {
	return e.RawError.Error()
}

func NewInternalServerError(err error) error {
	return &internalServerError{
		RawError: err,
	}
}

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
	case *badRequestError:
		statusCode = http.StatusBadRequest
		errorMessage = err.RawError
	case *authenticationError:
		statusCode = http.StatusUnauthorized
		errorMessage = err.RawError
	case *notFoundError:
		statusCode = http.StatusNotFound
		errorMessage = err.RawError
	case *conflictError:
		statusCode = http.StatusConflict
		errorMessage = err.RawError
	case *internalServerError:
		statusCode = http.StatusInternalServerError
		errorMessage = err.RawError
	default:
		statusCode = http.StatusInternalServerError
		errorMessage = xerrors.New("Internal Server Error")
	}

	return &httpError{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	}
}

func ErrorResponseByJSON(w http.ResponseWriter, err error) {
	httpError := newHTTPError(err)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpError.StatusCode)
	if err := json.NewEncoder(w).Encode(httpError); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
