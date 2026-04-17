// Package rabbitmq implements a message broker using RabbitMQ.
// It handles exchange/queue declaration, publishing with delayed expiration,
// and consuming with configurable workers and retry strategies.
package rabbitmq

import (
	"Kairos/internal/config"
	"Kairos/internal/logger"

	"context"
	"fmt"

	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	mainExchange = "mainExchange"     // mainExchange is the direct exchange used for all booking messages.
	contentType  = "application/json" // contentType is the MIME type for published messages.
	exchangeKind = "direct"           // exchangeKind is the type of the main exchange.
)

// Broker represents the RabbitMQ-based message broker.
// It holds a client, a consumer, a producer, and a cancellation callback.
type Broker struct {
	logger     logger.Logger                                    // structured logger
	config     config.Broker                                    // broker configuration
	Consumer   *rabbitmq.Consumer                               // message consumer
	producer   *rabbitmq.Publisher                              // message publisher
	cancelFunc func(ctx context.Context, bookingID int64) error // callback for cancellation
	client     *rabbitmq.RabbitClient                           // underlying RabbitMQ client
}

// SetCancelFunc assigns the function that will be called when a cancellation
// message is delivered. It is typically set after the service layer is constructed.
func (b *Broker) SetCancelFunc(fn func(ctx context.Context, bookingID int64) error) {
	b.cancelFunc = fn
}

// NewBroker creates a new RabbitMQ broker, establishes a connection,
// declares the main exchange and the queue, and initialises a consumer
// and a publisher. It returns the broker or an error if any step fails.
func NewBroker(logger logger.Logger, config config.Broker, cancelFunc func(ctx context.Context, bookingID int64) error) (*Broker, error) {

	client, err := rabbitmq.NewClient(rabbitmq.ClientConfig{

		URL:            config.URL,
		ConnectionName: config.ConnectionName,
		ConnectTimeout: config.ConnectTimeout,

		ReconnectStrat: retry.Strategy{
			Attempts: config.Reconnect.Attempts,
			Delay:    config.Reconnect.Delay,
			Backoff:  config.Reconnect.Backoff},

		ConsumingStrat: retry.Strategy{
			Attempts: config.Consumer.Attempts,
			Delay:    config.Reconnect.Delay,
			Backoff:  config.Reconnect.Backoff}})

	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}

	err = client.DeclareExchange(mainExchange, exchangeKind, true, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	err = client.DeclareQueue(config.QueueName, mainExchange, config.QueueName, true, false, true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	producer := rabbitmq.NewPublisher(client, mainExchange, contentType)

	b := &Broker{
		logger:     logger,
		config:     config,
		Consumer:   nil,
		producer:   producer,
		cancelFunc: cancelFunc,
		client:     client}

	b.Consumer = rabbitmq.NewConsumer(client, rabbitmq.ConsumerConfig{
		Queue:         config.QueueName,
		ConsumerTag:   config.Consumer.ConsumerTag,
		AutoAck:       config.Consumer.AutoAck,
		Ask:           rabbitmq.AskConfig{Multiple: false},
		Nack:          rabbitmq.NackConfig{Multiple: false, Requeue: true},
		Args:          amqp.Table{},
		Workers:       config.Consumer.Workers,
		PrefetchCount: config.Consumer.PrefetchCount,
	}, func(ctx context.Context, msg amqp.Delivery) error { return b.handler(ctx, msg) })

	return b, nil

}

// Shutdown closes the underlying RabbitMQ client connection.
// It logs an error if the closure fails, otherwise logs a successful shutdown.
func (b *Broker) Shutdown() {
	if err := b.client.Close(); err != nil {
		b.logger.LogError("rabbit — failed to shutdown gracefully", err, "layer", "broker.rabbitMQ")
	} else {
		b.logger.LogInfo("rabbit — shutdown complete", "layer", "broker.rabbitMQ")
	}
}
