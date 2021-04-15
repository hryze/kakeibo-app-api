package appcontext

import "context"

func SetAppError(ctx context.Context, appErr error) context.Context {
	return context.WithValue(ctx, appErrorKey, appErr)
}

func GetAppError(ctx context.Context) (error, bool) {
	appErr, ok := ctx.Value(appErrorKey).(error)

	return appErr, ok
}
