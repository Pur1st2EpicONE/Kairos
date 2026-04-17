package v1

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/helpers"
)

// ConfirmBooking handles POST /api/v1/events/:id/confirm.
// It expects a JWT in the Authorization header to identify the user,
// and an event ID in the URL path. It calls the service layer to confirm
// the user's booking for that event. On success, it returns a 200 OK with
// status "confirmed". On failure, it responds with the appropriate error.
func (h *Handler) ConfirmBooking(c *ginext.Context) {

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

	if err := h.service.ConfirmBooking(c.Request.Context(), userID, eventID); err != nil {
		RespondError(c, err)
		return
	}

	respondOK(c, models.StatusConfirmed)

}
