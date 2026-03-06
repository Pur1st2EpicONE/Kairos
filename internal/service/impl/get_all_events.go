package impl

import (
	"Kairos/internal/models"
	"context"
)

func (c *CoreService) GetAllEvents(ctx context.Context) []models.Event {
	events, err := c.storage.GetAllEvents(ctx)
	if err != nil {
		c.logger.LogError("service — failed to get events from DB", err, "layer", "service.impl")
	}
	return events
}
