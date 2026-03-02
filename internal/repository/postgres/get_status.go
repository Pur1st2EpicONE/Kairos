package postgres

import (
	"Kairos/internal/errs"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

// GetStatus returns the current status of a notification by its ID.
func (c *CoreStorage) GetStatus(ctx context.Context, notificationID string) (string, error) {

	query := `

    SELECT status
    FROM Notifications
    WHERE uuid = $1;`

	row, err := c.db.QueryRowWithRetry(ctx, retry.Strategy{
		Attempts: c.config.QueryRetryStrategy.Attempts,
		Delay:    c.config.QueryRetryStrategy.Delay,
		Backoff:  c.config.QueryRetryStrategy.Backoff}, query, notificationID)

	if err != nil {
		return "", fmt.Errorf("failed to execute query: %w", err)
	}

	var status string
	if err := row.Scan(&status); err != nil {
		return "", errs.ErrNotificationNotFound
	}

	return status, nil

}
