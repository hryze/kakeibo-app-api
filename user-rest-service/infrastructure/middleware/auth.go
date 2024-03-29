package middleware

import (
	"net/http"
	"regexp"

	"golang.org/x/xerrors"

	"github.com/hryze/kakeibo-app-api/user-rest-service/apierrors"
	"github.com/hryze/kakeibo-app-api/user-rest-service/appcontext"
	"github.com/hryze/kakeibo-app-api/user-rest-service/config"
	"github.com/hryze/kakeibo-app-api/user-rest-service/interfaces/presenter"
	"github.com/hryze/kakeibo-app-api/user-rest-service/usecase/sessionstore"
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
			userID, err := sessionStore.FetchUserByUserID(sessionID)
			if err != nil {
				presenter.ErrorJSON(w, err)
				return
			}

			ctx := appcontext.SetUserID(r.Context(), userID.Value())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

const (
	getGroupUserIDListHandlerPathFormat                = `^/groups/(?P<v0>[0-9]+)/users$`
	verifyGroupAffiliationHandlerPathFormat            = `^/groups/(?P<v0>[0-9]+)/users/(?P<v1>[\S]{1,10})/verify$`
	verifyGroupAffiliationOfUsersListHandlerPathFormat = `^/groups/(?P<v0>[0-9]+)/users/verify$`
)

var (
	skipAuthMiddlewarePaths = [...]string{
		"/readyz",
		"/signup",
		"/login",
		"/logout",
	}

	skipAuthMiddlewareHandlers = [...]*regexp.Regexp{
		regexp.MustCompile(getGroupUserIDListHandlerPathFormat),
		regexp.MustCompile(verifyGroupAffiliationHandlerPathFormat),
		regexp.MustCompile(verifyGroupAffiliationOfUsersListHandlerPathFormat),
	}
)

func skipAuthMiddleware(r *http.Request) bool {
	requestPath := r.URL.Path

	for _, path := range skipAuthMiddlewarePaths {
		if requestPath == path {
			return true
		}
	}

	if r.Method == http.MethodGet {
		for _, regex := range skipAuthMiddlewareHandlers {
			if regex.MatchString(requestPath) {
				return true
			}
		}
	}

	return false
}
