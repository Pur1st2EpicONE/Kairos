// Package v1 implements the JSON API handlers for version 1 of the Kairos API.
// It provides endpoints for authentication, event management, and booking operations.
package v1

import (
	"Kairos/internal/config"
	"Kairos/internal/service"
)

// Handler aggregates the API dependencies: HTTP server configuration and the service layer.
// It is used to define all v1 route handlers as methods.
type Handler struct {
	config  config.Server   // HTTP server configuration (timeouts, etc.)
	service service.Service // Composite service containing auth, core, and booking logic
}

// NewHandler creates a new v1 API handler with the given configuration and service.
func NewHandler(config config.Server, service service.Service) *Handler {
	return &Handler{config: config, service: service}
}
