package router

import (
	"net/http"

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

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Accept-Language"},
		AllowCredentials: true,
	})

	if err := http.ListenAndServe(":8082", corsWrapper.Handler(router)); err != nil {
		return err
	}

	return nil
}
