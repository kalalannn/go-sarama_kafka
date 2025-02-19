package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

func main() {
	brokers := []string{"localhost:29091", "localhost:29092", "localhost:29093", "localhost:29094", "localhost:29095"}

	var topic string
	var count int

	flag.IntVar(&count, "count", 100, "Count of events")
	flag.StringVar(&topic, "topic", "test-1", "Topic name")
	flag.Parse()

	config := sarama.NewConfig()
	// config.Producer.Idempotent = true
	// config.Producer.RequiredAcks = sarama.WaitForAll
	// config.Net.MaxOpenRequests = 1
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 100 * time.Millisecond

	producer, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Failed to close producer: %v", err)
		}
	}()

	for i := 0; i < count; i++ {
		key := uuid.New().String()
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(key),
			Value: sarama.StringEncoder(fmt.Sprintf("Value: %d", i)),
		}
		if _, _, err := producer.SendMessage(msg); err != nil {
			log.Fatalf("Failed to produce message: %v", err)
		}
	}
	log.Println("Messages sent successfully")
}
