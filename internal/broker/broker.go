package broker

import (
	rabbitmq "Kairos/internal/broker/rabbitMQ"
	"Kairos/internal/config"
	"Kairos/internal/logger"
	"Kairos/internal/models"
	"context"
)

type Broker interface {
	Consume()
	Produce(booking *models.Booking) error
	SetCancelFunc(fn func(ctx context.Context, bookingID int64) error)
	Shutdown()
}

func NewBroker(logger logger.Logger, config config.Broker, cancelFunc func(ctx context.Context, bookingID int64) error) (Broker, error) {
	return rabbitmq.NewBroker(logger, config, cancelFunc)
}
