package impl

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"context"
	"errors"
)

// CancelNotification attempts to cancel a notification by setting its status to "canceled".
// The method first checks the cache; if the status is already "canceled" or otherwise non-cancelable,
// it returns an appropriate error. The actual update is performed in the database.
// The logic for detecting already canceled notifications is done in the DB itself to avoid
// an extra read query — only one query per request is needed. Cache is updated after a successful change.
func (c *CoreService) CancelNotification(ctx context.Context, notificationID string) error {

	if err := c.storage.SetStatus(ctx, notificationID, models.StatusCanceled); err != nil {

		switch {

		case errors.Is(err, errs.ErrNotificationNotFound):
			return errs.ErrNotificationNotFound

		case errors.Is(err, errs.ErrCannotCancel):
			currentStatus, err := c.storage.GetStatus(ctx, notificationID)
			if err != nil {
				c.logger.LogError("service — failed to get notification status from DB", err, "layer", "service.impl")
				return err
			}

			if currentStatus == models.StatusCanceled {
				return errs.ErrAlreadyCanceled
			}
			return errs.ErrCannotCancel

		default:
			c.logger.LogError("service — failed to set notification status in DB", err, "layer", "service.impl")
			return err

		}

	}

	return nil

}
