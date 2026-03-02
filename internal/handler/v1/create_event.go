package v1

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"

	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) CreateEvent(c *ginext.Context) {

	var request CreateEventDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		RespondError(c, errs.ErrInvalidJSON)
		return
	}

	date, err := parseTime(request.Date)
	if err != nil {
		RespondError(c, err)
		return
	}

	eventId, err := h.service.CreateEvent(c.Request.Context(), &models.Event{
		Title:       request.Title,
		Description: request.Description,
		Date:        date,
		TotalSeats:  request.TotalSeats})

	if err != nil {
		RespondError(c, err)
		return
	}

	respondOK(c, eventId)

}
