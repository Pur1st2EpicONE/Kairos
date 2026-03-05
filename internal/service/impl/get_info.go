package impl

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (c *CoreService) GetInfo(ctx context.Context, eventID string) (*models.Event, error) {

	event, err := c.storage.GetInfo(ctx, eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println(err)
			return nil, errs.ErrEventNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return event, nil

}
