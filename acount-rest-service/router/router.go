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
	router.HandleFunc("/categories", h.GetCategoriesList).Methods("GET")
	router.HandleFunc("/categories/custom-categories", h.PostCustomCategory).Methods("POST")
	router.HandleFunc("/categories/custom-categories/{id:[0-9]+}", h.PutCustomCategory).Methods("PUT")
	router.HandleFunc("/categories/custom-categories/{id:[0-9]+}", h.DeleteCustomCategory).Methods("DELETE")
	router.HandleFunc("/transactions/{month:[0-9]{4}-[0-9]{2}}", h.GetMonthlyTransactionsList).Methods("GET")
	router.HandleFunc("/transactions", h.PostTransaction).Methods("POST")
	router.HandleFunc("/transactions/{id:[0-9]+}", h.PutTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{id:[0-9]+}", h.DeleteTransaction).Methods("DELETE")
	router.HandleFunc("/transactions/search", h.SearchTransactionsList).Methods("GET")
	router.HandleFunc("/standard-budgets", h.PostInitStandardBudgets).Methods("POST")
	router.HandleFunc("/standard-budgets", h.GetStandardBudgets).Methods("GET")
	router.HandleFunc("/standard-budgets", h.PutStandardBudgets).Methods("PUT")
	router.HandleFunc("/custom-budgets/{year_month:[0-9]{4}-[0-9]{2}}", h.PostCustomBudgets).Methods("POST")

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
