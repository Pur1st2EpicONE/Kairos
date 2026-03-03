// Package v1 provides version 1 of the Kairos API handlers for notifications.
// It includes endpoints to create, query, and cancel notifications via HTTP.
package v1

import (
	"Kairos/internal/service"
)

// Handler is the v1 API handler for notifications.
// It wraps the service layer and provides HTTP endpoints for CRUD operations.
type Handler struct {
	service service.Service
}

// NewHandler creates a new v1 Handler with the provided service.
func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}
