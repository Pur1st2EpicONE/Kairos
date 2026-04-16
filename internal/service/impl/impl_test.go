package impl_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	broker_mocks "Kairos/internal/broker/mocks"
	"Kairos/internal/config"
	"Kairos/internal/errs"
	logger_mocks "Kairos/internal/logger/mocks"
	"Kairos/internal/models"
	notifier_mocks "Kairos/internal/notifier/mocks"
	repository_mocks "Kairos/internal/repository/mocks"
	"Kairos/internal/service/impl"

	"github.com/golang-jwt/jwt"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_CreateUser(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := logger_mocks.NewMockLogger(controller)
	mockStorage := repository_mocks.NewMockAuthStorage(controller)

	config := config.Service{
		TokenTTL:          15 * time.Minute,
		TokenSignedString: "aboba-token-string",
	}

	service := impl.NewAuthService(mockLogger, config, mockStorage)

	t.Run("password too long", func(t *testing.T) {
		_, err := service.CreateUser(context.Background(), models.User{
			Login:    "test",
			Password: string(make([]byte, 100)),
		})
		require.ErrorIs(t, err, errs.ErrPasswordTooLong)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockStorage.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(int64(0), &pq.Error{Code: "23505"})
		mockLogger.EXPECT().LogError(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		_, err := service.CreateUser(context.Background(), models.User{Login: "test", Password: "pass123"})
		require.ErrorIs(t, err, errs.ErrUserAlreadyExists)
	})

	t.Run("storage generic error", func(t *testing.T) {
		mockStorage.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("db error"))
		mockLogger.EXPECT().LogError(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		_, err := service.CreateUser(context.Background(), models.User{Login: "test", Password: "pass123"})
		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		mockStorage.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(int64(42), nil)
		id, err := service.CreateUser(context.Background(), models.User{Login: "test", Password: "pass123"})
		require.NoError(t, err)
		require.Equal(t, int64(42), id)
	})

}

func TestAuthService_CreateToken(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := logger_mocks.NewMockLogger(controller)
	mockStorage := repository_mocks.NewMockAuthStorage(controller)

	config := config.Service{
		TokenTTL:          15 * time.Minute,
		TokenSignedString: "aboba-token-string",
	}

	service := impl.NewAuthService(mockLogger, config, mockStorage)

	token, err := service.CreateToken(123)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Contains(t, token, ".")

}

func TestAuthService_GetUserId(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := logger_mocks.NewMockLogger(controller)
	mockStorage := repository_mocks.NewMockAuthStorage(controller)

	config := config.Service{TokenTTL: 15 * time.Minute, TokenSignedString: "key"}
	service := impl.NewAuthService(mockLogger, config, mockStorage)

	t.Run("empty login", func(t *testing.T) {
		_, err := service.GetUserId(context.Background(), models.User{Password: "pass"})
		require.ErrorIs(t, err, errs.ErrEmptyLogin)
	})

	t.Run("empty password", func(t *testing.T) {
		_, err := service.GetUserId(context.Background(), models.User{Login: "test"})
		require.ErrorIs(t, err, errs.ErrEmptyPassword)
	})

	t.Run("user not found", func(t *testing.T) {
		mockStorage.EXPECT().GetUserByLogin(gomock.Any(), "test").Return(models.User{}, sql.ErrNoRows)
		_, err := service.GetUserId(context.Background(), models.User{Login: "test", Password: "pass123"})
		require.ErrorIs(t, err, errs.ErrInvalidCredentials)
	})

	t.Run("storage error (not NoRows)", func(t *testing.T) {
		mockStorage.EXPECT().GetUserByLogin(gomock.Any(), "test").Return(models.User{}, errors.New("db connection lost"))
		mockLogger.EXPECT().LogError("service — failed to get userID by login", gomock.Any(), gomock.Any()).AnyTimes()
		_, err := service.GetUserId(context.Background(), models.User{Login: "test", Password: "pass123"})
		require.Error(t, err)
		require.NotErrorIs(t, err, errs.ErrInvalidCredentials)
	})

	t.Run("wrong password", func(t *testing.T) {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("realpass"), bcrypt.DefaultCost)
		mockStorage.EXPECT().GetUserByLogin(gomock.Any(), "test").Return(models.User{
			ID:       1,
			Login:    "test",
			Password: string(hashed),
		}, nil)
		_, err := service.GetUserId(context.Background(), models.User{Login: "test", Password: "wrong"})
		require.ErrorIs(t, err, errs.ErrInvalidCredentials)
	})

	t.Run("success", func(t *testing.T) {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
		mockStorage.EXPECT().GetUserByLogin(gomock.Any(), "test").Return(models.User{
			ID:       42,
			Login:    "test",
			Password: string(hashed),
		}, nil)
		id, err := service.GetUserId(context.Background(), models.User{Login: "test", Password: "pass123"})
		require.NoError(t, err)
		require.Equal(t, int64(42), id)
	})

}

func TestAuthService_ParseToken(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := logger_mocks.NewMockLogger(controller)
	mockStorage := repository_mocks.NewMockAuthStorage(controller)

	config := config.Service{
		TokenTTL:          15 * time.Minute,
		TokenSignedString: "aboba-token-string",
	}
	service := impl.NewAuthService(mockLogger, config, mockStorage)

	t.Run("invalid token format", func(t *testing.T) {
		_, err := service.ParseToken("not-a-jwt")
		require.ErrorIs(t, err, errs.ErrInvalidToken)
	})

	t.Run("invalid claims or token not valid", func(t *testing.T) {
		invalidToken := "eyJhbGciOiJIUzI1NiPOk6yJV_adQssw5c"
		_, err := service.ParseToken(invalidToken)
		require.ErrorIs(t, err, errs.ErrInvalidToken)
	})

	t.Run("subject is not number", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Subject:   "not-a-number",
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		})
		signedToken, _ := token.SignedString([]byte(config.TokenSignedString))
		_, err := service.ParseToken(signedToken)
		require.ErrorIs(t, err, errs.ErrInvalidUserID)
	})

	t.Run("success", func(t *testing.T) {
		token, _ := service.CreateToken(12345)
		id, err := service.ParseToken(token)
		require.NoError(t, err)
		require.Equal(t, int64(12345), id)
	})

}

func TestCoreService_CreateEvent(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := logger_mocks.NewMockLogger(controller)
	mockStorage := repository_mocks.NewMockCoreStorage(controller)
	mockBroker := broker_mocks.NewMockBroker(controller)
	mockNotifier := notifier_mocks.NewMockNotifier(controller)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().LogError(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockNotifier.EXPECT().Notify(gomock.Any()).AnyTimes().Return(nil)

	service := impl.NewCoreService(mockLogger, mockBroker, mockStorage, mockNotifier)

	validEvent := func() *models.Event {
		return &models.Event{
			UserID:      1,
			Title:       "Test Event",
			Description: "Desc",
			Date:        time.Now().UTC().Add(48 * time.Hour),
			Seats:       50,
			BookingTTL:  2 * time.Hour,
		}
	}

	t.Run("validation errors", func(t *testing.T) {
		cases := []struct {
			name string
			mod  func(*models.Event)
			err  error
		}{
			{"missing title", func(e *models.Event) { e.Title = "" }, errs.ErrMissingTitle},
			{"title too short", func(e *models.Event) { e.Title = "ab" }, errs.ErrTitleTooShort},
			{"title too long", func(e *models.Event) { e.Title = string(make([]byte, 201)) }, errs.ErrTitleTooLong},
			{"description too long", func(e *models.Event) { e.Description = string(make([]byte, 2001)) }, errs.ErrDescriptionTooLong},
			{"missing date (zero)", func(e *models.Event) { e.Date = time.Time{} }, errs.ErrMissingDate},
			{"date in past", func(e *models.Event) { e.Date = time.Now().UTC().Add(-time.Hour) }, errs.ErrDateInPast},
			{"date too soon", func(e *models.Event) { e.Date = time.Now().UTC().Add(12 * time.Hour) }, errs.ErrDateTooSoon},
			{"date too far", func(e *models.Event) { e.Date = time.Now().UTC().AddDate(2, 0, 0) }, errs.ErrDateTooFar},
			{"invalid seats (zero)", func(e *models.Event) { e.Seats = 0 }, errs.ErrInvalidSeatCount},
			{"too many seats", func(e *models.Event) { e.Seats = 10001 }, errs.ErrTooManySeats},
			{"booking TTL too short", func(e *models.Event) { e.BookingTTL = 30 * time.Second }, errs.ErrBookingTTLTooShort},
			{"booking TTL too long", func(e *models.Event) { e.BookingTTL = 25 * time.Hour }, errs.ErrBookingTTLTooLong},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				e := validEvent()
				tc.mod(e)
				_, err := service.CreateEvent(context.Background(), e)
				require.ErrorIs(t, err, tc.err)
			})
		}
	})

	t.Run("invalid userID (foreign key)", func(t *testing.T) {
		mockStorage.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(&pq.Error{Code: "23503"})
		_, err := service.CreateEvent(context.Background(), validEvent())
		require.ErrorIs(t, err, errs.ErrInvalidUserID)
	})

	t.Run("storage generic error", func(t *testing.T) {
		mockStorage.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
		_, err := service.CreateEvent(context.Background(), validEvent())
		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		mockStorage.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(nil)
		id, err := service.CreateEvent(context.Background(), validEvent())
		require.NoError(t, err)
		require.NotEmpty(t, id)
	})

}

func TestCoreService_CreateBooking(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := logger_mocks.NewMockLogger(controller)
	mockStorage := repository_mocks.NewMockCoreStorage(controller)
	mockBroker := broker_mocks.NewMockBroker(controller)
	mockNotifier := notifier_mocks.NewMockNotifier(controller)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().LogError(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockNotifier.EXPECT().Notify(gomock.Any()).AnyTimes().Return(nil)

	service := impl.NewCoreService(mockLogger, mockBroker, mockStorage, mockNotifier)

	eventID := "123e4567-e89b-12d3-a456-426614174000"
	userID := int64(1)

	t.Run("event not found", func(t *testing.T) {
		mockStorage.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(*sql.Tx, context.Context) error) error { return fn(nil, ctx) })
		mockStorage.EXPECT().GetEventForBooking(gomock.Any(), gomock.Any(), eventID).Return(nil, sql.ErrNoRows)
		_, err := service.CreateBooking(context.Background(), userID, eventID)
		require.ErrorIs(t, err, errs.ErrEventNotFound)
	})

	t.Run("event full", func(t *testing.T) {
		mockStorage.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(*sql.Tx, context.Context) error) error { return fn(nil, ctx) })
		mockStorage.EXPECT().GetEventForBooking(gomock.Any(), gomock.Any(), eventID).Return(&models.Event{Seats: 0}, nil)
		_, err := service.CreateBooking(context.Background(), userID, eventID)
		require.ErrorIs(t, err, errs.ErrEventFull)
	})

	t.Run("booking already exists", func(t *testing.T) {
		mockStorage.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(*sql.Tx, context.Context) error) error { return fn(nil, ctx) })
		mockStorage.EXPECT().GetEventForBooking(gomock.Any(), gomock.Any(), eventID).Return(&models.Event{Seats: 10, DBID: 5, BookingTTL: time.Hour}, nil)
		mockStorage.EXPECT().CreateBooking(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), &pq.Error{Code: "23505"})
		_, err := service.CreateBooking(context.Background(), userID, eventID)
		require.ErrorIs(t, err, errs.ErrBookingAlreadyExists)
	})

	t.Run("broker error", func(t *testing.T) {
		mockStorage.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(*sql.Tx, context.Context) error) error { return fn(nil, ctx) })
		mockStorage.EXPECT().GetEventForBooking(gomock.Any(), gomock.Any(), eventID).Return(&models.Event{Seats: 10, DBID: 5, BookingTTL: time.Hour}, nil)
		mockStorage.EXPECT().CreateBooking(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(777), nil)
		mockBroker.EXPECT().Produce(gomock.Any()).Return(errors.New("kafka down"))
		_, err := service.CreateBooking(context.Background(), userID, eventID)
		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		mockStorage.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(*sql.Tx, context.Context) error) error { return fn(nil, ctx) })
		mockStorage.EXPECT().GetEventForBooking(gomock.Any(), gomock.Any(), eventID).Return(&models.Event{Seats: 10, DBID: 5, BookingTTL: time.Hour}, nil)
		mockStorage.EXPECT().CreateBooking(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(777), nil)
		mockBroker.EXPECT().Produce(gomock.Any()).Return(nil)
		mockStorage.EXPECT().UpdateEventSeats(gomock.Any(), gomock.Any(), false, int64(5)).Return(nil)
		id, err := service.CreateBooking(context.Background(), userID, eventID)
		require.NoError(t, err)
		require.Equal(t, int64(777), id)
	})

}

func TestCoreService_ConfirmBooking(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := logger_mocks.NewMockLogger(controller)
	mockStorage := repository_mocks.NewMockCoreStorage(controller)
	mockBroker := broker_mocks.NewMockBroker(controller)
	mockNotifier := notifier_mocks.NewMockNotifier(controller)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().LogError(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockNotifier.EXPECT().Notify(gomock.Any()).AnyTimes().Return(nil)

	service := impl.NewCoreService(mockLogger, mockBroker, mockStorage, mockNotifier)

	t.Run("booking not found", func(t *testing.T) {
		mockStorage.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(*sql.Tx, context.Context) error) error { return fn(nil, ctx) })
		mockStorage.EXPECT().GetBookingForConfirm(gomock.Any(), gomock.Any(), int64(1), "event-uuid").Return(nil, sql.ErrNoRows)
		err := service.ConfirmBooking(context.Background(), 1, "event-uuid")
		require.ErrorIs(t, err, errs.ErrBookingNotFound)
	})

	t.Run("already confirmed", func(t *testing.T) {
		mockStorage.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(*sql.Tx, context.Context) error) error { return fn(nil, ctx) })
		mockStorage.EXPECT().GetBookingForConfirm(gomock.Any(), gomock.Any(), int64(1), "event-uuid").Return(&models.Booking{Status: models.StatusConfirmed}, nil)
		err := service.ConfirmBooking(context.Background(), 1, "event-uuid")
		require.ErrorIs(t, err, errs.ErrAlreadyConfirmed)
	})

	t.Run("expired", func(t *testing.T) {
		mockStorage.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(*sql.Tx, context.Context) error) error { return fn(nil, ctx) })
		mockStorage.EXPECT().GetBookingForConfirm(gomock.Any(), gomock.Any(), int64(1), "event-uuid").Return(&models.Booking{Status: models.StatusExpired}, nil)
		err := service.ConfirmBooking(context.Background(), 1, "event-uuid")
		require.ErrorIs(t, err, errs.ErrBookingExpired)
	})

	t.Run("success", func(t *testing.T) {
		mockStorage.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(*sql.Tx, context.Context) error) error { return fn(nil, ctx) })
		mockStorage.EXPECT().GetBookingForConfirm(gomock.Any(), gomock.Any(), int64(1), "event-uuid").Return(&models.Booking{ID: 999, Status: models.StatusPending}, nil)
		mockStorage.EXPECT().UpdateBookingStatus(gomock.Any(), gomock.Any(), int64(999), models.StatusConfirmed).Return(nil)
		err := service.ConfirmBooking(context.Background(), 1, "event-uuid")
		require.NoError(t, err)
	})

}

func TestCoreService_CancelBooking(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := logger_mocks.NewMockLogger(controller)
	mockStorage := repository_mocks.NewMockCoreStorage(controller)
	mockBroker := broker_mocks.NewMockBroker(controller)
	mockNotifier := notifier_mocks.NewMockNotifier(controller)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().LogError(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockNotifier.EXPECT().Notify(gomock.Any()).AnyTimes().Return(nil)

	service := impl.NewCoreService(mockLogger, mockBroker, mockStorage, mockNotifier)

	t.Run("no rows → silent success", func(t *testing.T) {
		mockStorage.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(*sql.Tx, context.Context) error) error { return fn(nil, ctx) })
		mockStorage.EXPECT().CancelBooking(gomock.Any(), gomock.Any(), int64(999)).Return(int64(0), sql.ErrNoRows)
		err := service.CancelBooking(context.Background(), 999)
		require.NoError(t, err)
	})

	t.Run("cancel error", func(t *testing.T) {
		mockStorage.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(*sql.Tx, context.Context) error) error { return fn(nil, ctx) })
		mockStorage.EXPECT().CancelBooking(gomock.Any(), gomock.Any(), int64(999)).Return(int64(0), errors.New("db fail"))
		err := service.CancelBooking(context.Background(), 999)
		require.Error(t, err)
	})

	t.Run("update seats error", func(t *testing.T) {
		mockStorage.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(*sql.Tx, context.Context) error) error { return fn(nil, ctx) })
		mockStorage.EXPECT().CancelBooking(gomock.Any(), gomock.Any(), int64(999)).Return(int64(5), nil)
		mockStorage.EXPECT().UpdateEventSeats(gomock.Any(), gomock.Any(), true, int64(5)).Return(errors.New("seat update fail"))
		err := service.CancelBooking(context.Background(), 999)
		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		mockStorage.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(*sql.Tx, context.Context) error) error { return fn(nil, ctx) })
		mockStorage.EXPECT().CancelBooking(gomock.Any(), gomock.Any(), int64(999)).Return(int64(5), nil)
		mockStorage.EXPECT().UpdateEventSeats(gomock.Any(), gomock.Any(), true, int64(5)).Return(nil)
		err := service.CancelBooking(context.Background(), 999)
		require.NoError(t, err)
	})

}

func TestCoreService_GetInfo(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := logger_mocks.NewMockLogger(controller)
	mockStorage := repository_mocks.NewMockCoreStorage(controller)
	mockBroker := broker_mocks.NewMockBroker(controller)
	mockNotifier := notifier_mocks.NewMockNotifier(controller)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().LogError(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockNotifier.EXPECT().Notify(gomock.Any()).AnyTimes().Return(nil)

	service := impl.NewCoreService(mockLogger, mockBroker, mockStorage, mockNotifier)

	t.Run("not found", func(t *testing.T) {
		mockStorage.EXPECT().GetInfo(gomock.Any(), "event-uuid").Return(nil, sql.ErrNoRows)
		_, err := service.GetInfo(context.Background(), "event-uuid")
		require.ErrorIs(t, err, errs.ErrEventNotFound)
	})

	t.Run("storage error", func(t *testing.T) {
		mockStorage.EXPECT().GetInfo(gomock.Any(), "event-uuid").Return(nil, errors.New("db error"))
		_, err := service.GetInfo(context.Background(), "event-uuid")
		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		expected := &models.Event{ID: "event-uuid", Title: "Test"}
		mockStorage.EXPECT().GetInfo(gomock.Any(), "event-uuid").Return(expected, nil)
		event, err := service.GetInfo(context.Background(), "event-uuid")
		require.NoError(t, err)
		require.Equal(t, expected, event)
	})

}

func TestCoreService_GetAllEvents(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := logger_mocks.NewMockLogger(controller)
	mockStorage := repository_mocks.NewMockCoreStorage(controller)
	mockBroker := broker_mocks.NewMockBroker(controller)
	mockNotifier := notifier_mocks.NewMockNotifier(controller)

	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().LogError(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockNotifier.EXPECT().Notify(gomock.Any()).AnyTimes().Return(nil)

	service := impl.NewCoreService(mockLogger, mockBroker, mockStorage, mockNotifier)

	t.Run("success", func(t *testing.T) {
		data := []models.Event{{ID: "1"}, {ID: "2"}}
		mockStorage.EXPECT().GetAllEvents(gomock.Any()).Return(data, nil)
		events := service.GetAllEvents(context.Background())
		require.Equal(t, data, events)
	})

	t.Run("storage error", func(t *testing.T) {
		mockStorage.EXPECT().GetAllEvents(gomock.Any()).Return(nil, errors.New("db error"))
		service.GetAllEvents(context.Background())
	})

}
