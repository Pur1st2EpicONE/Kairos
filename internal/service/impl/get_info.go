package impl

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"context"
	"database/sql"
	"errors"
)

// GetInfo retrieves public event information by its UUID.
// It returns the event or an error. If the event is not found, it returns ErrEventNotFound.
// Other storage errors are logged and returned as is.
func (c *CoreService) GetInfo(ctx context.Context, eventID string) (*models.Event, error) {

	event, err := c.storage.GetInfo(ctx, eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrEventNotFound
		}
		c.logger.LogError("service — failed to get event info from storage", err, "layer", "service.impl")
		return nil, err
	}

	return event, nil

}
