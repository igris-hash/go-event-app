package utils

import (
	"database/sql"
	"log"
	"os"
)

// InitializeDB creates the database file if it doesn't exist and sets up the tables
func InitializeDB(dbPath string) (*sql.DB, error) {
	// Check if database file exists
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		log.Printf("Creating new database file at %s\n", dbPath)
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Create tables
	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create events table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			date DATETIME NOT NULL,
			capacity INTEGER NOT NULL,
			creator_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (creator_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	// Create registrations table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS registrations (
			user_id INTEGER NOT NULL,
			event_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			PRIMARY KEY (user_id, event_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (event_id) REFERENCES events(id)
		)
	`)
	return err
}
