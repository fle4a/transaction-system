package main

import (
	"context"
	"log"
	"net/http"

	"github.com/fle4a/transaction-system/invoice-service/src/configs"
	"github.com/fle4a/transaction-system/invoice-service/src/internal/api"
	"github.com/fle4a/transaction-system/invoice-service/src/internal/core"
	"github.com/go-chi/chi/v5"
)

var db *core.DBPool

func main() {
	config, err := configs.ReadConfig()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	db = core.NewPool(context.Background(), core.CreateDBURL(config))
	if err := db.Init(); err != nil {
		log.Println(err)
		panic(err)
	}
	defer db.Close()
	r := chi.NewRouter()
	r.Post("/invoice", api.InvoiceHandler(db))
	log.Printf("Starting server on %s", config.Server.Addr)
	http.ListenAndServe(config.Server.Addr, r)
}
