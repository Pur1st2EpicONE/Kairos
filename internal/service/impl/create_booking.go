package impl

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

// CreateBooking creates a new booking for a user and event.
// It runs within a transaction: retrieves the event (with row lock), checks seat availability,
// creates a pending booking, publishes a delayed cancellation message to the broker,
// sends a Telegram notification asynchronously, and decrements the available seats.
func (c *CoreService) CreateBooking(ctx context.Context, userID int64, eventID string) (int64, error) {

	var bookingID int64

	err := c.storage.Transaction(ctx, func(tx *sql.Tx, ctx context.Context) error {

		event, err := c.storage.GetEventForBooking(tx, ctx, eventID)
		if err != nil {
			return c.wrap(err)
		}
		if event.Seats <= 0 {
			return errs.ErrEventFull
		}

		var innerErr error
		booking := initBooking(userID, event.DBID, event.BookingTTL)

		bookingID, innerErr = c.storage.CreateBooking(tx, ctx, booking)
		if innerErr != nil {
			return c.wrap(innerErr)
		}

		booking.ID = bookingID

		if err := c.broker.Produce(booking); err != nil {
			return c.wrap(err)
		}

		go func() {
			if err := c.notifier.Notify(models.Notification{Channel: models.Telegram, Message: models.StatusCreated}); err != nil {
				c.logger.LogError("service — failed to send booking expiration notification", err, "layer", "service.impl")
			}
		}()

		if err := c.storage.UpdateEventSeats(tx, ctx, false, event.DBID); err != nil {
			c.logger.LogError("service — failed to update event seats", err, "layer", "service.impl")
		}

		return nil

	})

	if err != nil {
		return 0, err
	}
	return bookingID, nil

}

// initBooking creates a new Booking struct with pending status and expiration
// time based on the event's booking TTL. It does not set the booking ID.
func initBooking(userID int64, eventID int64, bookingTTL time.Duration) *models.Booking {
	now := time.Now()
	return &models.Booking{
		UserID:    userID,
		EventID:   eventID,
		Status:    models.StatusPending,
		CreatedAt: now,
		ExpiresAt: now.Add(bookingTTL),
	}
}

// wrap converts common database errors into domain errors:
// - sql.ErrNoRows -> ErrEventNotFound
// - PostgreSQL unique violation (23505) -> ErrBookingAlreadyExists
// Other errors are logged and returned as is.
func (c *CoreService) wrap(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return errs.ErrEventNotFound
	}
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return errs.ErrBookingAlreadyExists
	}
	c.logger.LogError("service — failed to create booking", err, "layer", "service.impl")
	return err
}
