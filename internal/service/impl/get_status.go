package impl

import (
	"Kairos/internal/errs"
	"context"
	"errors"
)

// GetStatus retrieves the current status of a notification.
// It first checks the cache, and if not found, falls back to the database.
// After fetching from the database, the status is updated in the cache for future calls.
// This ensures a single source of truth while optimizing for read performance.
func (c *CoreService) GetStatus(ctx context.Context, notificationID string) (string, error) {

	status, err := c.storage.GetStatus(ctx, notificationID)
	if err != nil {
		if errors.Is(err, errs.ErrNotificationNotFound) {
			c.logger.Debug("service — notification status fetched from DB", "notificationID", notificationID, "layer", "service.impl")
		}
		c.logger.LogError("service — failed to get notification status from DB", err, "notificationID", notificationID, "layer", "service.impl")
		return "", err
	}

	return status, nil

}
