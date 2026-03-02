package postgres

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

// SetStatus updates the status of a notification in the database.
// If the status is "canceled", only notifications that are currently
// "pending" or "running late" can be canceled. The check for already canceled
// notifications is done directly in the database rather than in the service layer,
// because a separate query would be needed anyway to get the current status.
// This way, each request results in a single query: try to update, and if the DB
// reports zero affected rows, return the appropriate error.
func (c *CoreStorage) SetStatus(ctx context.Context, notificationID string, status string) error {

	query := `
    
	UPDATE Notifications
    SET status = $1, updated_at = NOW()
    WHERE uuid = $2;`

	var args []any
	args = append(args, status, notificationID)

	if status == models.StatusCanceled {
		query = query[:len(query)-1] + " AND status IN ($3, $4);"
		args = append(args, models.StatusPending, models.StatusLate)
	}

	res, err := c.db.ExecWithRetry(ctx, retry.Strategy{
		Attempts: c.config.QueryRetryStrategy.Attempts,
		Delay:    c.config.QueryRetryStrategy.Delay,
		Backoff:  c.config.QueryRetryStrategy.Backoff}, query, args...)

	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of affected rows: %w", err)
	}

	if rows == 0 {
		if status == models.StatusCanceled {
			return errs.ErrCannotCancel
		}
		return errs.ErrNotificationNotFound
	}

	return nil

}
