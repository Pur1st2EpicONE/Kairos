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

type App struct {
	logger  logger.Logger
	logFile *os.File
	broker  broker.Broker
	server  server.Server
	ctx     context.Context
	cancel  context.CancelFunc
	storage *repository.Storage
}

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

func processFunc(service service.Service) func(ctx context.Context, bookingID int64) error {
	return func(ctx context.Context, bookingID int64) error {
		return service.CancelBooking(ctx, bookingID)
	}
}

func (a *App) Run() {

	go a.server.Run()
	go a.broker.Consume()

	<-a.ctx.Done()

	a.Stop()

}

func (a *App) Stop() {

	a.server.Shutdown()
	a.broker.Shutdown()

	a.storage.Close()

	if a.logFile != nil && a.logFile != os.Stdout {
		_ = a.logFile.Close()
	}

}
