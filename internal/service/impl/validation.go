package impl

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"strings"
	"time"
)

const (
	minBookingTTL        = 1 * time.Minute // Minimum allowed booking TTL
	maxBookingTTL        = 24 * time.Hour  // Maximum allowed booking TTL
	maxSeats             = 10000           // Maximum number of seats for an event
	maxDescriptionLength = 2000            // Maximum length of event description
)

// validateUser checks that the user's login and password are not empty.
// Returns ErrEmptyLogin or ErrEmptyPassword accordingly.
func validateUser(user models.User) error {
	if user.Login == "" {
		return errs.ErrEmptyLogin
	}
	if user.Password == "" {
		return errs.ErrEmptyPassword
	}
	return nil
}

// validateEvent validates all fields of an event.
// It checks title, description, date, seats, and booking TTL.
// Returns the first validation error encountered.
func validateEvent(event *models.Event) error {

	if err := validateTitle(event.Title); err != nil {
		return err
	}

	if err := validateDescription(event.Description); err != nil {
		return err
	}

	if err := validateDate(event.Date); err != nil {
		return err
	}

	if err := validateSeats(event.Seats); err != nil {
		return err
	}

	if err := validateBookingTTL(event.BookingTTL); err != nil {
		return err
	}

	return nil
}

// validateTitle checks that the title is not empty, has at least 3 characters,
// and at most 200 characters after trimming spaces.
func validateTitle(title string) error {

	title = strings.TrimSpace(title)

	if title == "" {
		return errs.ErrMissingTitle
	}

	if len(title) < 3 {
		return errs.ErrTitleTooShort
	}

	if len(title) > 200 {
		return errs.ErrTitleTooLong
	}

	return nil

}

// validateDescription ensures the description does not exceed the maximum length.
// An empty description is allowed.
func validateDescription(desc string) error {
	if len(desc) > maxDescriptionLength {
		return errs.ErrDescriptionTooLong
	}
	return nil
}

// validateDate checks that the event date:
// - is provided (not zero)
// - is not in the past
// - is at least 24 hours from now
// - is not more than one year in the future
// Returns appropriate errors for each violation.
func validateDate(t time.Time) error {

	if t.IsZero() {
		return errs.ErrMissingDate
	}

	now := time.Now().UTC()

	if t.Before(now) {
		return errs.ErrDateInPast
	}
	if t.Before(now.Add(24 * time.Hour)) {
		return errs.ErrDateTooSoon
	}

	if t.After(now.AddDate(1, 0, 0)) {
		return errs.ErrDateTooFar
	}

	return nil

}

// validateSeats checks that the number of seats is positive and does not exceed maxSeats.
func validateSeats(seats int) error {
	if seats <= 0 {
		return errs.ErrInvalidSeatCount
	}
	if seats > maxSeats {
		return errs.ErrTooManySeats
	}
	return nil
}

// validateBookingTTL checks that the booking time-to-live is between
// minBookingTTL and maxBookingTTL (inclusive).
func validateBookingTTL(ttl time.Duration) error {

	if ttl < minBookingTTL {
		return errs.ErrBookingTTLTooShort
	}

	if ttl > maxBookingTTL {
		return errs.ErrBookingTTLTooLong
	}

	return nil

}
