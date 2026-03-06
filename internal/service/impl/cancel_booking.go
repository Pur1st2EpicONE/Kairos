package impl

import (
	"Kairos/internal/models"
	"context"
	"database/sql"
	"errors"
)

const cancled = "Booking canceled"

func (c *CoreService) CancelBooking(ctx context.Context, bookingID int64) error {

	return c.storage.Transaction(ctx, func(tx *sql.Tx, ctx context.Context) error {

		eventID, err := c.storage.CancelBooking(tx, ctx, bookingID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.logger.Debug("booking confirmed, cancelation skipped", "bookingID", bookingID, "layer", "service.impl")
				return nil
			}
			c.logger.LogError("service — failed to cancel booking", err, "layer", "service.impl")
			return err
		}

		if err := c.storage.UpdateEventSeats(tx, ctx, true, eventID); err != nil {
			c.logger.LogError("service — failed to increment seats", err, "layer", "service.impl")
			return err
		}

		go func() {
			if err := c.notifier.Notify(models.Notification{Channel: models.Telegram, Message: cancled}); err != nil {
				c.logger.LogError("service — failed to send booking expiration notification", err, "layer", "service.impl")
			}
		}()

		c.logger.Debug("service — expired booking was canceled", "layer", "service.impl")

		return nil

	})

}
