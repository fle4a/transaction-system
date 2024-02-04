package main

import (
	"context"
	"log"
	"strings"

	"github.com/fle4a/transaction-system/withdraw-service/src/internal/core"
	"github.com/fle4a/transaction-system/withdraw-service/src/configs"
)

var (
	consumer *core.KafkaConsumer
	db       *core.DBPool
)



func main() {
	config, err := configs.ReadConfig()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	db = core.NewPool(context.Background(), core.CreateDBURL(config))
	consumer = core.NewConsumer(strings.Join(config.Kafka.BrokerList, ","), config.Kafka.ConsumerTopics, config.Kafka.GroupID, db)
	if err := db.Init(); err != nil {
		log.Println(err)
		panic(err)
	}
	defer db.Close()

	if err := consumer.Init(); err != nil {
		log.Println(err)
		panic(err)
	}
	defer consumer.Close()
	select {}
}
