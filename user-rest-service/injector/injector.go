package injector

import (
	"fmt"
	"os"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/handler"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure"
)

func InjectDB() infrastructure.SQLHandler {
	SQLHandler, err := infrastructure.NewSQLHandler()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return *SQLHandler
}

func InjectUserRepository() *infrastructure.UserRepository {
	SQLHandler := InjectDB()
	return infrastructure.NewUserRepository(SQLHandler)
}

func InjectUserHandler() *handler.UserHandler {
	userRepository := InjectUserRepository()
	return handler.NewUserHandler(userRepository)
}
