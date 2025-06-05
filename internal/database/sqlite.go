package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// SQLiteDB implements the Database interface for SQLite
type SQLiteDB struct {
	db     *sql.DB
	config *Config
}

// NewSQLiteDB creates a new SQLite database instance
func NewSQLiteDB(config *Config) *SQLiteDB {
	return &SQLiteDB{
		config: config,
	}
}

// Connect establishes a connection to the SQLite database
func (s *SQLiteDB) Connect() error {
	// Create database directory if it doesn't exist
	if err := os.MkdirAll(s.config.Host, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	dbPath := filepath.Join(s.config.Host, s.config.DBName)
	db, err := sql.Open("sqlite", dbPath) // Note: driver name is "sqlite" for modernc.org/sqlite
	if err != nil {
		return err
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return err
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	s.db = db
	return nil
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// BeginTx starts a new transaction
func (s *SQLiteDB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, nil)
}

// Exec executes a query without returning any rows
func (s *SQLiteDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

// Query executes a query that returns rows
func (s *SQLiteDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

// QueryRow executes a query that returns a single row
func (s *SQLiteDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.db.QueryRow(query, args...)
}

// DB returns the underlying sql.DB instance
func (s *SQLiteDB) DB() *sql.DB {
	return s.db
}
