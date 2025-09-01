package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os/signal"
	"syscall"
	"tech-wb-L0/backend/domain"
	"tech-wb-L0/backend/internal/repository"
)

type Consumer struct {
	reader  *kafka.Reader
	storage repository.OrderRepository
}

func NewConsumer(storage repository.OrderRepository) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9094"},
		GroupID:  "orders-consumer-group",
		Topic:    "orders",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	return &Consumer{
		r,
		storage,
	}
}

func (c *Consumer) StartConsumer() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	defer c.reader.Close()

	log.Println("Consumer started")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		var o domain.Order
		if err := json.Unmarshal(m.Value, &o); err != nil {
			log.Printf("bad json: %v", err)
		}

		err = c.storage.CreateOrder(ctx, &o)
		if err != nil {
			log.Printf("Error creating order: %v", err)
		} else {
			log.Printf("Consumed and created order: %v", o)
		}

	}

}
