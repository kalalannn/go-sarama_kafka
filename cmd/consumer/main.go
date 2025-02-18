package main

import (
	"context"
	"log"
	"math/rand"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

// max 10 goroutines for each partition (partitions_count = 3)
const maxMessageGoroutines = 5

type ConsumerHandler struct{}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var wg sync.WaitGroup
	guardChan := make(chan struct{}, maxMessageGoroutines)

	for message := range claim.Messages() {
		guardChan <- struct{}{}
		wg.Add(1)
		go func(message *sarama.ConsumerMessage) {
			defer wg.Done()
			defer func() { <-guardChan }()

			select {
			case <-session.Context().Done():
				return
			default:
				timeToProcess := time.Duration(rand.Intn(3) + 2)
				log.Printf("Processing ... (%s:%d:%d): %s (%ds)\n",
					message.Topic, message.Partition, message.Offset, message.Key, timeToProcess)

				// Process simulation
				time.Sleep(timeToProcess * time.Second)

				// Mark as processed
				session.MarkMessage(message, string(message.Key))

				log.Printf("Processed successfully. (%s)", message.Key)
			}
		}(message)
	}

	// Wait for all messages to be processed and commit
	wg.Wait()
	session.Commit()

	return nil
}

func main() {
	brokers := []string{"localhost:29092", "localhost:39092", "localhost:49092"}
	topic := "test"
	groupID := "test-group"

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.AutoCommit.Enable = false

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			log.Fatalf("Failed to close consumer group: %v", err)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	handler := &ConsumerHandler{}

	log.Println("Waiting for messages...")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
				log.Fatalf("Error from consumer group: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-ctx.Done()
	log.Println("Received interrupt, waiting for consumer group to finish...")

	wg.Wait()
	log.Println("Consumer gracefully stopped.")
}
