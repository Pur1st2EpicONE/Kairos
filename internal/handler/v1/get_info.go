package v1

import (
	"Kairos/internal/errs"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/helpers"
)

func (h *Handler) GetInfo(c *gin.Context) {

	eventIDStr := c.Param("id")
	if err := helpers.ParseUUID(eventIDStr); err != nil {
		RespondError(c, errs.ErrInvalidEventID)
		return
	}

	event, err := h.service.GetInfo(c.Request.Context(), eventIDStr)
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
