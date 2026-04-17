package postgres

import (
	"Kairos/internal/models"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

// GetAllEvents retrieves all events from the database, ordered by creation time (ascending).
// It returns a slice of Event structs containing public fields (uuid, title, description,
// event_date, available_seats). Returns an error if the query or row scanning fails.
func (c *CoreStorage) GetAllEvents(ctx context.Context) ([]models.Event, error) {

	var events []models.Event

	query := `
	
	SELECT uuid, title, description, event_date, available_seats
	FROM events
	ORDER BY created_at ASC;`

	rows, err := c.db.QueryWithRetry(ctx, retry.Strategy(c.config.QueryRetryStrategy), query)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Date, &event.Seats); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		events = append(events, event)
	}

	return events, nil

}
