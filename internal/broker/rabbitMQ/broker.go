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
	mainExchange = "mainExchange"
	contentType  = "application/json"
	exchangeKind = "direct"
)

type Broker struct {
	logger     logger.Logger
	config     config.Broker
	Consumer   *rabbitmq.Consumer
	producer   *rabbitmq.Publisher
	cancelFunc func(ctx context.Context, bookingID int64) error
	client     *rabbitmq.RabbitClient
}

func (b *Broker) SetCancelFunc(fn func(ctx context.Context, bookingID int64) error) {
	b.cancelFunc = fn
}

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

func (b *Broker) Shutdown() {
	if err := b.client.Close(); err != nil {
		b.logger.LogError("rabbit — failed to shutdown gracefully", err, "layer", "broker.rabbitMQ")
	} else {
		b.logger.LogInfo("rabbit — shutdown complete", "layer", "broker.rabbitMQ")
	}
}
