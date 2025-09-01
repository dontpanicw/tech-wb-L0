package main

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"kafka-producer/domain"
	"log"
	"os"
	"time"
)

func main() {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9094"),
		Topic:    "orders",
		Balancer: &kafka.Hash{},
	}
	defer w.Close()

	data, err := os.ReadFile("orders.json")
	if err != nil {
		log.Fatalf("failed to read orders.json: %v", err)
	}

	var orders []domain.Order
	if err := json.Unmarshal(data, &orders); err != nil {
		log.Fatalf("failed to unmarshal orders.json: %v", err)
	}
	for _, order := range orders {
		b, _ := json.Marshal(order)
		msg := &kafka.Message{
			Key:   []byte(order.OrderUID),
			Value: b,
		}

		if err := w.WriteMessages(context.Background(), *msg); err != nil {
			log.Fatalf("failed to write message to kafka: %v", err)
		}
		log.Printf("order written to kafka: %v", order.OrderUID)
		time.Sleep(1 * time.Second)
	}
	log.Printf("end")

}
