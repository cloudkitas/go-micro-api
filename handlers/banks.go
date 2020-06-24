package handlers

import (
	"context"
	"go-api/practice/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Banks struct {
	l *log.Logger
}

func NewBanks(l *log.Logger) *Banks {
	return &Banks{l}
}

func (b *Banks) GetBanks(rw http.ResponseWriter, r *http.Request) {
	b.l.Println("Handle GET Banks")
	bl := data.GetBanks()
	err := bl.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
	}
}

func (b *Banks) AddBanks(rw http.ResponseWriter, r *http.Request) {
	b.l.Println("Handle POST Bank")
	bnk := r.Context().Value(KeyBank{}).(data.Bank)
	data.AddBank(&bnk)
}

func (b Banks) UpdateBanks(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert ID", http.StatusBadRequest)
		return
	}
	if err == data.ErrBankNotFound {
		http.Error(rw, "Bank Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Bank Not Found", http.StatusInternalServerError)
		return
	}

	b.l.Println("Handle PUT Bank", id)
	bnk := r.Context().Value(KeyBank{}).(data.Bank)

	data.UpdateBank(id, &bnk)

}

type KeyBank struct{}

func (b Banks) MiddlewareValidateBank(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		bnk := data.Bank{}
		err := bnk.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal bank", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyBank{}, bnk)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
