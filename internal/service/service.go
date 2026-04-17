// Package service defines the business logic layer interfaces and provides
// a composite Service struct that aggregates authentication and core operations.
package service

import (
	"Kairos/internal/broker"
	"Kairos/internal/config"
	"Kairos/internal/logger"
	"Kairos/internal/models"
	"Kairos/internal/notifier"
	"Kairos/internal/repository"
	"Kairos/internal/service/impl"
	"context"
)

// AuthService defines the operations for user authentication and token management.
type AuthService interface {
	CreateUser(ctx context.Context, user models.User) (int64, error) // CreateUser registers a new user. Returns the generated user ID or an error.
	CreateToken(userID int64) (string, error)                        // CreateToken generates a JWT token for the given user ID.
	GetUserId(ctx context.Context, user models.User) (int64, error)  // GetUserId authenticates a user by login/password and returns the user ID.
	ParseToken(tokenString string) (int64, error)                    // ParseToken validates a JWT token and extracts the user ID.
}

// CoreService defines the business operations for events and bookings.
type CoreService interface {
	CreateEvent(ctx context.Context, event *models.Event) (string, error)           // CreateEvent creates a new event. Returns the event's public UUID or an error.
	CreateBooking(ctx context.Context, userID int64, eventID string) (int64, error) // CreateBooking creates a booking for a user and event. Returns the booking ID.
	ConfirmBooking(ctx context.Context, userID int64, eventID string) error         // ConfirmBooking confirms a pending booking for a user and event.
	CancelBooking(ctx context.Context, bookingID int64) error                       // CancelBooking cancels a booking by its ID (used by broker on expiration).
	GetInfo(ctx context.Context, eventID string) (*models.Event, error)             // GetInfo retrieves public event details by its UUID.
	GetAllEvents(ctx context.Context) []models.Event                                // GetAllEvents returns all events (for listing on the home page).
}

// Service composes AuthService and CoreService into a single struct,
// providing a unified dependency for the handler layer.
type Service struct {
	AuthService
	CoreService
}

// NewService constructs a Service instance by creating the concrete
// implementations of AuthService and CoreService with their dependencies.
func NewService(logger logger.Logger, config config.Service, broker broker.Broker, storage *repository.Storage, notifier notifier.Notifier) *Service {
	return &Service{
		AuthService: impl.NewAuthService(logger, config, storage.AuthStorage),
		CoreService: impl.NewCoreService(logger, broker, storage.CoreStorage, notifier),
	}
}
