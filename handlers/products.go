package handlers

import (
	"context"
	"go-api/practice/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	pl := data.GetProducts()
	err := pl.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
	}
}

func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

type KeyProduct struct{}

func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert ID", http.StatusBadRequest)
	}

	if err == data.ErrProductNotFound {
		http.Error(rw, "Unable to Find Product", http.StatusNotFound)
	}

	if err != nil {
		http.Error(rw, "Unable to find product", http.StatusInternalServerError)
	}

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Println("Handle PUT Product", id)
	data.UpdateProduct(id, &prod)

}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to umarshal json", http.StatusBadRequest)
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)

		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
