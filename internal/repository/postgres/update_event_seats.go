package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

func (c *CoreStorage) UpdateEventSeats(tx *sql.Tx, ctx context.Context, increment bool, eventID int64) error {

	op := "-"
	if increment {
		op = "+"
	}

	_, err := tx.ExecContext(ctx, fmt.Sprintf(`
        
	UPDATE events
    SET available_seats = available_seats %s 1
    WHERE id = $1`,

		op), eventID)

	return err

}
