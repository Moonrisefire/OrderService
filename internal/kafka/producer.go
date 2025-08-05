package kafka

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"time"

	"OrderService/internal/models"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
			Async:        false,
		},
	}
}

func (p *Producer) SendOrder(order *models.Order) error {
	order.DateCreated = order.DateCreated.UTC()

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(order.OrderUID),
		Value: data,
		Time:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return err
	}

	log.Printf("Заказ %s отправлен в Kafka %s\n", order.OrderUID, p.writer.Topic)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func EnsureTopicExists(brokers []string, topic string) error {
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	return controllerConn.CreateTopics(topicConfig)
}
