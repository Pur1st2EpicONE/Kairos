package postgres

import (
	"Kairos/internal/models"
	"context"
)

func (c *CoreStorage) CreateEvent(ctx context.Context, event *models.Event) error {

	_, err := c.db.ExecContext(ctx, `

	INSERT INTO events (uuid, title, description, event_date, total_seats)
	VALUES ($1, $2, $3, $4, $5)`,

		event.ID, event.Title, event.Description, event.Date, event.TotalSeats)
	if err != nil {
		return err
	}

	return nil

}
