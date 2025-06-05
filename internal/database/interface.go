package database

import (
	"context"
	"database/sql"
)

// Database represents the interface that any database implementation must satisfy
type Database interface {
	// Connection management
	Connect() error
	Close() error

	// Transaction management
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// Basic operations
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row

	// Get underlying database
	DB() *sql.DB
}

// Config holds the database configuration
type Config struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}
