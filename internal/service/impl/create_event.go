package impl

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"context"

	"github.com/lib/pq"
	"github.com/wb-go/wbf/helpers"
)

// CreateEvent creates a new event after validation and initialisation.
// It validates the event fields (title, date, seats, booking TTL), generates a UUID,
// and persists the event via the storage layer. Returns the generated event UUID
// or an error (e.g., ErrInvalidUserID if the user does not exist).
func (c *CoreService) CreateEvent(ctx context.Context, event *models.Event) (string, error) {

	if err := validateEvent(event); err != nil {
		return "", err
	}

	initialize(event)

	if err := c.storage.CreateEvent(ctx, event); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23503" {
				return "", errs.ErrInvalidUserID
			}
		}
		c.logger.LogError("service — failed to create event", err, "layer", "service.impl")
		return "", err
	}

	return event.ID, nil

}

// initialize sets the event's ID to a new UUID.
func initialize(event *models.Event) {
	event.ID = helpers.CreateUUID()
}
