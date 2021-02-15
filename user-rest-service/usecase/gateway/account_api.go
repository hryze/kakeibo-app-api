package gateway

import "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/userdomain"

type AccountApi interface {
	PostInitStandardBudgets(userID userdomain.UserID) error
}
