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
	router.HandleFunc("/signup", handler.ResponseByJSONMiddleware(h.SignUp)).Methods("POST")
	router.HandleFunc("/login", handler.ResponseByJSONMiddleware(h.Login)).Methods("POST")
	if err := http.ListenAndServe(":8080", router); err != nil {
		return err
	}

	return nil
}
