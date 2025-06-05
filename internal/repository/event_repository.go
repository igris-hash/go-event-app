package repository

import (
	"database/sql"
	"errors"

	"github.com/igris-hash/go-event-app/internal/database"
	"github.com/igris-hash/go-event-app/internal/model"
)

// EventRepository handles database operations for events
type EventRepository struct {
	db *sql.DB
}

// NewEventRepository creates a new EventRepository instance
func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

// InitTable initializes the events table
func (r *EventRepository) InitTable() error {
	_, err := r.db.Exec(database.CreateEventsTable)
	return err
}

// Create inserts a new event into the database
func (r *EventRepository) Create(event *model.Event) error {
	query := `
		INSERT INTO events (title, description, location, date, capacity, creator_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		event.Title,
		event.Description,
		event.Location,
		event.Date,
		event.Capacity,
		event.CreatorID,
		event.CreatedAt,
		event.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	event.ID = int(id)
	return nil
}

// GetByID retrieves an event by ID
func (r *EventRepository) GetByID(id int) (*model.Event, error) {
	event := &model.Event{}
	query := `
		SELECT id, title, description, location, date, capacity, creator_id, created_at, updated_at
		FROM events
		WHERE id = ?
	`
	err := r.db.QueryRow(query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Location,
		&event.Date,
		&event.Capacity,
		&event.CreatorID,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("event not found")
	}
	if err != nil {
		return nil, err
	}
	return event, nil
}

// GetAll retrieves all events from the database
func (r *EventRepository) GetAll() ([]*model.Event, error) {
	query := `
		SELECT id, title, description, location, date, capacity, creator_id, created_at, updated_at
		FROM events
		ORDER BY date DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.Event
	for rows.Next() {
		event := &model.Event{}
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Location,
			&event.Date,
			&event.Capacity,
			&event.CreatorID,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// Update updates event information
func (r *EventRepository) Update(event *model.Event) error {
	query := `
		UPDATE events
		SET title = ?, description = ?, location = ?, date = ?, capacity = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query,
		event.Title,
		event.Description,
		event.Location,
		event.Date,
		event.Capacity,
		event.UpdatedAt,
		event.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("event not found")
	}

	return nil
}

// Delete removes an event from the database
func (r *EventRepository) Delete(id int) error {
	query := `DELETE FROM events WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("event not found")
	}

	return nil
}

// RegisterUser registers a user for an event
func (r *EventRepository) RegisterUser(eventID, userID int) error {
	query := `
		INSERT INTO registrations (event_id, user_id, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
	`
	_, err := r.db.Exec(query, eventID, userID)
	return err
}

// UnregisterUser removes a user's registration for an event
func (r *EventRepository) UnregisterUser(eventID, userID int) error {
	query := `DELETE FROM registrations WHERE event_id = ? AND user_id = ?`
	result, err := r.db.Exec(query, eventID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("registration not found")
	}

	return nil
}

// GetRegisteredUsers retrieves all users registered for an event
func (r *EventRepository) GetRegisteredUsers(eventID int) ([]*model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.updated_at
		FROM users u
		JOIN registrations r ON r.user_id = u.id
		WHERE r.event_id = ?
	`
	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetRegisteredUsersCount retrieves the count of users registered for an event
func (r *EventRepository) GetRegisteredUsersCount(eventID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM registrations WHERE event_id = ?`
	err := r.db.QueryRow(query, eventID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// IsUserRegistered checks if a user is registered for an event
func (r *EventRepository) IsUserRegistered(eventID, userID int) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM registrations WHERE event_id = ? AND user_id = ?
		)
	`
	err := r.db.QueryRow(query, eventID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
