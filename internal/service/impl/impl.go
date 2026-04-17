// Package impl provides concrete implementations of the service interfaces.
package impl

import (
	"Kairos/internal/broker"
	"Kairos/internal/config"
	"Kairos/internal/logger"
	"Kairos/internal/notifier"
	"Kairos/internal/repository"
)

// AuthService implements the AuthService interface using the repository's AuthStorage.
type AuthService struct {
	logger  logger.Logger          // structured logger
	config  config.Service         // service configuration (JWT TTL, signing key)
	storage repository.AuthStorage // data access for users
}

// NewAuthService creates a new AuthService with the given dependencies.
func NewAuthService(logger logger.Logger, config config.Service, storage repository.AuthStorage) *AuthService {
	return &AuthService{logger: logger, config: config, storage: storage}
}

// CoreService implements the CoreService interface using the repository's CoreStorage,
// the message broker, and the notifier.
type CoreService struct {
	logger   logger.Logger          // structured logger
	broker   broker.Broker          // message broker for delayed cancellations
	storage  repository.CoreStorage // data access for events and bookings
	notifier notifier.Notifier      // notification sender (e.g., Telegram)
}

// NewCoreService creates a new CoreService with the given dependencies.
func NewCoreService(logger logger.Logger, broker broker.Broker, storage repository.CoreStorage, notifier notifier.Notifier) *CoreService {
	return &CoreService{logger: logger, broker: broker, storage: storage, notifier: notifier}
}
