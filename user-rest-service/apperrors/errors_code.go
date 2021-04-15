package apperrors

type code string

const (
	invalidParameter    code = "InvalidParameter"
	unauthorized        code = "Unauthorized"
	notFound            code = "NotFound"
	conflict            code = "Conflict"
	internalServerError code = "InternalServerError"
)

func (c code) value() string {
	return string(c)
}
