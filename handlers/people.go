package handlers

import (
	"context"
	"fmt"
	"go-api/practice/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Peoples struct {
	l *log.Logger
}

func NewPeoples(l *log.Logger) *Peoples {
	return &Peoples{l}
}

func (p *Peoples) GetPeople(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET People")
	pl := data.GetPeople()
	err := pl.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
	}
}

func (p *Peoples) AddPeople(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST People")
	peep := r.Context().Value(keyPeople{}).(data.People)
	data.AddPeople(&peep)
}

func (p Peoples) UpdatePeople(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert ID", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT People")
	peep := r.Context().Value(keyPeople{}).(data.People)
	err = data.UpdatePeople(id, &peep)
	if err == data.ErrPeopleNotFound {
		http.Error(rw, "People Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "People Not Found", http.StatusInternalServerError)
		return
	}
}

type keyPeople struct{}

func (p Peoples) MiddlewareValidatePeople(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		peep := data.People{}
		err := peep.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		// Validate People

		err = peep.Validate()
		if err != nil {
			p.l.Println("Error Validating People")
			http.Error(
				rw,
				fmt.Sprintf("Unable to validate People %s", err),
				http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), keyPeople{}, peep)

		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
