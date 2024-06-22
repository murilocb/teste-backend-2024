package products

import (
	"context"
	"encoding/json"
	"log"
	"ms-go/app/models"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

func InitKafkaProducer() {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP([]string{"kafka:29092"}...),
		Balancer: &kafka.LeastBytes{},
		Topic:    "go-to-rails",
	}
}

func CloseKafkaProducer() {
	if kafkaWriter != nil {
		kafkaWriter.Close()
		log.Println("Kafka Producer closed.")
	}
}

func CreateOrUpdateKafka(product *models.Product) error {
	jsonData, err := json.Marshal(product)
	if err != nil {
		return err
	}

	_, err = kafka.DialLeader(context.Background(), "tcp", "kafka:29092", "go-to-rails", 0)
	if err != nil {
		log.Printf("Error creating Kafka topic: %v", err)
	}

	return kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Value: jsonData,
	})
}
