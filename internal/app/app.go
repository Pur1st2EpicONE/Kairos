// Package app wires application components, provides lifecycle management
// and exposes the entry point for booting and running the service.
package app

import (
	"Kairos/internal/broker"
	"Kairos/internal/config"
	"Kairos/internal/handler"
	"Kairos/internal/logger"
	"Kairos/internal/notifier"
	"Kairos/internal/repository"
	"Kairos/internal/server"
	"Kairos/internal/service"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"context"

	"github.com/pressly/goose/v3"
	"github.com/wb-go/wbf/dbpg"
)

// App represents the application's composition root.
// It holds long-lived resources (logger, DB, broker, server) and
// the context/cancel function used for graceful shutdown.
type App struct {
	logger  logger.Logger       // structured logger used across layers
	logFile *os.File            // file handle where logs are written (nil if stdout)
	broker  broker.Broker       // message broker for asynchronous tasks
	server  server.Server       // HTTP server instance
	ctx     context.Context     // root context for shutdown coordination
	cancel  context.CancelFunc  // cancels ctx when a shutdown signal is received
	storage *repository.Storage // data storage abstraction backed by the database
}

// Boot loads configuration, initializes logger, connects to database,
// applies migrations, wires all components, and returns a fully constructed *App.
// If any critical step fails, Boot logs a fatal error and exits.
func Boot() *App {

	config, err := config.Load()
	if err != nil {
		log.Fatalf("app — failed to load configs: %v", err)
	}

	logger, logFile := logger.NewLogger(config.Logger)

	db, err := bootstrapDB(logger, config.Storage)
	if err != nil {
		logger.LogFatal("app — failed to connect to database", err, "layer", "app")
	}

	app, err := wireApp(db, logger, logFile, config)
	if err != nil {
		logger.LogFatal("app — failed to connect to broker", err, "layer", "app")
	}

	return app

}

// bootstrapDB establishes a database connection using repository.ConnectDB,
// runs pending Goose migrations, and returns the DB handle.
// It logs a successful connection and migration application.
func bootstrapDB(logger logger.Logger, config config.Storage) (*dbpg.DB, error) {

	db, err := repository.ConnectDB(config)
	if err != nil {
		return nil, err
	}

	logger.LogInfo("app — connected to database", "layer", "app")

	if err := goose.SetDialect(config.Dialect); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db.Master, config.MigrationsDir); err != nil {
		return nil, fmt.Errorf("failed to apply goose migrations: %w", err)
	}

	logger.Debug("app — migrations applied", "layer", "app")

	return db, nil

}

// wireApp constructs all application components (storage, notifier, broker,
// service, handler, server), creates a cancellable context, and returns
// the assembled *App. It also sets a cancellation callback on the broker
// to delegate booking cancellations to the service layer.
func wireApp(db *dbpg.DB, logger logger.Logger, logFile *os.File, config config.Config) (*App, error) {

	ctx, cancel := newContext(logger)

	storage := repository.NewStorage(logger, config.Storage, db)
	notifier := notifier.NewNotifier(config.Notifier)

	broker, err := broker.NewBroker(logger, config.Broker, nil)
	if err != nil {
		return nil, err
	}

	service := service.NewService(logger, config.Service, broker, storage, notifier)

	broker.SetCancelFunc(func(ctx context.Context, bookingID int64) error {
		return service.CancelBooking(ctx, bookingID)
	})

	server := server.NewServer(logger, config.Server, handler.NewHandler(config.Server, service), cancel)

	return &App{
		logger:  logger,
		logFile: logFile,
		broker:  broker,
		server:  server,
		ctx:     ctx,
		cancel:  cancel,
		storage: storage,
	}, nil

}

// newContext creates a context that is cancelled when the process receives
// SIGINT or SIGTERM. It logs the received signal and triggers graceful shutdown
// by calling the returned cancel function.
func newContext(logger logger.Logger) (context.Context, context.CancelFunc) {

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-sigCh
		sigString := sig.String()
		if sig == syscall.SIGTERM {
			sigString = "terminate" // sig.String() returns the SIGTERM string in past tense for some reason
		}
		logger.LogInfo("app — received signal "+sigString+", initiating graceful shutdown", "layer", "app")
		cancel()
	}()

	return ctx, cancel

}

// Run starts the HTTP server and the message broker consumer in background
// goroutines, then blocks until the application's context is cancelled.
// After cancellation it calls Stop to perform orderly shutdown.
func (a *App) Run() {

	go a.server.Run()
	go a.broker.Consume()

	<-a.ctx.Done()

	a.Stop()

}

// Stop performs an orderly shutdown of application components:
// it shuts down the HTTP server, stops the broker consumer,
// closes the database storage, and closes the log file if it is not os.Stdout.
func (a *App) Stop() {

	a.server.Shutdown()
	a.broker.Shutdown()

	a.storage.Close()

	if a.logFile != nil && a.logFile != os.Stdout {
		_ = a.logFile.Close()
	}

}
