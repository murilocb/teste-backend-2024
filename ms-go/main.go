package main

import (
	"log"
	consumer "ms-go/app/consumers"
	"ms-go/app/services/products"
	_ "ms-go/db"
	"ms-go/router"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	products.InitKafkaProducer()
	defer products.CloseKafkaProducer()

	go func() {
		consumer.CreateOrUpdateKafka()
	}()

	go func() {
		router.Run()
	}()
	<-sigs

	log.Println("Encerrando o serviço...")
}
