// Package v1 provides version 1 of the Kairos API handlers for notifications.
// It includes endpoints to create, query, and cancel notifications via HTTP.
package v1

import (
	"Kairos/internal/config"
	"Kairos/internal/service"
)

type Handler struct {
	config  config.Server
	service service.Service
}

func NewHandler(config config.Server, service service.Service) *Handler {
	return &Handler{config: config, service: service}
}
