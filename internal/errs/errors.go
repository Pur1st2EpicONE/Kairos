package errs

import "errors"

var (
	ErrInvalidJSON           = errors.New("invalid JSON format")
	ErrInvalidNotificationID = errors.New("missing or invalid notification ID")
	ErrMissingChannel        = errors.New("channel is required")
	ErrUnsupportedChannel    = errors.New("unsupported channel")
	ErrMissingDate           = errors.New("date is required")
	ErrInvalidDate           = errors.New("invalid date format, expected RFC3339")
	ErrDateInPast            = errors.New("date cannot be in the past")
	ErrDateTooFar            = errors.New("date is too far in the future")
	ErrDateTooSoon           = errors.New("date is too soon")

	ErrMissingTitle  = errors.New("title is required")
	ErrTitleTooShort = errors.New("title is too short")
	ErrTitleTooLong  = errors.New("title is too long")

	ErrDescriptionTooLong = errors.New("description is too long")

	ErrInvalidSeatCount = errors.New("total seats must be greater than zero")
	ErrTooManySeats     = errors.New("total seats exceeds maximum allowed")

	ErrMissingSendTo       = errors.New("send_to is required")
	ErrInvalidEmailFormat  = errors.New("invalid email format")
	ErrMissingEmailSubject = errors.New("email subject is required")
	ErrEmailSubjectTooLong = errors.New("email subject is too long")
	ErrRecipientTooLong    = errors.New("recipient exceeds maximum length")

	ErrNotificationNotFound = errors.New("notification with given ID not found")
	ErrAlreadyCanceled      = errors.New("notification is already canceled")
	ErrCannotCancel         = errors.New("notification cannot be canceled in its current state")

	ErrInternal             = errors.New("internal server error")
	ErrUrgentDeliveryFailed = errors.New("cannot schedule notification for immediate delivery — service is temporarily unavailable")

	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid login or password")

	ErrEmptyAuthHeader      = errors.New("empty auth header")
	ErrInvalidAuthHeader    = errors.New("invalid auth header")
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidEventID       = errors.New("empty or invalid event_id")
	ErrEventFull            = errors.New("no seats left for the event")
	ErrEventNotFound        = errors.New("event not found")
	ErrBookingAlreadyExists = errors.New("booking already exists")
	ErrInvalidUserID        = errors.New("invalid user_id: user does not exist")
	ErrInvalidBookingTTL    = errors.New("invalid booking expiration time")
)
