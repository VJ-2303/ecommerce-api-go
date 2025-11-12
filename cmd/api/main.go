package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/VJ-2303/CityStars/internal/data"
)

const version = "1.0.0"

// config holds configuration settings for the application,
// including server port, environment, database connection string, and JWT secret.
type config struct {
	port      int    // Port number for HTTP server
	env       string // Application environment ("development", "staging", "production")
	dsn       string // PostgreSQL database connection string
	jwtSecret string // Secret key for signing JWT tokens
}

// application aggregates the application's dependencies and configuration.
type application struct {
	config config       // Application configuration
	logger *slog.Logger // Structured logger instance
	models data.Models  // Data models for database access
}

func main() {
	var cfg config // Create a config struct to hold flag values

	// Get port from environment variable (Railway sets this)
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000" // Default port if not set
	}

	// Get database DSN from environment (Railway uses DATABASE_URL)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN") // Fallback
	}

	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "mysecretkey" // Default for development only
	}

	// Get environment from env variable
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "production"
	}

	// Parse command-line flags into the config struct.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", environment, "Environment (development|staging|production)")
	flag.StringVar(&cfg.jwtSecret, "jwt-secret", jwtSecret, "JWT secret string")
	flag.StringVar(&cfg.dsn, "db-dsn", dsn, "Postgres DB connection string")
	flag.Parse()

	// Initialize a new logger that writes structured logs to standard output.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Log configuration for debugging
	logger.Info("configuration loaded",
		"port", cfg.port,
		"env", cfg.env,
		"has_dsn", cfg.dsn != "",
		"has_jwt_secret", cfg.jwtSecret != "")

	// Attempt to open a database connection using the provided DSN.
	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	// Create the application struct, injecting configuration, logger, and models.
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	// Configure the HTTP server with custom timeouts and error logging.
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.port),                  // Bind to all interfaces for Railway
		Handler:      app.routes(),                                         // HTTP handler with application routes
		IdleTimeout:  time.Minute,                                          // Maximum idle connection duration
		ReadTimeout:  15 * time.Second,                                     // Maximum duration for reading the request
		WriteTimeout: 15 * time.Second,                                     // Maximum duration for writing the response
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError), // Use structured logger for server errors
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	// Start the HTTP server and block until it stops.
	err = srv.ListenAndServe()

	// Log any server error (e.g., shutdown or unexpected failure) and exit.
	logger.Error("server stopped", "error", err)
	os.Exit(1)
}
