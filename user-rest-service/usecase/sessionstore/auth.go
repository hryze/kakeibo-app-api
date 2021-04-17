package sessionstore

import "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/userdomain"

type SessionStore interface {
	PingSessionStore() error
	StoreUserBySessionID(sessionID string, userID userdomain.UserID) error
	DeleteUserBySessionID(sessionID string) error
	FetchUserByUserID(sessionID string) (userdomain.UserID, error)
}
