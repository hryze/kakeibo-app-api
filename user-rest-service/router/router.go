package router

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/injector"
	"github.com/rs/cors"
)

func Run() error {
	isLocal := flag.Bool("local", false, "Please specify -local flag")
	flag.Parse()

	if *isLocal {
		if err := godotenv.Load("../../.env"); err != nil {
			return err
		}
	}

	if len(os.Getenv("ALLOWED_ORIGIN")) == 0 || len(os.Getenv("COOKIE_DOMAIN")) == 0 || len(os.Getenv("ACCOUNT_HOST")) == 0 || len(os.Getenv("MYSQL_DSN")) == 0 || len(os.Getenv("REDIS_DSN")) == 0 {
		return errors.New("environment variable not defined")
	}

	h := injector.InjectDBHandler()

	router := mux.NewRouter()
	router.HandleFunc("/readyz", h.Readyz).Methods("GET")
	router.HandleFunc("/signup", h.SignUp).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")
	router.HandleFunc("/logout", h.Logout).Methods("DELETE")
	router.HandleFunc("/groups", h.GetGroupList).Methods("GET")
	router.HandleFunc("/groups", h.PostGroup).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}", h.PutGroup).Methods("PUT")
	router.HandleFunc("/groups/{group_id:[0-9]+}/users", h.PostGroupUnapprovedUser).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}/users", h.DeleteGroupApprovedUser).Methods("DELETE")
	router.HandleFunc("/groups/{group_id:[0-9]+}/users/approved", h.PostGroupApprovedUser).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}/users/unapproved", h.DeleteGroupUnapprovedUser).Methods("DELETE")
	router.HandleFunc("/groups/{group_id:[0-9]+}/users/{user_id:[\\S]{1,10}}", h.VerifyGroupAffiliation).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/users", h.VerifyGroupAffiliationOfUsersList).Methods("GET")

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Accept-Language"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Addr:    ":8080",
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
