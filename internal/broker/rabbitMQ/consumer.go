package rabbitmq

import (
	"Kairos/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	wbf "github.com/wb-go/wbf/rabbitmq"
)

// Consume starts the RabbitMQ consumer and system monitor in the background.
// It returns an error if the consumer fails to start for reasons other than client closure or context cancellation.
func (b *Broker) Consume() error {

	if err := b.Consumer.Start(b.client.Context()); err != nil &&
		!errors.Is(err, wbf.ErrClientClosed) && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil

}

// handler processes a single RabbitMQ delivery message.
// It unmarshals the JSON payload into a Notification, checks its status,
// attempts to send it via the notifier, and updates the status accordingly.
func (b *Broker) handler(ctx context.Context, msg amqp091.Delivery) error {

	var notification models.Notification

	if err := json.Unmarshal(msg.Body, &notification); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	status, err := b.storage.GetStatus(ctx, notification.ID)
	if err != nil {
		return err
	}

	if status != models.StatusCanceled {

		if err := b.notifier.Notify(notification); err != nil {
			if status != models.StatusFailed {
				b.updateStatus(ctx, notification.ID, notification.SendAt, models.StatusFailed)
			}
			return err
		}

		b.updateStatus(ctx, notification.ID, notification.SendAt, status)

	}

	return nil

}

// updateStatus updates the notification status in both cache and storage.
// It applies automatic transformations: Pending → Sent, Late or timed-out → FailedToSendInTime.
// Updates are retried according to the configured retry strategy.
func (b *Broker) updateStatus(ctx context.Context, notificationID string, sendAt time.Time, status string) {

}
