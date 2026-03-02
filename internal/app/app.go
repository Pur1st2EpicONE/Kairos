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
	"sync"
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

	db, err := connectDB(logger, config.Storage)
	if err != nil {
		logger.LogFatal("app — failed to connect to database", err, "layer", "app")
	}

	app, err := wireApp(db, logger, logFile, config)
	if err != nil {
		logger.LogFatal("app — failed to connect to broker", err, "layer", "app")
	}

	return app

}

func connectDB(logger logger.Logger, config config.Storage) (*dbpg.DB, error) {

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
	broker, err := broker.NewBroker(logger, config.Broker, storage, notifier)
	service := service.NewService(logger, broker, storage)
	handler := handler.NewHandler(service)
	server := server.NewServer(logger, config.Server, handler)

	if err != nil {
		return nil, err
	}

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

func (a *App) Run() {

	var wg sync.WaitGroup

	wg.Go(func() {
		if err := a.server.Run(); err != nil {
			a.logger.LogFatal("server run failed", err, "layer", "app")
		}
	})

	// wg.Go(func() {
	// 	if err := a.broker.Consume(); err != nil {
	// 		a.logger.LogFatal("consumer run failed", err, "layer", "app")
	// 	}
	// })

	<-a.ctx.Done()

	a.Stop(&wg)

}

func (a *App) Stop(wg *sync.WaitGroup) {

	a.server.Shutdown()
	a.broker.Shutdown()

	wg.Wait()

	a.storage.Close()

	if a.logFile != nil && a.logFile != os.Stdout {
		_ = a.logFile.Close()
	}

}
