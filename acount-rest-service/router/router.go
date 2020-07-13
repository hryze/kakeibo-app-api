package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paypay3/kakeibo-app-api/acount-rest-service/injector"
	"github.com/rs/cors"
)

func Run() error {
	h := injector.InjectDBHandler()

	router := mux.NewRouter()
	router.HandleFunc("/categories", h.GetCategories).Methods("GET")
	router.HandleFunc("/custom-category", h.PostCustomCategory).Methods("POST")
	router.HandleFunc("/custom-category", h.PutCustomCategory).Methods("PUT")
	router.HandleFunc("/custom-category", h.DeleteCustomCategory).Methods("DELETE")

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Accept-Language"},
		AllowCredentials: true,
	})

	if err := http.ListenAndServe(":8081", corsWrapper.Handler(router)); err != nil {
		return err
	}

	return nil
}
