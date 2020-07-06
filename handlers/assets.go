package handlers

import (
	"context"
	"go-api/practice/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Assets struct {
	l *log.Logger
}

func NewAssets(l *log.Logger) *Assets {
	return &Assets{l}
}

func (a *Assets) GetAssets(rw http.ResponseWriter, r *http.Request) {
	a.l.Println("Handle GET Assets")
	al := data.GetAssets()
	err := al.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
	}
}

func (a *Assets) AddAssets(rw http.ResponseWriter, r *http.Request) {
	a.l.Println("Handle POST Asset")
	asst := r.Context().Value(KeyAsset{}).(data.Asset)
	data.AddAssets(&asst)
}

func (a Assets) UpdateAssets(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert ID", http.StatusBadRequest)
		return
	}

	a.l.Println("Handle PUT Asset", id)
	asst := r.Context().Value(KeyAsset{}).(data.Asset)

	err = data.UpdateAsset(id, &asst)
	if err == data.ErrAssetNotFound {
		http.Error(rw, "Unable to Find Asset", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Unable to find Asset", http.StatusInternalServerError)
		return
	}

}

type KeyAsset struct{}

func (a Assets) MiddlewareValidateAssets(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		asst := data.Asset{}
		err := asst.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		}

		ctx := context.WithValue(r.Context(), KeyAsset{}, asst)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})

}
