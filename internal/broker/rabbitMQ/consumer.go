package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	wbf "github.com/wb-go/wbf/rabbitmq"
)

// Consume starts the consumer's event loop. It blocks until the consumer
// stops normally (client closed or context cancelled) and logs any unexpected
// error that is not due to cancellation.
func (b *Broker) Consume() {
	if err := b.Consumer.Start(b.client.Context()); err != nil &&
		!errors.Is(err, wbf.ErrClientClosed) && !errors.Is(err, context.Canceled) {
		b.logger.LogError("consumer — unexpected context error", err, "layer", "broker.rabbitMQ")
	}
}

// handler is the callback function invoked for each delivered message.
// It unmarshals the booking ID from the JSON body and calls the injected
// cancelFunc. If cancelFunc is not set or unmarshalling fails, an error is returned.
func (b *Broker) handler(ctx context.Context, msg amqp091.Delivery) error {

	var bookingID int64

	if err := json.Unmarshal(msg.Body, &bookingID); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	if b.cancelFunc == nil {
		return fmt.Errorf("cancelFunc is not set")
	}

	return b.cancelFunc(ctx, bookingID)

}
