package postgres

import (
	"context"
	"database/sql"
)

// UpdateBookingStatus changes the status of a booking (e.g., to 'confirmed' or 'expired').
// It executes the update within the provided transaction.
func (c *CoreStorage) UpdateBookingStatus(tx *sql.Tx, ctx context.Context, bookingID int64, status string) error {

	_, err := tx.ExecContext(ctx, `

	UPDATE bookings
	SET status = $1
	WHERE id = $2`,

		status, bookingID)

	return err

}
