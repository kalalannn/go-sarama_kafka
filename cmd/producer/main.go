package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

func main() {
	brokers := []string{"localhost:29092", "localhost:39092", "localhost:49092"}
	topic := "test"

	config := sarama.NewConfig()
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

	for i := 0; i < 21; i++ {
		key := uuid.New().String()
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(key),
			Value: sarama.StringEncoder(fmt.Sprintf("Hello World %d", i)),
		}
		if _, _, err := producer.SendMessage(msg); err != nil {
			log.Fatalf("Failed to produce message: %v", err)
		}
	}
	log.Println("Messages sent successfully")
}
