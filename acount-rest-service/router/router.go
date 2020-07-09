package router

import (
	"net/http"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/injector"

	"github.com/gorilla/mux"
)

func Run() error {
	h := injector.InjectDBHandler()

	router := mux.NewRouter()
	router.HandleFunc("/categories", h.GetCategories).Methods("GET")

	if err := http.ListenAndServe(":8081", router); err != nil {
		return err
	}

	return nil
}
