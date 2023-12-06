package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	numOfPartitions := 3

	consumers := make([]*kafka.Consumer, numOfPartitions)
	for i := 0; i < numOfPartitions; i++ {
		c, err := NewConsumer()
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		consumers[i] = c
	}

	for _, c := range consumers {
		err := c.Subscribe("users", nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Starting consumers")

	for _, c := range consumers {
		go Consume(c)
	}

	// Block forever
	select {}
}

func NewConsumer() (*kafka.Consumer, error) {
	configs := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "golang-consumer",
		"group.id":          "golang-group",
		"auto.offset.reset": "earliest",
	}

	c, err := kafka.NewConsumer(configs)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func Consume(c *kafka.Consumer) {
	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}

		fmt.Printf(
			"User \"%s\" received from topic \"%s\" [%d] at offset %v\n",
			string(msg.Value),
			*msg.TopicPartition.Topic,
			msg.TopicPartition.Partition,
			msg.TopicPartition.Offset,
		)
	}
}
