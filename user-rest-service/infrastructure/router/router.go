package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/auth"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/auth/imdb"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/externalapi"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/externalapi/client"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/persistence"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/persistence/rdb"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/injector"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/interfaces/handler"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/usecase"
)

func Run() error {
	redisHandler, err := imdb.NewRedisHandler()
	if err != nil {
		return err
	}
	defer redisHandler.Pool.Close()

	mySQLHandler, err := rdb.NewMySQLHandler()
	if err != nil {
		return err
	}
	defer mySQLHandler.Conn.Close()

	accountApiHandler := client.NewAccountApiHandler()

	userRepository := persistence.NewUserRepository(mySQLHandler)
	sessionStore := auth.NewSessionStore(redisHandler)
	accountApi := externalapi.NewAccountApi(accountApiHandler)
	userUsecase := usecase.NewUserUsecase(userRepository, sessionStore, accountApi)
	userHandler := handler.NewUserHandler(userUsecase)

	h := injector.InjectDBHandler()

	router := mux.NewRouter()
	router.HandleFunc("/readyz", h.Readyz).Methods(http.MethodGet)
	router.HandleFunc("/signup", userHandler.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/logout", userHandler.Logout).Methods(http.MethodDelete)
	router.HandleFunc("/user", h.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/groups", h.GetGroupList).Methods(http.MethodGet)
	router.HandleFunc("/groups", h.PostGroup).Methods(http.MethodPost)
	router.HandleFunc("/groups/{group_id:[0-9]+}", h.PutGroup).Methods(http.MethodPut)
	router.HandleFunc("/groups/{group_id:[0-9]+}/users", h.GetGroupUserIDList).Methods(http.MethodGet)
	router.HandleFunc("/groups/{group_id:[0-9]+}/users", h.PostGroupUnapprovedUser).Methods(http.MethodPost)
	router.HandleFunc("/groups/{group_id:[0-9]+}/users", h.DeleteGroupApprovedUser).Methods(http.MethodDelete)
	router.HandleFunc("/groups/{group_id:[0-9]+}/users/approved", h.PostGroupApprovedUser).Methods(http.MethodPost)
	router.HandleFunc("/groups/{group_id:[0-9]+}/users/unapproved", h.DeleteGroupUnapprovedUser).Methods(http.MethodDelete)
	router.HandleFunc("/groups/{group_id:[0-9]+}/users/{user_id:[\\S]{1,10}}/verify", h.VerifyGroupAffiliation).Methods(http.MethodGet)
	router.HandleFunc("/groups/{group_id:[0-9]+}/users/verify", h.VerifyGroupAffiliationOfUsersList).Methods(http.MethodGet)

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   config.Env.Cors.AllowedOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Accept-Language"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Env.Server.Port),
		Handler: corsWrapper.Handler(router),
	}

	errorCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			errorCh <- err
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errorCh:
		return err
	case s := <-signalCh:
		log.Printf("SIGNAL %s received", s.String())
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}
