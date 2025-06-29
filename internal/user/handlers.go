package user

// обработчики событий на шине сообщений

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func UserRegisteredHandler(service *Service) func(amqp.Delivery) {
	return func(msg amqp.Delivery) {
		var payload struct {
			UserID string `json:"userId"`
		}

		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			log.Printf("userRegistered msg parse error: %v", err)
			msg.Nack(false, true)
			return
		}

		user := &User{
			ID:       uuid.MustParse(payload.UserID),
			Username: "test",
		}

		if err := service.CreateUser(context.Background(), user); err != nil {
			log.Printf("userRegistered user creation error: %v", err)
			msg.Nack(false, true)
			return
		}

		msg.Ack(false)
	}
}
