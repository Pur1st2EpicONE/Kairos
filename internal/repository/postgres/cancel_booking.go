package postgres

import (
	"context"
	"database/sql"
)

// CancelBooking updates a pending booking's status to 'expired' and returns
// the associated event ID. It expects a transaction and a booking ID.
// If no row is updated (e.g., booking not found or already not pending),
// it returns sql.ErrNoRows.
func (c *CoreStorage) CancelBooking(tx *sql.Tx, ctx context.Context, bookingID int64) (int64, error) {

	var eventID int64
	err := tx.QueryRowContext(ctx, `

	UPDATE bookings
	SET status = 'expired'
	WHERE id = $1 AND status = 'pending'
	RETURNING event_id`,

		bookingID).Scan(&eventID)
	if err != nil {
		return 0, err
	}

	return eventID, nil

}
