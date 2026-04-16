package postgres_test

import (
	"Kairos/internal/config"
	"Kairos/internal/logger"
	"Kairos/internal/models"
	"Kairos/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	wbf "github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/dbpg"
)

var (
	testStorage *repository.Storage
	testDB      *dbpg.DB
)

const (
	validEventUUID  = "550e8400-e29b-41d4-a716-446655440000"
	nonExistentUUID = "11111111-1111-1111-1111-111111111111"
)

var fixedEventDate = time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

func TestMain(m *testing.M) {

	cfg := wbf.New()

	if err := cfg.LoadEnvFiles("../../../.env"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := cfg.LoadConfigFiles("../../../config.yaml"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var appCfg config.Config
	if err := cfg.Unmarshal(&appCfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	testCfg := config.Storage{
		Host:               appCfg.Storage.Host,
		Port:               appCfg.Storage.Port,
		Username:           os.Getenv("DB_USER"),
		Password:           os.Getenv("DB_PASSWORD"),
		DBName:             appCfg.Storage.DBName,
		SSLMode:            appCfg.Storage.SSLMode,
		MaxOpenConns:       appCfg.Storage.MaxOpenConns,
		MaxIdleConns:       appCfg.Storage.MaxIdleConns,
		ConnMaxLifetime:    appCfg.Storage.ConnMaxLifetime,
		QueryRetryStrategy: appCfg.Storage.QueryRetryStrategy,
	}

	logger, _ := logger.NewLogger(config.Logger{Debug: true})

	var err error
	testDB, err = repository.ConnectDB(testCfg)
	if err != nil {
		logger.LogFatal("postgres_test — failed to connect to test DB", err, "layer", "repository.postgres_test")
	}

	if err := migrate(testDB.Master); err != nil {
		logger.LogError("postgres_test — failed to run migrations", err, "layer", "repository.postgres_test")
		os.Exit(1)
	}

	testStorage = repository.NewStorage(logger, testCfg, testDB)

	exitCode := m.Run()
	testStorage.Close()
	os.Exit(exitCode)

}

func migrate(db *sql.DB) error {
	_ = goose.SetDialect("postgres")
	if err := goose.Up(db, "../../../migrations"); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}
	return nil
}

func setupTest(t *testing.T) {

	ctx := context.Background()
	_, err := testDB.Master.ExecContext(ctx, `

	TRUNCATE TABLE users, events, bookings 
	RESTART IDENTITY CASCADE`)

	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}

}

func TestCreateUser_Errors(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	user := models.User{Login: "testdup", Password: "pass"}

	_, err := testStorage.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("first CreateUser failed: %v", err)
	}

	_, err = testStorage.CreateUser(ctx, user)
	if err == nil {
		t.Fatalf("expected error for duplicate username, got nil")
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && string(pqErr.Code) == "23505" {
		t.Logf("expected unique violation error captured: %v", err)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

}

func TestGetUserByLogin(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	user := models.User{Login: "qwe", Password: "qweqweq123"}

	userID, err := testStorage.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	got, err := testStorage.GetUserByLogin(ctx, "qwe")
	if err != nil {
		t.Fatalf("GetUserByLogin failed: %v", err)
	}

	if got.ID != userID || got.Login != user.Login || got.Password != user.Password {
		t.Fatalf("unexpected user data: %+v", got)
	}

	_, err = testStorage.GetUserByLogin(ctx, "nonexistent")
	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows for non-existent user, got %v", err)
	}

}

func TestCreateEvent_Errors(t *testing.T) {

	setupTest(t)

	ctx := context.Background()
	event := &models.Event{
		ID:          validEventUUID,
		UserID:      999999999,
		Title:       "Invalid Event",
		Description: "test",
		Date:        fixedEventDate,
		Seats:       50,
		BookingTTL:  15 * time.Minute,
	}

	err := testStorage.CreateEvent(ctx, event)
	if err == nil {
		t.Fatalf("expected error for invalid userID, got nil")
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && string(pqErr.Code) == "23503" {
		t.Logf("expected foreign key error captured: %v", err)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

}

func TestCreateEvent_GetAllEvents(t *testing.T) {

	setupTest(t)

	ctx := context.Background()

	user := models.User{Login: "eventuser", Password: "pass"}
	userID, err := testStorage.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	event := &models.Event{
		ID:          validEventUUID,
		UserID:      userID,
		Title:       "Cool TesT tiTle",
		Description: "bla bla bla",
		Date:        fixedEventDate,
		Seats:       100,
		BookingTTL:  30 * time.Minute,
	}

	if err := testStorage.CreateEvent(ctx, event); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	events, err := testStorage.GetAllEvents(ctx)
	if err != nil {
		t.Fatalf("GetAllEvents failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	e := events[0]
	if e.ID != event.ID ||
		e.Title != event.Title ||
		e.Description != event.Description ||
		!e.Date.Equal(fixedEventDate) ||
		e.Seats != event.Seats {
		t.Fatalf("event data mismatch\ngot: %+v\nwant: %+v", e, event)
	}

}

func TestGetInfo(t *testing.T) {

	setupTest(t)

	ctx := context.Background()

	user := models.User{Login: "infouser", Password: "pass"}
	_, err := testStorage.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	event := &models.Event{
		ID:          validEventUUID,
		UserID:      1,
		Title:       "Info Event",
		Description: "info bla bla bla",
		Date:        fixedEventDate,
		Seats:       200,
		BookingTTL:  10 * time.Minute,
	}

	if err := testStorage.CreateEvent(ctx, event); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	got, err := testStorage.GetInfo(ctx, validEventUUID)
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}

	if got.Title != event.Title ||
		got.Description != event.Description ||
		!got.Date.Equal(fixedEventDate) ||
		got.Seats != event.Seats {
		t.Fatalf("event info mismatch\ngot: %+v\nwant: %+v", got, event)
	}

}

func TestCancelBooking(t *testing.T) {

	setupTest(t)

	ctx := context.Background()

	user := models.User{Login: "bookuser", Password: "pass"}
	userID, err := testStorage.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	event := &models.Event{
		ID:          validEventUUID,
		UserID:      userID,
		Title:       "Booking Cool TesT tiTle",
		Description: "for cancel test",
		Date:        fixedEventDate,
		Seats:       5,
		BookingTTL:  30 * time.Minute,
	}

	if err := testStorage.CreateEvent(ctx, event); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	var bookingID int64
	err = testStorage.Transaction(ctx, func(tx *sql.Tx, txCtx context.Context) error {
		ev, err := testStorage.GetEventForBooking(tx, txCtx, validEventUUID)
		if err != nil {
			return fmt.Errorf("GetEventForBooking failed: %w", err)
		}
		booking := &models.Booking{
			UserID:    userID,
			EventID:   ev.DBID,
			Status:    models.StatusPending,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}
		bookingID, err = testStorage.CreateBooking(tx, txCtx, booking)
		if err != nil {
			return fmt.Errorf("CreateBooking failed: %w", err)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("failed to create booking: %v", err)
	}

	err = testStorage.Transaction(ctx, func(tx *sql.Tx, txCtx context.Context) error {
		_, err := testStorage.CancelBooking(tx, txCtx, bookingID)
		return err
	})
	if err != nil {
		t.Fatalf("CancelBooking failed: %v", err)
	}

	var status string
	err = testDB.Master.QueryRowContext(ctx, `
	
	SELECT status 
	FROM bookings 
	WHERE id = $1`,

		bookingID).Scan(&status)
	if err != nil {
		t.Fatalf("verify query failed: %v", err)
	}
	if status != "expired" {
		t.Fatalf("expected status 'expired', got %s", status)
	}

}

func TestBookingFullLifecycle(t *testing.T) {

	setupTest(t)
	ctx := context.Background()

	user := models.User{Login: "lifecycle", Password: "pass"}
	userID, _ := testStorage.CreateUser(ctx, user)

	event := &models.Event{
		ID:          validEventUUID,
		UserID:      userID,
		Title:       "Lifecycle Event",
		Description: "full booking test",
		Date:        fixedEventDate,
		Seats:       10,
		BookingTTL:  30 * time.Minute,
	}
	_ = testStorage.CreateEvent(ctx, event)

	var bookingID int64

	err := testStorage.Transaction(ctx, func(tx *sql.Tx, txCtx context.Context) error {
		ev, err := testStorage.GetEventForBooking(tx, txCtx, validEventUUID)
		if err != nil {
			return err
		}
		if err := testStorage.UpdateEventSeats(tx, txCtx, false, ev.DBID); err != nil {
			return err
		}
		booking := &models.Booking{
			UserID:    userID,
			EventID:   ev.DBID,
			Status:    models.StatusPending,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}
		bookingID, err = testStorage.CreateBooking(tx, txCtx, booking)
		return err
	})

	if err != nil {
		t.Fatalf("full booking transaction failed: %v", err)
	}

	err = testStorage.Transaction(ctx, func(tx *sql.Tx, txCtx context.Context) error {
		booking, err := testStorage.GetBookingForConfirm(tx, txCtx, userID, validEventUUID)
		if err != nil {
			return err
		}
		return testStorage.UpdateBookingStatus(tx, txCtx, booking.ID, "confirmed")
	})
	if err != nil {
		t.Fatalf("confirm booking failed: %v", err)
	}

	var status string
	err = testDB.Master.QueryRowContext(ctx, `
	
	SELECT status 
	FROM bookings 
	WHERE id = $1`,

		bookingID).Scan(&status)
	if err != nil || status != "confirmed" {
		t.Fatalf("expected status 'confirmed', got %s", status)
	}

}

func TestUpdateEventSeats(t *testing.T) {

	setupTest(t)
	ctx := context.Background()

	userID, _ := testStorage.CreateUser(ctx, models.User{Login: "seatsuser", Password: "pass"})
	event := &models.Event{
		ID:         validEventUUID,
		UserID:     userID,
		Title:      "Seats Test",
		Seats:      5,
		BookingTTL: 30 * time.Minute,
	}
	_ = testStorage.CreateEvent(ctx, event)

	err := testStorage.Transaction(ctx, func(tx *sql.Tx, txCtx context.Context) error {
		ev, _ := testStorage.GetEventForBooking(tx, txCtx, validEventUUID)
		if err := testStorage.UpdateEventSeats(tx, txCtx, false, ev.DBID); err != nil {
			return err
		}
		return testStorage.UpdateEventSeats(tx, txCtx, true, ev.DBID)
	})

	if err != nil {
		t.Fatalf("UpdateEventSeats failed: %v", err)
	}

}

func TestGetBookingForConfirm_Error(t *testing.T) {

	setupTest(t)
	ctx := context.Background()

	userID, _ := testStorage.CreateUser(ctx, models.User{Login: "confirmuser", Password: "pass"})
	err := testStorage.Transaction(ctx, func(tx *sql.Tx, txCtx context.Context) error {
		_, err := testStorage.GetBookingForConfirm(tx, txCtx, userID, nonExistentUUID)
		return err
	})

	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}

}
