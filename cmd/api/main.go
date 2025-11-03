package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/VJ-2303/ecommerce-api-go/internal/data"
)

const version = "1.0.0"

type config struct {
	port      int    // Port number for http connection
	env       string // Environment(devolopment|staging|production)
	dsn       string // Postgres Database connection string
	jwtSecret string // JWT Secret string, For generating signed JWT Strings
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {
	var cfg config // Initialize an cfg struct

	// Scan the flags into the cfg struct by passing the references
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "devolopment", "Environment(devolopment|stating|production)")
	flag.StringVar(&cfg.jwtSecret, "jwt-secret", "mysecretkey", "JWT secret string")
	flag.StringVar(&cfg.dsn, "db-dsn", "", "Postgres DB Connection string")

	flag.Parse()

	// Initialize an New Logger which writes to the os.stdOut
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("Database connection pool established")

	// Initialize the application struct by providing the required fields
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	// Initialize the http.Server by providing custom configurations
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		// This ensures consistent error log in all over our application
		// Using the same log which we used in the application struct
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	// Calling ListenAndServe on our custom server, This Block the Control and Listen for Http Requests
	err = srv.ListenAndServe()

	// Return Error in scenarios Like Shutdown or some other un intented failures
	// Log the error message and exit the program
	logger.Error(err.Error())
	os.Exit(1)
}
