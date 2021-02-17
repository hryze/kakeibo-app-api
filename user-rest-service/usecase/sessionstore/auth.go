package sessionstore

import "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/userdomain"

type SessionStore interface {
	AddSessionID(sessionID string, userID userdomain.UserID) error
	DeleteSessionID(sessionID string) error
}
