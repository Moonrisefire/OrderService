package api

import (
	"encoding/json"
	"log"
	"net/http"

	"OrderService/internal/models"
	"OrderService/internal/service"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func CreateOrderHandler(service *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order models.Order

		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Невалидный JSON: "+err.Error(), http.StatusBadRequest)
			log.Printf("Ошибка декодирования JSON: %v", err)
			return
		}

		if err := validate.Struct(order); err != nil {
			http.Error(w, "Ошибка валидации: "+err.Error(), http.StatusBadRequest)
			log.Printf("Ошибка валидации заказа: %v", err)
			return
		}

		if err := service.ProcessOrder(&order); err != nil {
			http.Error(w, "Ошибка при отправке в Kafka: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Ошибка отправки заказа в Kafka: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		if err := json.NewEncoder(w).Encode(map[string]string{
			"status":    "accepted",
			"order_uid": order.OrderUID,
		}); err != nil {
			log.Printf("Ошибка кодирования ответа JSON: %v", err)
			http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
			return
		}
	}
}
