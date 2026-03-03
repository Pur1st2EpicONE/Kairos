package impl

import (
	"Kairos/internal/broker"
	"Kairos/internal/logger"
	"Kairos/internal/repository"
)

type AuthService struct {
	logger  logger.Logger
	storage repository.AuthStorage
}

func NewAuthService(logger logger.Logger, storage repository.AuthStorage) *AuthService {
	return &AuthService{logger: logger, storage: storage}
}

type CoreService struct {
	logger  logger.Logger
	broker  broker.Broker
	storage repository.CoreStorage
}

func NewCoreService(logger logger.Logger, broker broker.Broker, storage repository.CoreStorage) *CoreService {
	return &CoreService{logger: logger, broker: broker, storage: storage}
}
