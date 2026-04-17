package v1

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"time"

	"github.com/wb-go/wbf/ginext"
)

// CreateEvent handles POST /api/v1/events.
// It extracts the authenticated user ID from the context, binds the JSON request
// to CreateEventDTO, validates and parses the date and booking TTL, then calls
// the service layer to create a new event. Returns the new event ID on success.
func (h *Handler) CreateEvent(c *ginext.Context) {

	userID, ok := c.Request.Context().Value(models.UserIDKey).(int64)
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

	eventID, err := h.service.CreateEvent(c.Request.Context(), &models.Event{
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

	respondOK(c, eventID)

}
