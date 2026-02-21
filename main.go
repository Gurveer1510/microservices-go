package main

import (
	"context"
	"log"
	"microservices/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "http: ", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods("GET").Subrouter()
	putRouter := sm.Methods("PUT").Subrouter()
	postRouter := sm.Methods("POST").Subrouter()

	putRouter.Use(ph.MiddlewareProductsValidator)
	postRouter.Use(ph.MiddlewareProductsValidator)

	getRouter.HandleFunc("/", ph.GetProducts)
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	postRouter.HandleFunc("/", ph.AddProduct)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			l.Fatalf("Error starting server: %s\n", err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown:", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	if err := s.Shutdown(tc); err != nil {
		l.Fatalf("Error shutting down server: %s\n", err)
	}
}
