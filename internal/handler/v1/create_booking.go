package v1

import (
	"Kairos/internal/errs"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/helpers"
)

func (h *Handler) CreateBooking(c *gin.Context) {

	userIDRaw := c.Request.Context().Value("userID")
	userID, ok := userIDRaw.(int64)
	if !ok {
		RespondError(c, errs.ErrInvalidToken)
		return
	}

	eventIDStr := c.Param("id")
	if err := helpers.ParseUUID(eventIDStr); err != nil {
		RespondError(c, errs.ErrInvalidEventID)
		return
	}

	bookingID, err := h.service.CreateBooking(c.Request.Context(), userID, eventIDStr)
	if err != nil {
		RespondError(c, err)
		return
	}

	respondOK(c, bookingID)

}
