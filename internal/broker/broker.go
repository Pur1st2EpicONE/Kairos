// Package broker provides abstractions for message brokers used in the application.
// It defines the Broker interface and exposes a constructor for concrete implementations.
package broker

import (
	rabbitmq "Kairos/internal/broker/rabbitMQ"
	"Kairos/internal/config"
	"Kairos/internal/logger"
	"Kairos/internal/models"
	"Kairos/internal/notifier"
	"Kairos/internal/repository"
)

type Broker interface {
	Consume()
	Produce(booking *models.Booking) error
	Shutdown()
}

func NewBroker(logger logger.Logger, config config.Broker, storage *repository.Storage, notifier notifier.Notifier) (Broker, error) {
	return rabbitmq.NewBroker(logger, config, storage, notifier)
}
