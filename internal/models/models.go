package models

import "time"

type User struct {
	ID       int64
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Event struct {
	DBID        int64
	ID          string
	UserID      int64
	Title       string
	Description string
	Date        time.Time
	Seats       int
	BookingTTL  time.Duration
}

type Booking struct {
	ID        int64
	UserID    int64
	EventID   int64
	Status    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Notification struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

type userIDContextKey struct{}

var UserIDKey = userIDContextKey{}

const (
	Created = "Booking created"
	Cancled = "Booking canceled"
)

const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusExpired   = "expired"
)

const (
	Telegram = "telegram"
)
