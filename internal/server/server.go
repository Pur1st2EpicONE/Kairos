// Package server provides an abstraction for running and managing the HTTP server.
package server

import (
	"Kairos/internal/config"
	"Kairos/internal/logger"
	"Kairos/internal/server/httpserver"
	"context"
	"net/http"
)

type Server interface {
	Run()
	Shutdown()
}

func NewServer(logger logger.Logger, config config.Server, handler http.Handler, cancel context.CancelFunc) Server {
	return httpserver.NewServer(logger, config, handler, cancel)
}
