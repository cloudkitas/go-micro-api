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

type Users struct {
	l *log.Logger
}

func NewUsers(l *log.Logger) *Users {
	return &Users{l}
}

func (u *Users) GetUsers(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET User")
	ul := data.GetUsers()
	err := ul.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
	}

}

func (u *Users) AddUsers(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle POST User")
	usr := r.Context().Value(KeyUser{}).(data.User)
	data.AddUsers(&usr)
}

func (u Users) UpdateUsers(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unavke to convert ID", http.StatusBadRequest)
		return
	}
	if err == data.ErrUserNotFound {
		http.Error(rw, "ERROR: User not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "ERROR: User Not Found", http.StatusInternalServerError)
		return
	}

	usr := r.Context().Value(KeyUser{}).(data.User)

	u.l.Println("Handle PUT User", id)
	data.UpdateUser(id, &usr)

}

type KeyUser struct{}

func (u Users) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		usr := data.User{}

		err := usr.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		//validate users json schema

		err = usr.Validate()
		if err != nil {
			http.Error(
				rw,
				fmt.Sprintf("Error Validating User: %s", err),
				http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUser{}, usr)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
