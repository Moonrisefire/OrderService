package config

import "os"

var (
	KafkaBroker = getEnv("KAFKA_BROKER", "localhost:9092")
	KafkaTopic  = getEnv("KAFKA_TOPIC", "orders")

	PostgresDSN = getEnv("POSTGRES_DSN",
		"host=localhost user=order_user password=password dbname=orders_db port=5432 sslmode=disable")

	HTTPPort = getEnv("HTTP_PORT", "8080")
)

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
