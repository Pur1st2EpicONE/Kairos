// Package broker defines the message broker abstraction for asynchronous
// booking cancellation tasks. It provides an interface and a constructor
// that returns a RabbitMQ-based implementation.
package broker

import (
	rabbitmq "Kairos/internal/broker/rabbitMQ"
	"Kairos/internal/config"
	"Kairos/internal/logger"
	"Kairos/internal/models"
	"context"
)

// Broker defines the contract for a message broker that produces
// booking cancellation messages and consumes them to trigger cancellations.
type Broker interface {
	Consume()                                                          // Consume starts the consumer loop, processing incoming messages until the broker is shut down or the context is cancelled.
	Produce(booking *models.Booking) error                             // Produce publishes a booking to the broker with a delayed expiration queue.
	SetCancelFunc(fn func(ctx context.Context, bookingID int64) error) // SetCancelFunc sets the callback function that will be invoked when a cancellation message is received.
	Shutdown()                                                         // Shutdown gracefully closes the broker connection and stops all consumers and producers.
}

// NewBroker constructs a Broker instance using RabbitMQ as the underlying
// transport. It accepts a cancelFunc that will be wired later via SetCancelFunc.
// If the connection or exchange/queue declaration fails, an error is returned.
func NewBroker(logger logger.Logger, config config.Broker, cancelFunc func(ctx context.Context, bookingID int64) error) (Broker, error) {
	return rabbitmq.NewBroker(logger, config, cancelFunc)
}
