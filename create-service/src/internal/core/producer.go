package core

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/fle4a/transaction-system/create-service/src/internal/types"
)

type KafkaProducer struct {
	Producer           *kafka.Producer
	BrokerList         string
	Topic              string
	TransactionChannel <-chan types.Transaction
}

func NewProducer(brokerList string, topic string, chanel <-chan types.Transaction) *KafkaProducer {
	return &KafkaProducer{
		BrokerList:         brokerList,
		Topic:              topic,
		TransactionChannel: chanel,
	}
}

func (p *KafkaProducer) Init() error {
	c, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": p.BrokerList,
	})
	if err != nil {
		return err
	}
	p.Producer = c

	go p.initProduce()
	log.Printf("Producer start")
	return nil
}

func (p *KafkaProducer) initProduce() {
	for transaction := range p.TransactionChannel {
		bytes, err := json.Marshal(transaction)
		if err != nil {
			log.Println(err)
			continue
		}
		err = p.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &p.Topic, Partition: kafka.PartitionAny},
			Value:          bytes,
		}, nil)
		if err != nil {
			log.Println(err)
		}
	}
}

func (p *KafkaProducer) Close() {
	log.Println("Closing producer")
	p.Producer.Close()
}
