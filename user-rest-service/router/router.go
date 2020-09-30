package router

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/injector"
	"github.com/rs/cors"
)

func Run() error {
	isLocal := flag.Bool("local", false, "Please specify -env flag")
	flag.Parse()

	h := injector.InjectDBHandler(*isLocal)

	router := mux.NewRouter()
	router.Handle("/signup", http.HandlerFunc(h.SignUp)).Methods("POST")
	router.Handle("/login", http.HandlerFunc(h.Login)).Methods("POST")
	router.Handle("/logout", http.HandlerFunc(h.Logout)).Methods("DELETE")
	router.Handle("/groups", http.HandlerFunc(h.GetGroupList)).Methods("GET")
	router.Handle("/groups", http.HandlerFunc(h.PostGroup)).Methods("POST")
	router.Handle("/groups/{group_id:[0-9]+}", http.HandlerFunc(h.PutGroup)).Methods("PUT")
	router.Handle("/groups/{group_id:[0-9]+}/users", http.HandlerFunc(h.PostGroupUnapprovedUser)).Methods("POST")
	router.Handle("/groups/{group_id:[0-9]+}/users", http.HandlerFunc(h.DeleteGroupApprovedUser)).Methods("DELETE")
	router.Handle("/groups/{group_id:[0-9]+}/users/approved", http.HandlerFunc(h.PostGroupApprovedUser)).Methods("POST")
	router.Handle("/groups/{group_id:[0-9]+}/users/unapproved", http.HandlerFunc(h.DeleteGroupUnapprovedUser)).Methods("DELETE")
	router.Handle("/groups/{group_id:[0-9]+}/users/{user_id:[\\S]{1,10}}", http.HandlerFunc(h.VerifyGroupAffiliation)).Methods("GET")

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
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
