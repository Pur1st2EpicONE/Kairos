// Package postgres implements the repository interfaces using PostgreSQL.
// It provides concrete storage types for authentication and core operations.
package postgres

import (
	"Kairos/internal/config"
	"Kairos/internal/logger"

	"github.com/wb-go/wbf/dbpg"
)

// AuthStorage is the PostgreSQL implementation of repository.AuthStorage.
// It handles user creation and lookup.
type AuthStorage struct {
	db     *dbpg.DB       // database connection
	logger logger.Logger  // structured logger
	config config.Storage // storage configuration
}

// NewAuthStorage creates a new AuthStorage instance with the given dependencies.
func NewAuthStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) *AuthStorage {
	return &AuthStorage{db: db, logger: logger, config: config}
}

// CoreStorage is the PostgreSQL implementation of repository.CoreStorage.
// It handles all event and booking operations, including transactions.
type CoreStorage struct {
	db     *dbpg.DB       // database connection
	logger logger.Logger  // structured logger
	config config.Storage // storage configuration
}

// NewCoreStorage creates a new CoreStorage instance with the given dependencies.
func NewCoreStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) *CoreStorage {
	return &CoreStorage{db: db, logger: logger, config: config}
}
