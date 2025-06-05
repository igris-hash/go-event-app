package database

// User queries
const (
	CreateUsersTable = `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`

	InsertUser = `
		INSERT INTO users (username, email, password)
		VALUES (?, ?, ?)
	`

	GetUserByEmail = `
		SELECT id, username, email, password
		FROM users
		WHERE email = ?
	`

	GetUserByID = `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = ?
	`
)

// Event queries
const (
	CreateEventsTable = `
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			location TEXT NOT NULL,
			date TIMESTAMP NOT NULL,
			capacity INTEGER NOT NULL,
			creator_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`

	InsertEvent = `
		INSERT INTO events (title, description, location, date, capacity, creator_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	GetEventByID = `
		SELECT e.id, e.title, e.description, e.location, e.date, e.capacity,
			   e.creator_id, e.created_at, u.username as creator_name
		FROM events e
		JOIN users u ON e.creator_id = u.id
		WHERE e.id = ?
	`

	ListEvents = `
		SELECT e.id, e.title, e.description, e.location, e.date, e.capacity,
			   e.creator_id, e.created_at, u.username as creator_name
		FROM events e
		JOIN users u ON e.creator_id = u.id
		ORDER BY e.date DESC
		LIMIT ? OFFSET ?
	`

	UpdateEvent = `
		UPDATE events
		SET title = ?, description = ?, location = ?, date = ?, capacity = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND creator_id = ?
	`

	DeleteEvent = `
		DELETE FROM events
		WHERE id = ? AND creator_id = ?
	`
)

// Registration queries
const (
	CreateRegistrationsTable = `
		CREATE TABLE IF NOT EXISTS registrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			status TEXT NOT NULL CHECK(status IN ('pending', 'confirmed', 'cancelled')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(event_id, user_id)
		)
	`

	RegisterForEvent = `
		INSERT INTO registrations (event_id, user_id, status)
		VALUES (?, ?, 'pending')
	`

	UpdateRegistrationStatus = `
		UPDATE registrations
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE event_id = ? AND user_id = ?
	`

	GetEventRegistrations = `
		SELECT r.id, r.status, r.created_at, u.id as user_id, u.username, u.email
		FROM registrations r
		JOIN users u ON r.user_id = u.id
		WHERE r.event_id = ?
		ORDER BY r.created_at DESC
	`

	GetUserRegistrations = `
		SELECT r.id, r.status, r.created_at,
			   e.id as event_id, e.title, e.date, e.location
		FROM registrations r
		JOIN events e ON r.event_id = e.id
		WHERE r.user_id = ?
		ORDER BY e.date DESC
	`

	GetRegistrationCount = `
		SELECT COUNT(*)
		FROM registrations
		WHERE event_id = ? AND status = 'confirmed'
	`
)
