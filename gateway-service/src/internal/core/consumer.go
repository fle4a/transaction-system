package core

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/fle4a/transaction-system/gateway-service/src/internal/types"
)

type KafkaProducer struct {
	Producer   *kafka.Producer
	BrokerList string
	Topic      string
}

func NewProducer(brokerList string, topic string) *KafkaProducer {
	return &KafkaProducer{
		BrokerList: brokerList,
		Topic:      topic,
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

	return nil
}

func (p *KafkaProducer) Produce(data types.Transaction) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return err
	}
	err = p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.Topic, Partition: kafka.PartitionAny},
		Value:          bytes,
	}, nil)

	if err != nil {
		log.Println(err)
	}
	return err
}

func (p *KafkaProducer) Close() {
	log.Println("Closing producer")
	p.Producer.Close()
}
