package postgres

import (
	"Kairos/internal/models"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

// GetInfo retrieves public event details (title, description, date, available seats)
// by its UUID, without locking. Uses the configured query retry strategy.
// Returns the event or an error.
func (c *CoreStorage) GetInfo(ctx context.Context, eventUUID string) (*models.Event, error) {

	row, err := c.db.QueryRowWithRetry(ctx, retry.Strategy(c.config.QueryRetryStrategy), `

    SELECT title, description, event_date, available_seats
    FROM events
    WHERE uuid = $1`,

		eventUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to query row %w", err)
	}

	event := new(models.Event)
	err = row.Scan(
		&event.Title,
		&event.Description,
		&event.Date,
		&event.Seats)
	if err != nil {
		return nil, err
	}

	return event, nil

}
