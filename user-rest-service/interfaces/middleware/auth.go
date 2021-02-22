package middleware

import (
	"net/http"

	"github.com/gorilla/context"
	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/interfaces/presenter"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase/sessionstore"
)

func NewAuthMiddlewareFunc(sessionStore sessionstore.SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if skip := skipAuthMiddleware(r); skip {
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := r.Cookie(config.Env.Cookie.Name)
			if xerrors.Is(err, http.ErrNoCookie) {
				presenter.ErrorJSON(w, apierrors.NewAuthenticationError(apierrors.NewErrorString("このページを表示するにはログインが必要です")))
				return
			}

			sessionID := cookie.Value
			userID, err := sessionStore.FetchUserID(sessionID)
			if err != nil {
				presenter.ErrorJSON(w, err)
				return
			}

			context.Set(r, config.Env.RequestCtx.UserID, userID.Value())

			next.ServeHTTP(w, r)
		})
	}
}

func skipAuthMiddleware(r *http.Request) bool {
	var skipAuthMiddlewarePaths = [...]string{
		"/readyz",
		"/signup",
		"/login",
		"/logout",
	}

	requestPath := r.URL.Path
	for _, path := range skipAuthMiddlewarePaths {
		if requestPath == path {
			return true
		}
	}

	return false
}
