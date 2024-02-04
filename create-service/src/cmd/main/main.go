package main

import (
	"context"
	"log"
	"strings"

	"github.com/fle4a/transaction-system/create-service/src/configs"
	"github.com/fle4a/transaction-system/create-service/src/internal/core"
	"github.com/fle4a/transaction-system/create-service/src/internal/types"
)

var (
	consumer *core.KafkaConsumer
	producer *core.KafkaProducer
	db       *core.DBPool
)

func main() {
	config, err := configs.ReadConfig()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	transactionChannel := make(chan types.Transaction)
	db = core.NewPool(context.Background(), core.CreateDBURL(config))
	consumer = core.NewConsumer(strings.Join(config.Kafka.BrokerList, ","), config.Kafka.ConsumerTopics, config.Kafka.GroupID, db, transactionChannel)
	producer = core.NewProducer(strings.Join(config.Kafka.BrokerList, ","), config.Kafka.ProducerTopic, transactionChannel)
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

	if err := producer.Init(); err != nil {
		log.Println(err)
		panic(err)
	}
	defer producer.Close()

	select {}
}
