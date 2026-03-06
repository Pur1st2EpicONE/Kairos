package postgres

import (
	"Kairos/internal/models"
	"context"
	"database/sql"
	"time"
)

func (c *CoreStorage) GetEventForBooking(tx *sql.Tx, ctx context.Context, eventUUID string) (*models.Event, error) {

	row := tx.QueryRowContext(ctx, `

    SELECT id, userid, uuid, title, description, event_date, available_seats, booking_ttl
    FROM events
    WHERE uuid = $1
    FOR UPDATE`,

		eventUUID)
	event := new(models.Event)
	var ttlSeconds int

	err := row.Scan(
		&event.DBID,
		&event.UserID,
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Date,
		&event.Seats,
		&ttlSeconds)
	if err != nil {
		return nil, err
	}

	event.BookingTTL = time.Duration(ttlSeconds) * time.Second
	return event, nil

}
