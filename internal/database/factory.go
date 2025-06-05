package database

import "fmt"

const (
	SQLite = "sqlite"
	// Add other database types here as needed
	// PostgreSQL = "postgres"
	// MySQL = "mysql"
)

// NewDatabase creates a new database instance based on the driver type
func NewDatabase(config *Config) (Database, error) {
	switch config.Driver {
	case SQLite:
		return NewSQLiteDB(config), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}
}
