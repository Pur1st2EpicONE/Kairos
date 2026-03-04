package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	wbf "github.com/wb-go/wbf/rabbitmq"
)

func (b *Broker) Consume() {

	if err := b.Consumer.Start(b.client.Context()); err != nil &&
		!errors.Is(err, wbf.ErrClientClosed) && !errors.Is(err, context.Canceled) {
		b.logger.LogError("broker — consumer returned unexpected context error", err, "layer", "broker.rabbimq")
	}

}

func (b *Broker) handler(ctx context.Context, msg amqp091.Delivery) error {
	_ = ctx
	var bookingID int64

	if err := json.Unmarshal(msg.Body, &bookingID); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	fmt.Println("SUCCESS", bookingID)

	return nil

}

// updateStatus updates the notification status in both cache and storage.
// It applies automatic transformations: Pending → Sent, Late or timed-out → FailedToSendInTime.
// Updates are retried according to the configured retry strategy.
func (b *Broker) updateStatus(ctx context.Context, notificationID string, sendAt time.Time, status string) {

}
