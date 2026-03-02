// Package postgres provides a PostgreSQL-backed implementation of the repository layer.
package postgres

import (
	"Kairos/internal/config"
	"Kairos/internal/logger"

	"github.com/wb-go/wbf/dbpg"
)

// CoreStorage implements the CoreStorage interface using PostgreSQL.
type CoreStorage struct {
	db     *dbpg.DB       // Database connection pool
	logger logger.Logger  // Application logger
	config config.Storage // CoreStorage-related configuration
}

// NewStorage creates a new PostgreSQL storage instance.
func NewCoreStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) *CoreStorage {
	return &CoreStorage{db: db, logger: logger, config: config}
}
