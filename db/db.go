package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	// Opening the database connection to api.db (this file will be created if it doesn't exist)
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		// It's generally better to log a fatal error and exit here
		// rather than panic, as panicking can obscure shutdown processes.
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Set maximum number of open connections to the database.
	DB.SetMaxOpenConns(10)
	// Set maximum number of idle connections in the connection pool.
	DB.SetMaxIdleConns(5)

	// Call the function to create necessary tables.
	createTables()
}

func createTables() {
	// SQL statement to create the 'users' table if it doesn't already exist.
	createUserTable := `
		CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
		)
	`
	// Execute the SQL statement.
	_, err := DB.Exec(createUserTable)
	if err != nil {
		log.Fatalf("Could not create users table: %v", err) // Log error with details
	}

	// SQL statement to create the 'events' table if it doesn't already exist.
	// FIX: Added a comma after 'user_id INTEGER'.
	createEventTable := `
		CREATE TABLE IF NOT EXISTS events(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
		)
	`
	// Execute the SQL statement.
	_, eventErr := DB.Exec(createEventTable)
	if eventErr != nil {
		log.Fatalf("Could not create events table: %v", eventErr) // Log error with details
	}
}
