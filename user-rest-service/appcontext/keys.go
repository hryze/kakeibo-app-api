package appcontext

type contextKey string

const (
	userIDKey   contextKey = "USER_ID"
	appErrorKey contextKey = "APP_ERROR"
)
