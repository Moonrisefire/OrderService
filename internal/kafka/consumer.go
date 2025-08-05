package kafka

import (
	"context"
	"encoding/json"
	"log"

	"OrderService/internal/cache"
	"OrderService/internal/models"
	"OrderService/internal/service"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader       *kafka.Reader
	orderService *service.OrderService
	cache        *cache.Cache
}

func NewConsumer(brokers []string, topic, groupID string, s *service.OrderService, c *cache.Cache) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		orderService: s,
		cache:        c,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("Kafka consumer запущен")
	defer func() {
		if err := c.reader.Close(); err != nil {
			log.Println("Ошибка закрытия Kafka reader:", err)
		}
		log.Println("Kafka consumer остановлен")
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Контекст отменён — выходим из consumer-а")
			return
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Println("Ошибка чтения сообщения из Kafka:", err)
				continue
			}

			log.Printf("Получено сообщение от Kafka: key=%s value=%s\n", string(m.Key), string(m.Value))

			var order models.Order
			if err := json.Unmarshal(m.Value, &order); err != nil {
				log.Println("Ошибка парсинга JSON:", err)
				continue
			}

			if err := c.orderService.ProcessOrder(&order); err != nil {
				log.Println("Ошибка обработки заказа:", err)
				continue
			}

			c.cache.Set(order.OrderUID, order)
			log.Println("Заказ сохранён и закэширован:", order.OrderUID)
		}
	}
}
