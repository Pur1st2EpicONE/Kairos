package v1

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/helpers"
)

func (h *Handler) ConfirmBooking(c *ginext.Context) {

	userID, ok := c.Request.Context().Value("userID").(int64)
	if !ok {
		RespondError(c, errs.ErrInvalidToken)
		return
	}

	eventID := c.Param("id")
	if err := helpers.ParseUUID(eventID); err != nil {
		RespondError(c, errs.ErrInvalidEventID)
		return
	}

	if err := h.service.ConfirmBooking(c.Request.Context(), userID, eventID); err != nil {
		RespondError(c, err)
		return
	}

	respondOK(c, models.StatusConfirmed)

}
