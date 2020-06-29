package router

import (
	"net/http"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/injector"

	"github.com/gorilla/mux"
)

func Run() error {
	h := injector.InjectUserHandler()

	router := mux.NewRouter()
	router.Handle("/signup", http.HandlerFunc(h.SignUp)).Methods("POST")
	router.Handle("/login", http.HandlerFunc(h.Login)).Methods("POST")
	if err := http.ListenAndServe(":8080", router); err != nil {
		return err
	}

	return nil
}
