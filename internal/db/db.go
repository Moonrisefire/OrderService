package db

import (
	"log"

	"OrderService/internal/cache"
	"OrderService/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := "host=localhost user=order_user password=password dbname=Orders port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	err = DB.AutoMigrate(&models.Order{}, &models.Delivery{}, &models.Payment{}, &models.Item{})
	if err != nil {
		log.Fatal("Не удалось выполнить миграцию:", err)
	}
}

func LoadCache(c *cache.Cache) {
	var orders []models.Order
	if err := DB.Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders).Error; err != nil {
		log.Printf("Ошибка загрузки заказов из БД для кэша: %v", err)
		return
	}
	for _, o := range orders {
		c.Set(o.OrderUID, o)
	}
}

func SaveOrder(order *models.Order) error {
	err := DB.Create(order).Error
	if err != nil {
		log.Printf("Ошибка при сохранении заказа %s: %v", order.OrderUID, err)
		return err
	}
	return nil
}
