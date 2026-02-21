package handlers

import (
	"log"
	"microservices/data"
	"net/http"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		p.getProducts(rw, r)
		return
	case http.MethodPost:
		p.addProduct(rw, r)
		return
	case http.MethodPut:
		p.updateProduct(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to read body", http.StatusBadRequest)
	}
	data.AddProduct(prod)
	p.l.Printf("Product: %#v", prod)
}

func (p *Products) updateProduct(rw http.ResponseWriter, r *http.Request) {

	p.l.Println(r.URL.Path)
}
