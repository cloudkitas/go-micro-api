package handlers

import (
	"context"
	"go-api/practice/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Companies struct {
	l *log.Logger
}

func NewCompanies(l *log.Logger) *Companies {
	return &Companies{l}
}

func (c *Companies) GetCompanies(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET Companies")
	cl := data.GetCompany()
	err := cl.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
	}
}

func (c *Companies) AddCompany(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle POST Company")
	comp := r.Context().Value(KeyCompany{}).(data.Company)

	data.AddCompany(&comp)
}

func (c Companies) UpdateCompany(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert ID", http.StatusBadRequest)
		return
	}

	c.l.Println("Handle PUT Company", id)
	comp := r.Context().Value(KeyCompany{}).(data.Company)
	err = data.UpdateCompany(id, &comp)

	if err == data.ErrCompanyNotFound {
		http.Error(rw, "Company Not Found", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(rw, "Company Not Found", http.StatusNotFound)
		return
	}

}

type KeyCompany struct{}

func (c Companies) MiddlewareValidateCompany(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		comp := data.Company{}
		err := comp.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		// validate company

		ctx := context.WithValue(r.Context(), KeyCompany{}, comp)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
