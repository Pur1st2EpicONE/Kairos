// Package v1 implements the JSON API handlers for version 1 of the Kairos API.
// It provides endpoints for authentication, event management, and booking operations.
package v1

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/helpers"
)

// CreateBooking handles POST /api/v1/events/:id/book.
// It extracts the user ID from the JWT context and the event ID from the URL,
// then calls the service layer to create a booking for that user and event.
// On success it returns the booking ID. On failure it responds with an appropriate error.
func (h *Handler) CreateBooking(c *ginext.Context) {

	userID, ok := c.Request.Context().Value(models.UserIDKey).(int64)
	if !ok {
		RespondError(c, errs.ErrInvalidToken)
		return
	}

	eventID := c.Param("id")
	if err := helpers.ParseUUID(eventID); err != nil {
		RespondError(c, errs.ErrInvalidEventID)
		return
	}

	bookingID, err := h.service.CreateBooking(c.Request.Context(), userID, eventID)
	if err != nil {
		RespondError(c, err)
		return
	}

	respondOK(c, bookingID)

}
