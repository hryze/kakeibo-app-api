package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/paypay3/kakeibo-app-api/todo-rest-service/injector"
	"github.com/rs/cors"
)

func Run() error {
	h := injector.InjectDBHandler()

	router := mux.NewRouter()
	router.HandleFunc("/todo-list/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", h.GetDailyTodoList).Methods("GET")
	router.HandleFunc("/todo-list/{year_month:[0-9]{4}-[0-9]{2}}", h.GetMonthlyTodoList).Methods("GET")
	router.HandleFunc("/todo-list", h.PostTodo).Methods("POST")
	router.HandleFunc("/todo-list/{id:[0-9]+}", h.PutTodo).Methods("PUT")
	router.HandleFunc("/todo-list/{id:[0-9]+}", h.DeleteTodo).Methods("DELETE")
	router.HandleFunc("/todo-list/search", h.SearchTodoList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/todo-list/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", h.GetDailyGroupTodoList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/todo-list/{year_month:[0-9]{4}-[0-9]{2}}", h.GetMonthlyGroupTodoList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/todo-list", h.PostGroupTodo).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}/todo-list/{id:[0-9]+}", h.PutGroupTodo).Methods("PUT")
	router.HandleFunc("/groups/{group_id:[0-9]+}/todo-list/{id:[0-9]+}", h.DeleteGroupTodo).Methods("DELETE")
	router.HandleFunc("/groups/{group_id:[0-9]+}/todo-list/search", h.SearchGroupTodoList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/tasks/users", h.PostGroupTasksUser).Methods("POST")

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Accept-Language"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Addr:    ":8082",
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
