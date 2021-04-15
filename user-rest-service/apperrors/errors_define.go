package apperrors

var (
	InvalidParameter    = newBadRequest(invalidParameter)
	Unauthorized        = newUnauthorized(unauthorized)
	NotFound            = newNotFound(notFound)
	Conflict            = newConflict(conflict)
	InternalServerError = newInternalServerError(internalServerError)
)
