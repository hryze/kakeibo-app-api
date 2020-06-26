package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/handler"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/injector"
)

func Run() error {
	h := injector.InjectUserHandler()

	router := mux.NewRouter()
	router.HandleFunc("/user", handler.ErrorCheckHandler(h.SignUp)).Methods("POST")
	if err := http.ListenAndServe(":8080", router); err != nil {
		return err
	}

	return nil
}
