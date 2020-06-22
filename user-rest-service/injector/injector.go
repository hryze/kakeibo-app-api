package injector

import (
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/handler"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure"
)

func InjectDB() (infrastructure.SQLHandler, error) {
	sqlhandler, err := infrastructure.NewSQLHandler()
	return *sqlhandler, err
}

func InjectUserRepository() (repository.UserRepository, error) {
	sqlHandler, err := InjectDB()
	return infrastructure.NewUserRepository(sqlHandler), err
}

func InjectUserHandler() (handler.UserHandler, error) {
	userRepository, err := InjectUserRepository()
	return handler.NewUserHandler(userRepository), err
}
