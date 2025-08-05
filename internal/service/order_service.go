package service

import (
	"errors"
	"fmt"
	"time"

	"OrderService/internal/db"
	"OrderService/internal/models"
)

type OrderService struct{}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) ProcessOrder(order *models.Order) error {
	if order.OrderUID == "" {
		return errors.New("OrderUID не должен быть пустым")
	}

	if order.DateCreated.IsZero() {
		return errors.New("DateCreated должен быть задан")
	}
	if order.DateCreated.After(time.Now().UTC()) {
		return errors.New("DateCreated не может быть в будущем")
	}

	requiredStrings := map[string]string{
		"TrackNumber": order.TrackNumber,
		"Entry":       order.Entry,
		"Locale":      order.Locale,
		"CustomerID":  order.CustomerID,
		"DeliverySvc": order.DeliverySvc,
		"ShardKey":    order.ShardKey,
		"OofShard":    order.OofShard,
	}
	for field, val := range requiredStrings {
		if val == "" {
			return fmt.Errorf("Поле %s обязательно для заполнения", field)
		}
	}

	if order.SMID == 0 {
		return errors.New("SMID должен быть задан и отличаться от 0")
	}

	d := order.Delivery
	if d.Name == "" || d.Phone == "" || d.Zip == "" || d.City == "" || d.Address == "" || d.Region == "" || d.Email == "" {
		return errors.New("Все поля доставки должны быть заполнены")
	}

	p := order.Payment
	if p.Transaction == "" || p.Currency == "" || p.Provider == "" || p.Amount < 0 || p.PaymentDT == 0 ||
		p.Bank == "" || p.DeliveryCost < 0 || p.GoodsTotal < 0 || p.CustomFee < 0 {
		return errors.New("Некорректные данные оплаты")
	}

	if len(order.Items) == 0 {
		return errors.New("В заказе должен быть хотя бы один товар")
	}
	for i, item := range order.Items {
		if item.ChrtID == 0 {
			return fmt.Errorf("Item[%d]: ChrtID должен быть задан", i)
		}
		if item.TrackNumber == "" {
			return fmt.Errorf("Item[%d]: TrackNumber обязателен", i)
		}
		if item.Price < 0 {
			return fmt.Errorf("Item[%d]: Price не может быть отрицательным", i)
		}
		if item.Rid == "" || item.Name == "" || item.Size == "" || item.Brand == "" {
			return fmt.Errorf("Item[%d]: поля Rid, Name, Size и Brand обязательны", i)
		}
		if item.TotalPrice < 0 {
			return fmt.Errorf("Item[%d]: TotalPrice не может быть отрицательным", i)
		}
		if item.NMID == 0 {
			return fmt.Errorf("Item[%d]: NMID должен быть задан", i)
		}
	}

	var totalItemsPrice int
	for _, item := range order.Items {
		totalItemsPrice += item.TotalPrice
	}
	if totalItemsPrice != p.GoodsTotal {
		return errors.New("Сумма total_price товаров не совпадает с goods_total в оплате")
	}

	if err := db.SaveOrder(order); err != nil {
		return err
	}

	return nil
}
