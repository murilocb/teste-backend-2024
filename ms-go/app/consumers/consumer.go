package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ms-go/app/models"
	"ms-go/app/services/products"

	"github.com/segmentio/kafka-go"
)

func CreateOrUpdateKafka() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:29092"},
		Topic:     "rails-to-go",
		Partition: 0,
		MaxBytes:  10e6,
	})
	defer reader.Close()

	ctx := context.Background()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled. Exiting...")
			return
		default:
			err := handleMessage(ctx, reader)
			if err != nil {
				fmt.Printf("Error handling message: %v. Retrying in 5 seconds...\n", err)
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func handleMessage(ctx context.Context, reader *kafka.Reader) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	msg, err := reader.ReadMessage(ctx)
	if err != nil {
		return fmt.Errorf("error reading message: %w", err)
	}

	fmt.Println("Received message:", string(msg.Value))

	var product models.Product
	if err := json.Unmarshal(msg.Value, &product); err != nil {
		return fmt.Errorf("error decoding product: %w", err)
	}

	if err := product.Validate(); err != nil {
		return fmt.Errorf("error validating product: %w", err)
	}

	existingProduct, _ := products.Details(product)

	if existingProduct == nil {
		fmt.Println("Product not found. Creating...")
		createdProduct, createErr := products.Create(product, false)
		if createErr != nil {
			return fmt.Errorf("error creating product: %w", createErr)
		}
		fmt.Println("Product created successfully:", createdProduct)
	} else {
		fmt.Println("Product found. Updating...")
		updatedProduct, updateErr := products.Update(product, false)
		if updateErr != nil {
			return fmt.Errorf("error updating product: %w", updateErr)
		}
		fmt.Println("Product updated successfully:", updatedProduct)
	}

	return nil
}
