package postgres

import (
	"Kairos/internal/models"
	"context"

	"github.com/wb-go/wbf/retry"
)

// Cleanup removes outdated notifications from the database
// based on retention rules for each notification status.
func (c *CoreStorage) Cleanup(ctx context.Context) {

	query := `
	
        DELETE FROM Notifications 
        WHERE (status = $1 AND updated_at < NOW() - $2 * INTERVAL '1 second')
        OR (status = $3 AND updated_at < NOW() - $4 * INTERVAL '1 second')
        OR (status IN ($5, $6) AND updated_at < NOW() - $7 * INTERVAL '1 second');`

	_, err := c.db.ExecWithRetry(ctx, retry.Strategy{
		Attempts: c.config.QueryRetryStrategy.Attempts,
		Delay:    c.config.QueryRetryStrategy.Delay,
		Backoff:  c.config.QueryRetryStrategy.Backoff}, query,

		models.StatusCanceled, int(c.config.RetentionStrategy.Canceled.Seconds()),
		models.StatusSent, int(c.config.RetentionStrategy.Completed.Seconds()),
		models.StatusFailed, models.StatusFailedToSendInTime, int(c.config.RetentionStrategy.Failed.Seconds()),
	)

	if err != nil {
		c.logger.LogError("postgres — failed to delete old notifications", err, "layer", "repository.postgres")
	} else {
		c.logger.Debug("postgres — old notifications cleaned", "layer", "repository.postgres")
	}

}
