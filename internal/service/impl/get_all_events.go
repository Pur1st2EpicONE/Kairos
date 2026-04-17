package impl

import (
	"Kairos/internal/models"
	"context"
)

// GetAllEvents retrieves all events from the storage layer.
// If an error occurs, it is logged and an empty slice is returned.
// This method is used for rendering the home page.
func (c *CoreService) GetAllEvents(ctx context.Context) []models.Event {
	events, err := c.storage.GetAllEvents(ctx)
	if err != nil {
		c.logger.LogError("service — failed to get events from DB", err, "layer", "service.impl")
	}
	return events
}
