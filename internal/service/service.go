package service

import (
	"Kairos/internal/broker"
	"Kairos/internal/logger"
	"Kairos/internal/models"
	"Kairos/internal/notifier"
	"Kairos/internal/repository"
	"Kairos/internal/service/impl"
	"context"
)

type AuthService interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	CreateToken(userID int64) (string, error)
	GetUserId(ctx context.Context, user models.User) (int64, error)
	ParseToken(tokenString string) (int64, error)
}

type CoreService interface {
	CreateEvent(ctx context.Context, event *models.Event) (string, error)
	CreateBooking(ctx context.Context, userID int64, eventID string) (int64, error)
	ConfirmBooking(ctx context.Context, userID int64, eventID string) error
	CancelBooking(ctx context.Context, bookingID int64) error
	GetInfo(ctx context.Context, eventID string) (*models.Event, error)
}

type Service struct {
	AuthService
	CoreService
}

func NewService(logger logger.Logger, broker broker.Broker, storage *repository.Storage, notifier notifier.Notifier) *Service {
	return &Service{
		AuthService: impl.NewAuthService(logger, storage.AuthStorage),
		CoreService: impl.NewCoreService(logger, broker, storage.CoreStorage, notifier),
	}
}
