package apperrors

import (
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

func newBadRequest(code code) *appError {
	var e appError
	e.code = code
	e.status = http.StatusBadRequest
	e.level = levelInfo

	return &e
}

func newUnauthorized(code code) *appError {
	var e appError
	e.code = code
	e.status = http.StatusUnauthorized
	e.level = levelInfo

	return &e
}

func newNotFound(code code) *appError {
	var e appError
	e.code = code
	e.status = http.StatusNotFound
	e.level = levelInfo

	return &e
}

func newConflict(code code) *appError {
	var e appError
	e.code = code
	e.status = http.StatusConflict
	e.level = levelInfo

	return &e
}

func newInternalServerError(code code) *appError {
	var e appError
	e.code = code
	e.status = http.StatusInternalServerError
	e.level = levelError

	return &e
}

func (e appError) create(msg string) *appError {
	e.message = msg
	e.frame = xerrors.Caller(2)

	return &e
}

func (e appError) Wrap(err error, msg ...string) *appError {
	var m string
	if len(msg) != 0 {
		m = msg[0]
	} else {
		m = e.Code()
	}

	ne := e.create(m)
	ne.next = err

	return ne
}

func (e *appError) SetInfoMessage(err error) *appError {
	e.infoMessage = err

	return e
}

func (e *appError) Code() string {
	if e.code != "" {
		return e.code.value()
	}

	next := AsAppError(e.next)
	if next != nil {
		return next.Code()
	}

	return "NotDefined"
}

func (e *appError) Status() int {
	if e.status != 0 {
		return e.status
	}

	next := AsAppError(e.next)
	if next != nil {
		return next.Status()
	}

	return http.StatusInternalServerError
}

func (e *appError) InfoMessage() error {
	if e.infoMessage != nil {
		return e.infoMessage
	}

	next := AsAppError(e.next)
	if next != nil {
		return next.InfoMessage()
	}

	return nil
}
