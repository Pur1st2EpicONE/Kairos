package impl

import (
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
