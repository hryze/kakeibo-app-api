package injector

import (
	"github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/handler"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure"
)

func InjectDB() (infrastructure.SQLHandler, error) {
	SQLHandler, err := infrastructure.NewSQLHandler()
	return *SQLHandler, err
}

func InjectUserRepository() (repository.UserRepository, error) {
	SQLHandler, err := InjectDB()
	return infrastructure.NewUserRepository(SQLHandler), err
}

func InjectUserHandler() (handler.UserHandler, error) {
	userRepository, err := InjectUserRepository()
	return handler.NewUserHandler(userRepository), err
}