package router

import (
	"net/http"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/handler"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/injector"

	"github.com/gorilla/mux"
)

func Run() error {
	userHandler := injector.InjectUserHandler()
	h := handler.NewResponseHandler(userHandler)

	router := mux.NewRouter()
	router.Handle("/signup", http.HandlerFunc(h.ResponseByJSONHandler)).Methods("POST")
	router.Handle("/login", http.HandlerFunc(h.ResponseByJSONHandler)).Methods("POST")
	if err := http.ListenAndServe(":8080", router); err != nil {
		return err
	}

	return nil
}
