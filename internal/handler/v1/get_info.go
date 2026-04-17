package v1

import (
	"Kairos/internal/errs"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/helpers"
)

// GetInfo handles GET /api/v1/events/:id.
// It validates the event ID from the URL, fetches event details via the service,
// and returns the event information as InfoResponseDTO. Returns 404 if not found.
func (h *Handler) GetInfo(c *ginext.Context) {

	eventID := c.Param("id")
	if err := helpers.ParseUUID(eventID); err != nil {
		RespondError(c, errs.ErrInvalidEventID)
		return
	}

	event, err := h.service.GetInfo(c.Request.Context(), eventID)
	if err != nil {
		RespondError(c, err)
		return
	}

	respondOK(c, InfoResponseDTO{
		Title:       event.Title,
		Description: event.Description,
		Date:        event.Date,
		Seats:       event.Seats})

}
