package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

// max 10 goroutines for each partition (partitions_count = 10)
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
				timeToProcess := time.Duration(rand.Intn(3) + 1)
				// log.Printf("...(%s:%d:%d): %s (%ds)\n",
				// 	message.Topic, message.Partition, message.Offset, message.Key, timeToProcess)
				log.Printf("...(%s:%d:%d): (%ds)\n",
					message.Topic, message.Partition, message.Offset, timeToProcess)

				// Process simulation
				time.Sleep(timeToProcess * time.Second)

				// Mark as processed
				session.MarkMessage(message, string(message.Key))

				// log.Printf("OK.(%s)", message.Key)
				log.Printf("OK(%s:%d:%d): (%ds)\n",
					message.Topic, message.Partition, message.Offset, timeToProcess)
			}
		}(message)
	}

	// Wait for all messages to be processed and commit
	wg.Wait()
	session.Commit()

	return nil
}

func main() {
	log.SetFlags(log.Ltime)

	brokers := []string{"localhost:29091", "localhost:29092", "localhost:29093", "localhost:29094", "localhost:29095"}

	var topic string
	var groupID string
	flag.StringVar(&topic, "topic", "test-1", "Topic name")
	flag.StringVar(&groupID, "group", "test-1-group", "Consumer group")
	flag.Parse()

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

	log.Printf("Waiting for messages %s (%s) ...", topic, groupID)
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
