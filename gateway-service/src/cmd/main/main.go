package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/fle4a/transaction-system/gateway-service/src/internal/api"
	"github.com/fle4a/transaction-system/gateway-service/src/configs"
	"github.com/fle4a/transaction-system/gateway-service/src/internal/core"
	"github.com/go-chi/chi/v5"
)

var producer *core.KafkaProducer

func main() {
	config, err := configs.ReadConfig()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	producer = core.NewProducer(strings.Join(config.Kafka.BrokerList, ","), config.Kafka.Topic)
	if err := producer.Init(); err != nil {
		log.Println(err)
		panic(err)
	}
	defer producer.Close()
	r := chi.NewRouter()

	r.Post("/withdraw", api.TransactionHandler(producer))
	r.Post("/balance", api.BalanceHandler)

	log.Printf("Starting server on %s", config.Server.Addr)
	http.ListenAndServe(config.Server.Addr, r)
}
