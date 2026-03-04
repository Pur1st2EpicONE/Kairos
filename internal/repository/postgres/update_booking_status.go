package postgres

import (
	"context"
	"database/sql"
)

func (c *CoreStorage) UpdateBookingStatus(tx *sql.Tx, ctx context.Context, bookingID int64, status string) error {

	_, err := tx.ExecContext(ctx, `

	UPDATE bookings
	SET status = $1
	WHERE id = $2`,

		status, bookingID)

	return err

}
