package postgres

import (
	"Kairos/internal/models"
	"context"
	"time"
)

func (c *CoreStorage) GetInfo(ctx context.Context, eventUUID string) (*models.Event, error) {

	event := new(models.Event)
	var ttlSeconds int

	row := c.db.Master.QueryRowContext(ctx, `

    SELECT title, description, event_date, available_seats
    FROM events
    WHERE uuid = $1`,

		eventUUID)
	err := row.Scan(
		&event.Title,
		&event.Description,
		&event.Date,
		&event.Seats)
	if err != nil {
		return nil, err
	}

	event.BookingTTL = time.Duration(ttlSeconds) * time.Second
	return event, nil

}
