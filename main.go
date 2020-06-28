package main

import (
	"context"
	"go-api/practice/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "Products-Api", log.LstdFlags)
	sm := mux.NewRouter()

	ph := handlers.NewProducts(l)
	bh := handlers.NewBanks(l)
	uh := handlers.NewUsers(l)
	ah := handlers.NewAssets(l)

	GetRouter := sm.Methods(http.MethodGet).Subrouter()
	GetRouter.HandleFunc("/", ph.GetProducts)

	BankGetRouter := sm.Methods(http.MethodGet).Subrouter()
	BankGetRouter.HandleFunc("/banks", bh.GetBanks)

	UserGetRouter := sm.Methods(http.MethodGet).Subrouter()
	UserGetRouter.HandleFunc("/users", uh.GetUsers)

	AssetGetRouter := sm.Methods(http.MethodGet).Subrouter()
	AssetGetRouter.HandleFunc("/assets", ah.GetAssets)

	PostRouter := sm.Methods(http.MethodPost).Subrouter()
	PostRouter.HandleFunc("/", ph.AddProducts)
	PostRouter.Use(ph.MiddlewareValidateProduct)

	BankPostRouter := sm.Methods(http.MethodPost).Subrouter()
	BankPostRouter.HandleFunc("/banks", bh.AddBanks)
	BankPostRouter.Use(bh.MiddlewareValidateBank)

	UserPostRouter := sm.Methods(http.MethodPost).Subrouter()
	UserPostRouter.HandleFunc("/users", uh.AddUsers)
	UserPostRouter.Use(uh.MiddlewareValidateUser)

	AssetPostRouter := sm.Methods(http.MethodPost).Subrouter()
	AssetPostRouter.HandleFunc("/assets", ah.AddAssets)
	AssetPostRouter.Use(ah.MiddlewareValidateAssets)

	PutRouter := sm.Methods(http.MethodPut).Subrouter()
	PutRouter.HandleFunc("/{id:[0-9]+}", ah.UpdateAssets)
	PutRouter.Use(ah.MiddlewareValidateAssets)

	s := http.Server{
		Addr:         ":9090",
		Handler:      sm,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	SigChan := make(chan os.Signal)

	signal.Notify(SigChan, os.Interrupt)
	signal.Notify(SigChan, os.Kill)

	sig := <-SigChan
	log.Println("Notified To shutdown", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
