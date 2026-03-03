package v1

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateEvent(c *gin.Context) {

	userIDRaw := c.Request.Context().Value("userID")
	userID, ok := userIDRaw.(int64)
	if !ok {
		RespondError(c, errs.ErrInvalidToken)
		return
	}

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

	bookingTTL, err := time.ParseDuration(request.BookingTTL)
	if err != nil {
		RespondError(c, errs.ErrInvalidBookingTTL)
		return
	}

	eventId, err := h.service.CreateEvent(c.Request.Context(), &models.Event{
		UserID:      userID,
		Title:       request.Title,
		Description: request.Description,
		Date:        date,
		Seats:       request.Seats,
		BookingTTL:  bookingTTL})

	if err != nil {
		RespondError(c, err)
		return
	}

	respondOK(c, eventId)

}
