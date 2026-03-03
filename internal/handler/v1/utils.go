package v1

import (
	"Kairos/internal/errs"
	"errors"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
)

// parseTime parses a string in RFC3339 format into a UTC time.Time value.
// Returns ErrMissingSendAt if the string is empty, or ErrInvalidSendAt if parsing fails.
func parseTime(timeStr string) (time.Time, error) {

	if timeStr == "" {
		return time.Time{}, errs.ErrMissingDate
	}

	validTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, errs.ErrInvalidDate
	}

	return validTime.UTC(), nil

}

// respondOK sends a JSON HTTP 200 response with the given payload.
func respondOK(c *ginext.Context, response any) {
	c.JSON(http.StatusOK, ginext.H{"result": response})
}

// respondError maps an error to an HTTP status and sends it as a JSON response.
// Uses mapErrorToStatus to determine the appropriate status code.
func RespondError(c *ginext.Context, err error) {
	if err != nil {
		status, msg := mapErrorToStatus(err)
		c.AbortWithStatusJSON(status, ginext.H{"error": msg})
	}
}

// mapErrorToStatus converts a known error to an appropriate HTTP status code and message.
// Returns 400 for validation errors, 404 for not found, and 500 for internal errors.
func mapErrorToStatus(err error) (int, string) {

	switch {

	case errors.Is(err, errs.ErrInvalidJSON),
		errors.Is(err, errs.ErrInvalidNotificationID),
		errors.Is(err, errs.ErrMissingChannel),
		errors.Is(err, errs.ErrUnsupportedChannel),
		errors.Is(err, errs.ErrMissingDate),
		errors.Is(err, errs.ErrInvalidDate),
		errors.Is(err, errs.ErrDateInPast),
		errors.Is(err, errs.ErrDateTooFar),
		errors.Is(err, errs.ErrDateTooSoon),
		errors.Is(err, errs.ErrMissingSendTo),
		errors.Is(err, errs.ErrMissingEmailSubject),
		errors.Is(err, errs.ErrEmailSubjectTooLong),
		errors.Is(err, errs.ErrInvalidEmailFormat),
		errors.Is(err, errs.ErrCannotCancel),
		errors.Is(err, errs.ErrAlreadyCanceled),
		errors.Is(err, errs.ErrRecipientTooLong),
		errors.Is(err, errs.ErrMissingTitle),
		errors.Is(err, errs.ErrTitleTooShort),
		errors.Is(err, errs.ErrTitleTooLong),
		errors.Is(err, errs.ErrDescriptionTooLong),
		errors.Is(err, errs.ErrInvalidSeatCount),
		errors.Is(err, errs.ErrInvalidUserID),
		errors.Is(err, errs.ErrInvalidBookingTTL),
		errors.Is(err, errs.ErrTooManySeats):

		return http.StatusBadRequest, err.Error()

	case errors.Is(err, errs.ErrNotificationNotFound):
		return http.StatusNotFound, err.Error()

	case errors.Is(err, errs.ErrUserAlreadyExists),
		errors.Is(err, errs.ErrBookingAlreadyExists):
		return http.StatusConflict, err.Error()

	case errors.Is(err, errs.ErrInvalidToken),
		errors.Is(err, errs.ErrEmptyAuthHeader),
		errors.Is(err, errs.ErrInvalidAuthHeader),
		errors.Is(err, errs.ErrInvalidCredentials):
		return http.StatusUnauthorized, err.Error()

	default:
		if errors.Is(err, errs.ErrUrgentDeliveryFailed) {
			return http.StatusInternalServerError, err.Error()
		}
		return http.StatusInternalServerError, errs.ErrInternal.Error()
	}

}
