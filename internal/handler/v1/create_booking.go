package v1

import (
	"Kairos/internal/errs"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/helpers"
)

func (h *Handler) CreateBooking(c *ginext.Context) {

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

	bookingID, err := h.service.CreateBooking(c.Request.Context(), userID, eventID)
	if err != nil {
		RespondError(c, err)
		return
	}

	respondOK(c, bookingID)

}
