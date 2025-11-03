package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// openDB creates an sql connection using the provided connection string
func openDB(dsn string) (*sql.DB, error) {
	// opens an connection and using postgres as the driver name
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(15 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the server to check if it available to use
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
