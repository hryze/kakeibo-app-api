package middleware

import (
	"log"
	"net/http"

	"github.com/hryze/kakeibo-app-api/user-rest-service/appcontext"
	"github.com/hryze/kakeibo-app-api/user-rest-service/apperrors"
	"github.com/hryze/kakeibo-app-api/user-rest-service/config"
)

func NewLoggingMiddlewareFunc() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			err, ok := appcontext.GetAppError(r.Context())
			if !ok {
				// Successfully processed
				return
			}

			appErr := apperrors.AsAppError(err)
			if appErr == nil {
				log.Print("failed to type assertion for appError")
				return
			}

			if config.Env.Log.Debug {
				log.Printf("%+v", appErr)
				return
			}

			if appErr.IsLevelError() || appErr.IsLevelCritical() {
				// Transfer logs to CloudWatch Logs
				log.Printf("%+v", appErr)
				return
			}
		})
	}
}
