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

// Broker defines the interface for a message broker used by the application.
// It supports consuming messages, producing notifications, and graceful shutdown.
type Broker interface {
	Consume() error                                 // Consume starts processing messages from the broker.
	Produce(notification models.Notification) error // Produce sends a notification message to the broker.
	Shutdown()                                      // Shutdown gracefully stops the broker and releases resources.
}

// NewBroker creates a new Broker instance. Currently, it returns a RabbitMQ-based broker.
// It initializes the broker with the provided logger, configuration, cache, storage, and notifier.
func NewBroker(logger logger.Logger, config config.Broker, storage *repository.Storage, notifier notifier.Notifier) (Broker, error) {
	return rabbitmq.NewBroker(logger, config, storage, notifier)
}
