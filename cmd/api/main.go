package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port      int    // Port number for http connection
	env       string // Environment(devolopment|staging|production)
	jwtSecret string // JWT Secret string, For generating signed JWT Strings
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config // Initialize an cfg struct

	// Scan the flags into the cfg struct by passing the references
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "devolopment", "Environment(devolopment|stating|production)")
	flag.StringVar(&cfg.jwtSecret, "jwt-secret", "mysecretkey", "JWT secret string")

	flag.Parse()

	// Initialize an New Logger which writes to the os.stdOut
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize the application struct by providing the required fields
	app := &application{
		config: cfg,
		logger: logger,
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
	err := srv.ListenAndServe()

	// Return Error in scenarios Like Shutdown or some other un intented failures
	// Log the error message and exit the program
	logger.Error(err.Error())
	os.Exit(1)
}
