package impl

import (
	"Kairos/internal/models"
	"context"

	"github.com/wb-go/wbf/helpers"
)

func (c *CoreService) CreateEvent(ctx context.Context, event *models.Event) (string, error) {

	if err := validateCreate(event); err != nil {
		return "", err
	}

	initialize(event)

	if err := c.storage.CreateEvent(ctx, event); err != nil {
		c.logger.LogError("service — failed to create event", err, "layer", "service.impl")
		return "", err
	}

	return event.ID, nil

}

func initialize(event *models.Event) {
	event.ID = helpers.CreateUUID()
}
