package rabbitmq

import (
	"Kairos/internal/models"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/retry"
)

// Produce publishes a booking cancellation message with a delayed expiration.
// It creates a temporary queue for the booking ID with a TTL equal to the time
// until the booking expires. The message is published to that queue, which
// automatically expires after the TTL plus an additional buffer.
// If a queue for the same booking already exists (PreconditionFailed), it
// skips redeclaration and returns nil. The operation is retried according to
// the producer retry strategy.
func (b *Broker) Produce(booking *models.Booking) error {

	sendAt := max(time.Until(booking.ExpiresAt), 0)

	return retry.DoContext(b.client.Context(), retry.Strategy{
		Attempts: b.config.Producer.Attempts,
		Delay:    b.config.Producer.Delay,
		Backoff:  b.config.Producer.Backoff}, func() error {

		queueArgs := amqp.Table{
			"x-message-ttl":             int64(sendAt.Milliseconds()),
			"x-dead-letter-exchange":    mainExchange,
			"x-dead-letter-routing-key": b.config.QueueName,
			"x-expires":                 int64(sendAt.Milliseconds() + b.config.Producer.MessageQueueTTL.Milliseconds()),
		}

		err := b.client.DeclareQueue(strconv.FormatInt(booking.ID, 10), mainExchange, strconv.FormatInt(booking.ID, 10), false, true, true, queueArgs)
		if err != nil {
			if amqpErr, ok := err.(*amqp.Error); ok && amqpErr.Code == amqp.PreconditionFailed { // exception 406
				b.logger.Debug("producer — booking is already in queue, skipping",
					"bookingID", booking.ID, "layer", "broker.rabbitMQ")
				return nil
			}
			return fmt.Errorf("failed to declare queue: %w", err)
		}

		ch, err := b.client.GetChannel()
		if err != nil {
			return fmt.Errorf("failed to get channel: %w", err)
		}
		defer func() { _ = ch.Close() }()

		body, err := json.Marshal(booking.ID)
		if err != nil {
			return fmt.Errorf("failed to marshal bookingID to json: %w", err)
		}

		pub := amqp.Publishing{ContentType: contentType, Body: body}

		if err := ch.PublishWithContext(b.client.Context(), mainExchange, strconv.FormatInt(booking.ID, 10), false, false, pub); err != nil {
			return fmt.Errorf("failed to publish with context: %w", err)
		}

		return nil

	})

}
