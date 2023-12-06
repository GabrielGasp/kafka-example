package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	producer, err := NewProducer()
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	users, err := GetUsers()
	if err != nil {
		log.Fatal(err)
	}

	deliveryChan := make(chan kafka.Event)

	go ReportDelivery(deliveryChan)

	for _, user := range users {
		err = Produce(user.Name, "users", producer, []byte(user.ID), deliveryChan)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(10 * time.Millisecond) // Sleep between each message to make communication more visible
	}

	// Block forever
	select {}
}

func NewProducer() (*kafka.Producer, error) {
	configs := &kafka.ConfigMap{
		"bootstrap.servers":   "localhost:9092",
		"acks":                "all",
		"max.in.flight":       5,
		"enable.idempotence":  true,
		"delivery.timeout.ms": 10000,
	}

	p, err := kafka.NewProducer(configs)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func Produce(content, topic string, producer *kafka.Producer, key []byte, deliverChan chan kafka.Event) error {
	msg := &kafka.Message{
		Value: []byte(content),
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key: key,
	}

	err := producer.Produce(msg, deliverChan)
	if err != nil {
		return err
	}

	return nil
}

func ReportDelivery(deliveryChan chan kafka.Event) {
	for e := range deliveryChan {
		switch ev := e.(type) {

		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				// Do whatever you need when delivery fails (retry, etc.)
				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
			} else {
				// Do whatever you need after confirming delivery
				fmt.Printf(
					"User \"%s\" delivered to topic \"%s\" [%d] at offset %v\n",
					string(ev.Value),
					*ev.TopicPartition.Topic,
					ev.TopicPartition.Partition,
					ev.TopicPartition.Offset,
				)
			}

		default:
			fmt.Printf("Ignored event: %s\n", ev)
		}
	}
}
