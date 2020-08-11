package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	router.HandleFunc("/transactions/{year_month:[0-9]{4}-[0-9]{2}}", h.GetMonthlyTransactionsList).Methods("GET")
	router.HandleFunc("/transactions", h.PostTransaction).Methods("POST")
	router.HandleFunc("/transactions/{id:[0-9]+}", h.PutTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{id:[0-9]+}", h.DeleteTransaction).Methods("DELETE")
	router.HandleFunc("/transactions/search", h.SearchTransactionsList).Methods("GET")
	router.HandleFunc("/standard-budgets", h.PostInitStandardBudgets).Methods("POST")
	router.HandleFunc("/standard-budgets", h.GetStandardBudgets).Methods("GET")
	router.HandleFunc("/standard-budgets", h.PutStandardBudgets).Methods("PUT")
	router.HandleFunc("/custom-budgets/{year_month:[0-9]{4}-[0-9]{2}}", h.GetCustomBudgets).Methods("GET")
	router.HandleFunc("/custom-budgets/{year_month:[0-9]{4}-[0-9]{2}}", h.PostCustomBudgets).Methods("POST")
	router.HandleFunc("/custom-budgets/{year_month:[0-9]{4}-[0-9]{2}}", h.PutCustomBudgets).Methods("PUT")
	router.HandleFunc("/custom-budgets/{year_month:[0-9]{4}-[0-9]{2}}", h.DeleteCustomBudgets).Methods("DELETE")
	router.HandleFunc("/budgets/{year:[0-9]{4}}", h.GetYearlyBudgets).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/categories", h.GetGroupCategoriesList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/categories/custom-categories", h.PostGroupCustomCategory).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}/categories/custom-categories/{id:[0-9]+}", h.PutGroupCustomCategory).Methods("PUT")
	router.HandleFunc("/groups/{group_id:[0-9]+}/categories/custom-categories/{id:[0-9]+}", h.DeleteGroupCustomCategory).Methods("DELETE")
	router.HandleFunc("/groups/{group_id:[0-9]+}/transactions/{year_month:[0-9]{4}-[0-9]{2}}", h.GetMonthlyGroupTransactionsList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/transactions", h.PostGroupTransaction).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}/transactions/{id:[0-9]+}", h.PutGroupTransaction).Methods("PUT")
	router.HandleFunc("/groups/{group_id:[0-9]+}/transactions/{id:[0-9]+}", h.DeleteGroupTransaction).Methods("DELETE")
	router.HandleFunc("/groups/{group_id:[0-9]+}/transactions/search", h.SearchGroupTransactionsList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/standard-budgets", h.PostInitGroupStandardBudgets).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}/standard-budgets", h.GetGroupStandardBudgets).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/standard-budgets", h.PutGroupStandardBudgets).Methods("PUT")
	router.HandleFunc("/groups/{group_id:[0-9]+}/custom-budgets/{year_month:[0-9]{4}-[0-9]{2}}", h.GetGroupCustomBudgets).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/custom-budgets/{year_month:[0-9]{4}-[0-9]{2}}", h.PostGroupCustomBudgets).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}/custom-budgets/{year_month:[0-9]{4}-[0-9]{2}}", h.PutGroupCustomBudgets).Methods("PUT")
	router.HandleFunc("/groups/{group_id:[0-9]+}/custom-budgets/{year_month:[0-9]{4}-[0-9]{2}}", h.DeleteGroupCustomBudgets).Methods("DELETE")
	router.HandleFunc("/groups/{group_id:[0-9]+}/budgets/{year:[0-9]{4}}", h.GetYearlyGroupBudgets).Methods("GET")

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Accept-Language"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Addr:    ":8081",
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
