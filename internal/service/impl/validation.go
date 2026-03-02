package impl

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"strings"
	"time"
)

func validateCreate(event *models.Event) error {

	if err := validateTitle(event.Title); err != nil {
		return err
	}

	if err := validateDescription(event.Description); err != nil {
		return err
	}

	if err := validateDate(event.Date); err != nil {
		return err
	}

	if err := validateSeats(event.TotalSeats); err != nil {
		return err
	}

	return nil
}

func validateDescription(desc string) error {
	if len(desc) > 2000 {
		return errs.ErrDescriptionTooLong
	}
	return nil
}

func validateSeats(seats int) error {
	if seats <= 0 {
		return errs.ErrInvalidSeatCount
	}
	if seats > 10000 {
		return errs.ErrTooManySeats
	}
	return nil
}

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

func validateDate(t time.Time) error {

	if t.IsZero() {
		return errs.ErrMissingDate
	}

	now := time.Now().UTC()

	if t.Before(now) {
		return errs.ErrDateInPast
	}
	if t.Before(now.Add(5 * time.Minute)) {
		return errs.ErrDateTooSoon
	}

	if t.After(now.AddDate(1, 0, 0)) {
		return errs.ErrDateTooFar
	}

	return nil

}
