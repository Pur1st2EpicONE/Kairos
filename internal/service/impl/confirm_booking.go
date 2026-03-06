package impl

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"context"
	"database/sql"
	"errors"
)

func (c *CoreService) ConfirmBooking(ctx context.Context, userID int64, eventID string) error {

	return c.storage.Transaction(ctx, func(tx *sql.Tx, ctx context.Context) error {

		booking, err := c.storage.GetBookingForConfirm(tx, ctx, userID, eventID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errs.ErrBookingNotFound
			}
			c.logger.LogError("service — failed to get booking for confirmation", err, "layer", "service.impl")
			return err
		}

		if booking.Status == models.StatusConfirmed {
			return errs.ErrAlreadyConfirmed
		}

		if booking.Status == models.StatusExpired {
			return errs.ErrBookingExpired
		}

		if err := c.storage.UpdateBookingStatus(tx, ctx, booking.ID, models.StatusConfirmed); err != nil {
			c.logger.LogError("service — failed to update booking status", err, "layer", "service.impl")
			return err
		}

		return nil

	})

}
