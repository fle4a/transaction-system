package core

import (
	"encoding/json"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/fle4a/transaction-system/create-service/src/internal/types"
)

type KafkaConsumer struct {
	Consumer           *kafka.Consumer
	BrokerList         string
	TopicList          []string
	GroupID            string
	Db                 *DBPool
	TransactionChannel chan<- types.Transaction
}

func NewConsumer(brokerList string, topicList []string, groupID string, db *DBPool, chanel chan<- types.Transaction) *KafkaConsumer {
	return &KafkaConsumer{
		BrokerList:         brokerList,
		TopicList:          topicList,
		GroupID:            groupID,
		Db:                 db,
		TransactionChannel: chanel,
	}
}

func (kc *KafkaConsumer) Init() error {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     kc.BrokerList,
		"broker.address.family": "v4",
		"group.id":              kc.GroupID,
		"auto.offset.reset":     "earliest",
	})
	if err != nil {
		return err
	}
	kc.Consumer = c

	kc.Consumer.SubscribeTopics(kc.TopicList, nil)
	
	go kc.initConsume()
	log.Printf("Consumer start")
	return nil
}

func (kc *KafkaConsumer) initConsume() {
	for {
		ev := kc.Consumer.Poll(100)
		if ev == nil {
			continue
		}
		switch e := ev.(type) {
		case *kafka.Message:
			var transaction types.Transaction
			if err := json.Unmarshal(e.Value, &transaction); err != nil {
				log.Println(err)
				continue
			}
			if err := kc.Db.createTransaction(transaction); err != nil {
				log.Println(err)
				continue
			}
			kc.TransactionChannel <- transaction
		case kafka.Error:
			log.Printf("%% Error: %v: %v\n", e.Code(), e)
			if e.Code() == kafka.ErrAllBrokersDown {
				os.Exit(1)
			}
		default:
			log.Printf("Ignored %v\n", e)
		}
	}
}

func (kc *KafkaConsumer) Close() {
	log.Println("Closing consumer")
	kc.Consumer.Close()
}
