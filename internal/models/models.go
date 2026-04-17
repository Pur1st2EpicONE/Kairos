// Package models defines the core data structures used across the Kairos service.
// It includes domain entities (User, Event, Booking), status constants, and context keys.
package models

import "time"

// User represents an authenticated user of the system.
// It stores login credentials and the internal database ID.
type User struct {
	ID       int64  // Unique database identifier
	Login    string `json:"login"`    // User's login name (unique)
	Password string `json:"password"` // Plain-text password (hashed before storage)
}

// Event represents a bookable event created by a user.
// It contains event details, capacity, and the time window for booking confirmations.
type Event struct {
	DBID        int64         // Internal database primary key
	ID          string        // Public unique identifier (UUID)
	UserID      int64         // ID of the user who created the event
	Title       string        // Event title
	Description string        // Event description
	Date        time.Time     // Scheduled date and time (UTC)
	Seats       int           // Total number of available seats
	BookingTTL  time.Duration // Time-to-live for a booking before it expires
}

// Booking represents a user's reservation for a specific event.
// It tracks the booking status, creation time, and expiration deadline.
type Booking struct {
	ID        int64     // Unique database identifier
	UserID    int64     // ID of the user who made the booking
	EventID   int64     // ID of the booked event (references Event.DBID)
	Status    string    // Current status: pending, confirmed, or expired
	CreatedAt time.Time // Timestamp when the booking was created
	ExpiresAt time.Time // Deadline after which the booking is automatically cancelled
}

// Notification represents a message to be sent via a specific channel.
// It is used by the notifier component to deliver alerts (e.g., Telegram).
type Notification struct {
	Channel string `json:"channel"` // Delivery channel (e.g., "telegram")
	Message string `json:"message"` // Message content
}

// userIDContextKey is an unexported type used to avoid context key collisions.
type userIDContextKey struct{}

// UserIDKey is the context key for storing the authenticated user's ID.
// It is used to pass the user ID from authentication middleware to handlers.
var UserIDKey = userIDContextKey{}

// Standard booking status constants used throughout the application.
const (
	StatusCreated   = "created"   // Booking was just created
	StatusPending   = "pending"   // Booking created but not yet confirmed
	StatusConfirmed = "confirmed" // Booking confirmed by the user
	StatusExpired   = "expired"   // Booking expired without confirmation
	StatusCancled   = "cancelled" // Booking was just cancelled
)

// Notification channel constants.
const (
	Telegram = "telegram" // Send notifications via Telegram bot
)
