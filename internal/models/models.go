// Package models defines the core data structures used in Chronos,
// including notifications and their associated constants.
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
	Channel string   `json:"channel"` // Delivery channel
	Subject string   `json:"subject"` // Subject or title of the notification
	Message string   `json:"message"` // Main content of the notification
	SendTo  []string `json:"send_to"` // List of recipients
}

const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusExpired   = "expired"
)

const (
	Email    = "email"    // Email channel
	Stdout   = "stdout"   // Standard output
	Telegram = "telegram" // Telegram bot channel
)

const (
	MaxEmailLength   = 254 // Maximum length for email addresses
	MaxSubjectLength = 254 // Maximum length for email subject
	MaxMessageLength = 254 // Maximum length for message content
)
