// Package config provides configuration loading and parsing for the Kairos service.
// It supports YAML configuration files, environment variable overrides, and
// provides strongly-typed structs for all application settings.
package config

import (
	"fmt"
	"os"
	"time"

	wbf "github.com/wb-go/wbf/config"
)

// Config is the root configuration structure aggregating all subsections.
type Config struct {
	Logger   Logger   `mapstructure:"logger"`   // Logger configuration (debug mode, log directory)
	Notifier Notifier `mapstructure:"notifier"` // Notifier settings for Telegram/email alerts
	Server   Server   `mapstructure:"server"`   // HTTP server settings (port, timeouts)
	Service  Service  `mapstructure:"service"`  // Business logic layer settings (JWT TTL, signing key)
	Storage  Storage  `mapstructure:"database"` // Database connection and migration settings
	Broker   Broker   `mapstructure:"broker"`   // Message broker (RabbitMQ) settings
}

// Notifier holds configuration for sending notifications via Telegram.
type Notifier struct {
	TelegramToken    string // Bot token for Telegram API
	TelegramReceiver string // Chat ID or username that receives Telegram messages
}

// Logger defines logging behaviour: debug mode and log file directory.
type Logger struct {
	Debug  bool   `mapstructure:"debug_mode"`    // Enable debug-level logging
	LogDir string `mapstructure:"log_directory"` // Directory where log files are stored
}

// Server configures the HTTP server: port, timeouts, and shutdown grace period.
type Server struct {
	Port            string        `mapstructure:"port"`             // TCP port to listen on
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`     // Maximum duration for reading the entire request
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`    // Maximum duration before timing out writes
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"` // Maximum size of request headers
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"` // Grace period for server shutdown
}

// Service holds application‑specific business logic settings.
type Service struct {
	TokenTTL          time.Duration `mapstructure:"token_ttl"` // JWT token validity duration
	TokenSignedString string        // Secret key used to sign JWT tokens (loaded from env)
}

// Broker configures the RabbitMQ message broker: connection, queues,
// producers, consumers, and health check intervals.
type Broker struct {
	URL                 string        `mapstructure:"url"`                  // RabbitMQ connection URL (amqp://...)
	QueueName           string        `mapstructure:"queue_name"`           // Main queue name for booking cancellations
	ConnectionName      string        `mapstructure:"connection_name"`      // Human-readable connection name
	ConnectTimeout      time.Duration `mapstructure:"connect_timeout"`      // Timeout for establishing connection
	Reconnect           Producer      `mapstructure:"reconnect"`            // Retry strategy for reconnecting
	Producer            Producer      `mapstructure:"producer"`             // Producer‑specific retry and TTL settings
	Consumer            Consumer      `mapstructure:"consumer"`             // Consumer concurrency and ack settings
	CleanupInterval     time.Duration `mapstructure:"cleanup_interval"`     // Interval for cleaning up expired queues
	HealthcheckInterval time.Duration `mapstructure:"healthcheck_interval"` // Interval for broker health checks
}

// Producer defines retry behaviour and message queue TTL for the producer.
type Producer struct {
	Attempts        int           `mapstructure:"attempts"`          // Number of publish attempts
	Delay           time.Duration `mapstructure:"delay"`             // Initial delay between retries
	Backoff         float64       `mapstructure:"backoff"`           // Exponential backoff factor
	MessageQueueTTL time.Duration `mapstructure:"message_queue_ttl"` // Time-to-live for temporary queues
}

// Consumer configures message consumption: retries, concurrency, and acknowledgement.
type Consumer struct {
	Attempts      int           `mapstructure:"attempts"`       // Number of handler retries
	Delay         time.Duration `mapstructure:"delay"`          // Initial delay between retries
	Backoff       float64       `mapstructure:"backoff"`        // Exponential backoff factor
	ConsumerTag   string        `mapstructure:"consumer_tag"`   // Unique identifier for this consumer
	AutoAck       bool          `mapstructure:"auto_ack"`       // Automatically acknowledge messages
	Workers       int           `mapstructure:"workers"`        // Number of concurrent message processors
	PrefetchCount int           `mapstructure:"prefetch_count"` // Maximum number of unacknowledged messages
}

// Storage holds PostgreSQL connection parameters, migration settings,
// connection pool limits, and query retry strategy.
type Storage struct {
	Dialect            string             `mapstructure:"goose_dialect"`              // Goose migration dialect
	MigrationsDir      string             `mapstructure:"goose_migrations_directory"` // Directory containing migration files
	Host               string             `mapstructure:"host"`                       // Database host
	Port               string             `mapstructure:"port"`                       // Database port
	Username           string             `mapstructure:"username"`                   // Database user
	Password           string             `mapstructure:"password"`                   // Database password
	DBName             string             `mapstructure:"dbname"`                     // Database name
	SSLMode            string             `mapstructure:"sslmode"`                    // SSL mode (disable, require, verify-ca, verify-full)
	MaxOpenConns       int                `mapstructure:"max_open_conns"`             // Maximum number of open connections
	MaxIdleConns       int                `mapstructure:"max_idle_conns"`             // Maximum number of idle connections
	ConnMaxLifetime    time.Duration      `mapstructure:"conn_max_lifetime"`          // Maximum lifetime of a connection
	QueryRetryStrategy QueryRetryStrategy `mapstructure:"query_retry_strategy"`       // Retry strategy for failed queries
}

// QueryRetryStrategy defines how to retry database queries on transient errors.
type QueryRetryStrategy struct {
	Attempts int           `mapstructure:"attempts"` // Number of retry attempts
	Delay    time.Duration `mapstructure:"delay"`    // Initial delay between attempts
	Backoff  float64       `mapstructure:"backoff"`  // Exponential backoff multiplier
}

// Load reads configuration from YAML files and environment variables,
// then returns a populated Config struct. It uses the wbf/config helper.
// If loading or unmarshalling fails, an error is returned.
func Load() (Config, error) {

	cfg := wbf.New()

	if err := cfg.LoadConfigFiles("./config.yaml"); err != nil {
		return Config{}, err
	}

	if err := cfg.LoadEnvFiles(".env"); err != nil && !cfg.GetBool("docker") {
		return Config{}, err
	}

	var conf Config

	if err := cfg.Unmarshal(&conf); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	loadEnvs(&conf)

	return conf, nil

}

// loadEnvs overrides specific configuration fields with values from environment
// variables. It is called after unmarshalling YAML to apply secrets and
// environment‑specific overrides.
func loadEnvs(conf *Config) {

	conf.Storage.Username = os.Getenv("DB_USER")
	conf.Storage.Password = os.Getenv("DB_PASSWORD")

	conf.Notifier.TelegramToken = os.Getenv("TG_BOT_TOKEN")
	conf.Notifier.TelegramReceiver = os.Getenv("TG_CHAT_ID")

	conf.Service.TokenSignedString = os.Getenv("JWT_TOKEN_SIGNED_STRING")

}
