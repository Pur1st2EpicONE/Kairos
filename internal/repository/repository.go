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

type AuthStorage interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
}

type CoreStorage interface {
	Transact(ctx context.Context, fn func(tx *sql.Tx, ctx context.Context) error) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEventForBooking(tx *sql.Tx, ctx context.Context, eventUUID string) (*models.Event, error)
	CreateBooking(tx *sql.Tx, ctx context.Context, booking *models.Booking) (int64, error)
	UpdateEventSeats(tx *sql.Tx, ctx context.Context, eventDBID int64) error
	GetBookingForConfirm(tx *sql.Tx, ctx context.Context, userID int64, eventUUID string) (*models.Booking, error)
	UpdateBookingStatus(tx *sql.Tx, ctx context.Context, bookingID int64, status string) error
	Close()
}

type Storage struct {
	AuthStorage
	CoreStorage
}

func NewStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) *Storage {
	return &Storage{
		AuthStorage: postgres.NewAuthStorage(logger, config, db),
		CoreStorage: postgres.NewCoreStorage(logger, config, db),
	}
}

// ConnectDB establishes a connection to the Postgres database using the provided configuration.
// It returns a dbpg.DB instance ready for queries.
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
