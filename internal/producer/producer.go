package producer

import "OrderService/internal/models"

type Producer interface {
	SendOrder(order *models.Order) error
	Close() error
}
