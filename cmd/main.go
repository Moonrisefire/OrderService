package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"OrderService/config"
	"OrderService/internal/api"
	"OrderService/internal/cache"
	"OrderService/internal/db"
	"OrderService/internal/kafka"
	"OrderService/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	db.Init()

	c := cache.NewCache()
	go c.Cleanup(5 * time.Minute)

	db.LoadCache(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := kafka.EnsureTopicExists([]string{config.KafkaBroker}, config.KafkaTopic); err != nil {
		log.Fatalf("Не удалось создать Kafka топик: %v", err)
	}

	producer := kafka.NewProducer([]string{config.KafkaBroker}, config.KafkaTopic)
	defer func() {
		if err := producer.Close(); err != nil {
			log.Println("Ошибка при закрытии Kafka producer:", err)
		}
	}()

	orderService := service.NewOrderService()

	consumer := kafka.NewConsumer(
		[]string{config.KafkaBroker},
		config.KafkaTopic,
		"order-service",
		orderService,
		c,
	)
	go consumer.Start(ctx)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Handle("/*", http.FileServer(http.Dir("internal/frontend")))
	r.Get("/order/{id}", api.GetOrderHandler(c))
	r.Post("/orders", api.CreateOrderHandler(orderService))

	srv := &http.Server{
		Addr:    ":" + config.HTTPPort,
		Handler: r,
	}

	go func() {
		log.Printf("HTTP сервер запущен на порту :%s", config.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершаем работу...")

	cancel()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	log.Println("Сервер завершил работу корректно")
}
