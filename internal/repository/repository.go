// Package repository provides data access abstractions and a concrete storage
// implementation backed by PostgreSQL. It defines interfaces for authentication
// and core business operations, and exposes a composite Storage struct.
package repository

import (
	"Kairos/internal/config"
	"Kairos/internal/logger"
	"Kairos/internal/models"
	"Kairos/internal/repository/postgres"
	"context"
	"database/sql"
	"fmt"

	"github.com/wb-go/wbf/dbpg"
)

// AuthStorage defines the data operations for user authentication:
// creating a new user and retrieving a user by login.
type AuthStorage interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)       // CreateUser stores a new user in the database. It returns the generated user ID.
	GetUserByLogin(ctx context.Context, login string) (models.User, error) // GetUserByLogin retrieves a user by their unique login name.
}

// CoreStorage defines the data operations for events and bookings,
// including transactional support and query methods for the core domain.
type CoreStorage interface {
	Transaction(ctx context.Context, fn func(tx *sql.Tx, ctx context.Context) error) error                         // Transaction executes a function within a database transaction.
	CreateEvent(ctx context.Context, event *models.Event) error                                                    // CreateEvent persists a new event record.
	GetEventForBooking(tx *sql.Tx, ctx context.Context, eventUUID string) (*models.Event, error)                   // GetEventForBooking retrieves an event by its public UUID for booking operations.
	CreateBooking(tx *sql.Tx, ctx context.Context, booking *models.Booking) (int64, error)                         // CreateBooking inserts a new booking record and returns its ID.
	UpdateEventSeats(tx *sql.Tx, ctx context.Context, increment bool, eventID int64) error                         // UpdateEventSeats increments or decrements the available seats of an event.
	GetBookingForConfirm(tx *sql.Tx, ctx context.Context, userID int64, eventUUID string) (*models.Booking, error) // GetBookingForConfirm retrieves a pending booking for a user and event.
	UpdateBookingStatus(tx *sql.Tx, ctx context.Context, bookingID int64, status string) error                     // UpdateBookingStatus changes the status of a booking (e.g., to confirmed or expired).
	CancelBooking(tx *sql.Tx, ctx context.Context, bookingID int64) (int64, error)                                 // CancelBooking sets a booking's status to cancelled and returns the associated event ID.
	GetInfo(ctx context.Context, eventUUID string) (*models.Event, error)                                          // GetInfo returns event details by its public UUID (read-only).
	GetAllEvents(ctx context.Context) ([]models.Event, error)                                                      // GetAllEvents returns all events (used for listing on the home page).
	Close()                                                                                                        // Close releases the underlying database connection pool.
}

// Storage composes both AuthStorage and CoreStorage, providing a unified
// data access object for the service layer.
type Storage struct {
	AuthStorage
	CoreStorage
}

// NewStorage constructs a Storage instance by creating the PostgreSQL-specific
// implementations of AuthStorage and CoreStorage. It returns the composed Storage.
func NewStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) *Storage {
	return &Storage{
		AuthStorage: postgres.NewAuthStorage(logger, config, db),
		CoreStorage: postgres.NewCoreStorage(logger, config, db),
	}
}

// ConnectDB establishes a connection to the PostgreSQL database using the
// provided configuration. It builds the DSN, applies connection pool settings,
// and verifies connectivity with a ping. Returns the DB handle or an error.
func ConnectDB(config config.Storage) (*dbpg.DB, error) {

	options := &dbpg.Options{
		MaxOpenConns:    config.MaxOpenConns,
		MaxIdleConns:    config.MaxIdleConns,
		ConnMaxLifetime: config.ConnMaxLifetime,
	}

	db, err := dbpg.New(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode), nil, options)
	if err != nil {
		return nil, fmt.Errorf("database driver not found or DSN invalid: %w", err)
	}

	if err := db.Master.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil

}
