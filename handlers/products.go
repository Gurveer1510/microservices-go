package handlers

import (
	"context"
	"log"
	"microservices/data"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	prod, ok := r.Context().Value(KeyProduct{}).(*data.Product)
	if !ok {
		// p.l.Println(prod)
		http.Error(rw, "Missing Product in request body", http.StatusBadRequest)
		return
	}
	data.AddProduct(prod)
	// p.l.Printf("Product: %#v", prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to read id", http.StatusBadRequest)
	}
	prod, ok := r.Context().Value(KeyProduct{}).(*data.Product)
	if !ok {
		http.Error(rw, "Missing Product", http.StatusBadRequest)
		return
	}
	prod.ID = id
	data.UpdateProduct(prod)
	// p.l.Printf("Product: %#v", prod)
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductsValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to read body", http.StatusBadRequest)
			return
		}
		
		err = prod.ProductValidator()
		if err != nil {
			validationErrors, ok := err.(validator.ValidationErrors)
			if ok {
				http.Error(rw, validationErrors.Error(), http.StatusBadRequest)
				return
			}
		}
		
		p.l.Println("IN THE MIDDLEWARE:")
		p.l.Printf("Product: %v", prod)
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
