package apierrors

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

type BadRequestError struct {
	RawError error
}

func (e *BadRequestError) Error() string {
	return e.RawError.Error()
}

func NewBadRequestError(err error) error {
	return &BadRequestError{
		RawError: err,
	}
}

type AuthenticationError struct {
	RawError error
}

func (e *AuthenticationError) Error() string {
	return e.RawError.Error()
}

func NewAuthenticationError(err error) error {
	return &AuthenticationError{
		RawError: err,
	}
}

type NotFoundError struct {
	RawError error
}

func (e *NotFoundError) Error() string {
	return e.RawError.Error()
}

func NewNotFoundError(err error) error {
	return &NotFoundError{
		RawError: err,
	}
}

type ConflictError struct {
	RawError error
}

func (e *ConflictError) Error() string {
	return e.RawError.Error()
}

func NewConflictError(err error) error {
	return &ConflictError{
		RawError: err,
	}
}

type InternalServerError struct {
	RawError error
}

func (e *InternalServerError) Error() string {
	return e.RawError.Error()
}

func NewInternalServerError(err error) error {
	return &InternalServerError{
		RawError: err,
	}
}
