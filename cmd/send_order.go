package main

import (
	"encoding/json"
	"log"
	"os"

	"OrderService/internal/kafka"
	"OrderService/internal/models"
)

func main() {
	producer := kafka.NewProducer([]string{"localhost:9092"}, "orders")
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Ошибка при закрытии Kafka producer: %v", err)
		}
	}()

	data, err := os.ReadFile("model3.json")
	if err != nil {
		log.Fatal("Ошибка чтения файла:", err)
	}

	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		log.Fatal("Ошибка парсинга JSON:", err)
	}

	if err := producer.SendOrder(&order); err != nil {
		log.Fatal("Ошибка отправки в Kafka:", err)
	}

	log.Println("Заказ успешно отправлен в Kafka")
}
