package impl

import (
	"Kairos/internal/broker"
	"Kairos/internal/config"
	"Kairos/internal/logger"
	"Kairos/internal/notifier"
	"Kairos/internal/repository"
)

type AuthService struct {
	logger  logger.Logger
	config  config.Service
	storage repository.AuthStorage
}

func NewAuthService(logger logger.Logger, config config.Service, storage repository.AuthStorage) *AuthService {
	return &AuthService{logger: logger, config: config, storage: storage}
}

type CoreService struct {
	logger   logger.Logger
	broker   broker.Broker
	storage  repository.CoreStorage
	notifier notifier.Notifier
}

func NewCoreService(logger logger.Logger, broker broker.Broker, storage repository.CoreStorage, notifier notifier.Notifier) *CoreService {
	return &CoreService{logger: logger, broker: broker, storage: storage, notifier: notifier}
}
