package postgres

import (
	"Kairos/internal/models"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

// MarkLates updates notifications that are past their scheduled send time
// from "pending" to "running late" and returns their IDs.
// It is used by sysmon, the consumer goroutine that monitors broker health,
// so if the broker is not healthy, users can see which notifications are delayed.
func (c *CoreStorage) MarkLates(ctx context.Context) ([]string, error) {

	query := `

        UPDATE Notifications
        SET status = $1, updated_at = NOW()
        WHERE status = $2 AND send_at < NOW()
        RETURNING uuid;`

	rows, err := c.db.QueryWithRetry(ctx, retry.Strategy{
		Attempts: c.config.QueryRetryStrategy.Attempts,
		Delay:    c.config.QueryRetryStrategy.Delay,
		Backoff:  c.config.QueryRetryStrategy.Backoff},
		query, models.StatusLate, models.StatusPending)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var uuids []string

	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		uuids = append(uuids, uuid)
	}

	return uuids, nil

}
