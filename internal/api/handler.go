package api

import (
	"encoding/json"
	"net/http"

	"OrderService/internal/cache"
	"OrderService/internal/db"
	"OrderService/internal/models"

	"github.com/go-chi/chi/v5"
)

func GetOrderHandler(c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "id не указан", http.StatusBadRequest)
			return
		}

		var order models.Order
		var found bool

		if cachedOrder, ok := c.Get(id); ok {
			order = cachedOrder
			found = true
		} else {
			res := db.DB.Preload("Delivery").Preload("Payment").Preload("Items").
				Where("order_uid = ?", id).
				First(&order)

			if res.Error != nil {
				http.Error(w, "не найден", http.StatusNotFound)
				return
			}

			c.Set(id, order)
			found = true
		}

		if found {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(order); err != nil {
				http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
			}
		}
	}
}
