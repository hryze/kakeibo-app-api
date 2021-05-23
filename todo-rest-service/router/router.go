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

	"github.com/hryze/kakeibo-app-api/todo-rest-service/config"
	"github.com/hryze/kakeibo-app-api/todo-rest-service/injector"
)

func Run() error {
	h := injector.InjectDBHandler()

	router := mux.NewRouter()
	router.HandleFunc("/readyz", h.Readyz).Methods("GET")
	router.HandleFunc("/todo-list/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", h.GetDailyTodoList).Methods("GET")
	router.HandleFunc("/todo-list/{year_month:[0-9]{4}-[0-9]{2}}", h.GetMonthlyTodoList).Methods("GET")
	router.HandleFunc("/todo-list/expired", h.GetExpiredTodoList).Methods("GET")
	router.HandleFunc("/todo-list", h.PostTodo).Methods("POST")
	router.HandleFunc("/todo-list/{id:[0-9]+}", h.PutTodo).Methods("PUT")
	router.HandleFunc("/todo-list/{id:[0-9]+}", h.DeleteTodo).Methods("DELETE")
	router.HandleFunc("/todo-list/search", h.SearchTodoList).Methods("GET")
	router.HandleFunc("/shopping-list/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}/daily", h.GetDailyShoppingDataByDay).Methods("GET")
	router.HandleFunc("/shopping-list/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}/categories", h.GetDailyShoppingDataByCategory).Methods("GET")
	router.HandleFunc("/shopping-list/{year_month:[0-9]{4}-[0-9]{2}}/daily", h.GetMonthlyShoppingDataByDay).Methods("GET")
	router.HandleFunc("/shopping-list/{year_month:[0-9]{4}-[0-9]{2}}/categories", h.GetMonthlyShoppingDataByCategory).Methods("GET")
	router.HandleFunc("/shopping-list/expired", h.GetExpiredShoppingList).Methods("GET")
	router.HandleFunc("/shopping-list/regular", h.PostRegularShoppingItem).Methods("POST")
	router.HandleFunc("/shopping-list/regular/{id:[0-9]+}", h.PutRegularShoppingItem).Methods("PUT")
	router.HandleFunc("/shopping-list/regular/{id:[0-9]+}", h.DeleteRegularShoppingItem).Methods("DELETE")
	router.HandleFunc("/shopping-list", h.PostShoppingItem).Methods("POST")
	router.HandleFunc("/shopping-list/{id:[0-9]+}", h.PutShoppingItem).Methods("PUT")
	router.HandleFunc("/shopping-list/{id:[0-9]+}", h.DeleteShoppingItem).Methods("DELETE")
	router.HandleFunc("/shopping-list/categories", h.PutShoppingListCustomCategoryIdToMediumCategoryId).Methods("PUT")
	router.HandleFunc("/groups/{group_id:[0-9]+}/todo-list/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", h.GetDailyGroupTodoList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/todo-list/{year_month:[0-9]{4}-[0-9]{2}}", h.GetMonthlyGroupTodoList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/todo-list/expired", h.GetExpiredGroupTodoList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/todo-list", h.PostGroupTodo).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}/todo-list/{id:[0-9]+}", h.PutGroupTodo).Methods("PUT")
	router.HandleFunc("/groups/{group_id:[0-9]+}/todo-list/{id:[0-9]+}", h.DeleteGroupTodo).Methods("DELETE")
	router.HandleFunc("/groups/{group_id:[0-9]+}/todo-list/search", h.SearchGroupTodoList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/shopping-list/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}/daily", h.GetDailyGroupShoppingDataByDay).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/shopping-list/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}/categories", h.GetDailyGroupShoppingDataByCategory).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/shopping-list/{year_month:[0-9]{4}-[0-9]{2}}/daily", h.GetMonthlyGroupShoppingDataByDay).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/shopping-list/{year_month:[0-9]{4}-[0-9]{2}}/categories", h.GetMonthlyGroupShoppingDataByCategory).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/shopping-list/expired", h.GetExpiredGroupShoppingList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/shopping-list/regular", h.PostGroupRegularShoppingItem).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}/shopping-list/regular/{id:[0-9]+}", h.PutGroupRegularShoppingItem).Methods("PUT")
	router.HandleFunc("/groups/{group_id:[0-9]+}/shopping-list/regular/{id:[0-9]+}", h.DeleteGroupRegularShoppingItem).Methods("DELETE")
	router.HandleFunc("/groups/{group_id:[0-9]+}/shopping-list", h.PostGroupShoppingItem).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}/shopping-list/{id:[0-9]+}", h.PutGroupShoppingItem).Methods("PUT")
	router.HandleFunc("/groups/{group_id:[0-9]+}/shopping-list/{id:[0-9]+}", h.DeleteGroupShoppingItem).Methods("DELETE")
	router.HandleFunc("/groups/{group_id:[0-9]+}/shopping-list/categories", h.PutGroupShoppingListCustomCategoryIdToMediumCategoryId).Methods("PUT")
	router.HandleFunc("/groups/{group_id:[0-9]+}/tasks/users", h.GetGroupTasksListForEachUser).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/tasks/users", h.PostGroupTasksUsersList).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}/tasks/users", h.DeleteGroupTasksUsersList).Methods("DELETE")
	router.HandleFunc("/groups/{group_id:[0-9]+}/tasks", h.GetGroupTasksList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/tasks", h.PostGroupTask).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}/tasks/{id:[0-9]+}", h.PutGroupTask).Methods("PUT")
	router.HandleFunc("/groups/{group_id:[0-9]+}/tasks/{id:[0-9]+}", h.DeleteGroupTask).Methods("DELETE")

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   config.Env.Cors.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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
