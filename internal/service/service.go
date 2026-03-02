// Package service contains the business logic for notifications.
// It coordinates operations between the broker, cache, and storage layers.
package service

import (
	"Kairos/internal/broker"
	"Kairos/internal/logger"
	"Kairos/internal/models"
	"Kairos/internal/repository"
	"Kairos/internal/service/impl"
	"context"
)

type AuthService interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	CreateToken(userID int) (string, error)
	GetUserId(ctx context.Context, user models.User) (int, error)
	ParseToken(tokenString string) (int64, error)
}

type CoreService interface {
	CreateEvent(ctx context.Context, event *models.Event) (string, error)
	GetAllStatuses(ctx context.Context) []models.Notification             // GetAllStatuses retrieves all notifications with their current status. Used for frontend display; not optimized.
	GetStatus(ctx context.Context, notificationID string) (string, error) // GetStatus returns the current status of a specific notification by ID.
	CancelNotification(ctx context.Context, notificationID string) error  // CancelNotification attempts to cancel a notification by ID.
}

type Service struct {
	AuthService
	CoreService
}

func NewService(logger logger.Logger, broker broker.Broker, storage *repository.Storage) *Service {
	return &Service{
		AuthService: impl.NewAuthService(logger, storage.AuthStorage),
		CoreService: impl.NewCoreService(logger, broker, storage.CoreStorage),
	}
}
