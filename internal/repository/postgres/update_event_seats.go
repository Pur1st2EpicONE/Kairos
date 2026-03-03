package postgres

import (
	"context"
	"database/sql"
)

func (c *CoreStorage) UpdateEventSeats(tx *sql.Tx, ctx context.Context, eventDBID int64) error {
	_, err := tx.ExecContext(ctx, `
        UPDATE events 
        SET available_seats = available_seats - 1 
        WHERE id = $1
    `, eventDBID)
	return err
}
