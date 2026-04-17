package v1

import "time"

// RegisterDTO represents the request body for user registration.
type RegisterDTO struct {
	Login    string `json:"login"`    // User's login name
	Password string `json:"password"` // User's password (plain text, will be hashed by service)
}

// LoginDTO represents the request body for user login.
type LoginDTO struct {
	Login    string `json:"login"`    // User's login name
	Password string `json:"password"` // User's password
}

// CreateEventDTO represents the request body for creating a new event.
type CreateEventDTO struct {
	ID          string `json:"event_id"`    // Optional: event ID
	Title       string `json:"title"`       // Event title
	Description string `json:"description"` // Event description
	Date        string `json:"date"`        // Event date in RFC3339 format
	Seats       int    `json:"seats"`       // Total number of available seats
	BookingTTL  string `json:"booking_ttl"` // Time-to-live for bookings (e.g., "30m", "1h")
}

// InfoResponseDTO represents the response body for event information.
type InfoResponseDTO struct {
	ID          string    `json:"event_id"`    // Event unique identifier
	Title       string    `json:"title"`       // Event title
	Description string    `json:"description"` // Event description
	Date        time.Time `json:"date"`        // Event date and time (UTC)
	Seats       int       `json:"seats"`       // Remaining seats
}
