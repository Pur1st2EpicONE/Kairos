package impl

import (
	"Kairos/internal/models"
	"context"
)

// GetAllStatuses retrieves all notifications with their current status and scheduled times.
// This method is intended purely for frontend purposes and is not optimized for high-volume usage.
// Errors are logged but not returned to the caller, since the frontend can tolerate partial failures.
func (c *CoreService) GetAllStatuses(ctx context.Context) []models.Notification {
	statuses, err := c.storage.GetAllStatuses(ctx)
	if err != nil {
		c.logger.LogError("service — failed to get notification statuses from DB", err, "layer", "service.impl")
	}
	return statuses
}
