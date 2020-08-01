package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/injector"
	"github.com/rs/cors"
)

func Run() error {
	h := injector.InjectDBHandler()

	router := mux.NewRouter()
	router.Handle("/signup", http.HandlerFunc(h.SignUp)).Methods("POST")
	router.Handle("/login", http.HandlerFunc(h.Login)).Methods("POST")
	router.Handle("/logout", http.HandlerFunc(h.Logout)).Methods("DELETE")
	router.Handle("/group", http.HandlerFunc(h.PostGroup)).Methods("POST")
	router.Handle("/group/{id:[0-9]+}", http.HandlerFunc(h.PutGroup)).Methods("PUT")

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Accept-Language"},
		AllowCredentials: true,
	})

	if err := http.ListenAndServe(":8080", corsWrapper.Handler(router)); err != nil {
		return err
	}

	return nil
}
