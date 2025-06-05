package service

import (
	"errors"
	"time"

	"github.com/igris-hash/go-event-app/internal/model"
	"github.com/igris-hash/go-event-app/internal/repository"
)

// EventService handles business logic for events
type EventService struct {
	repo *repository.EventRepository
}

// NewEventService creates a new EventService instance
func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

// Create creates a new event
func (s *EventService) Create(event model.Event) (*model.Event, error) {
	// Set timestamps
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	// Validate capacity
	if event.Capacity <= 0 {
		return nil, errors.New("capacity must be greater than 0")
	}

	// Create event
	err := s.repo.Create(&event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// GetByID retrieves an event by ID
func (s *EventService) GetByID(id int) (*model.Event, error) {
	return s.repo.GetByID(id)
}

// GetAll retrieves all events
func (s *EventService) GetAll() ([]*model.Event, error) {
	return s.repo.GetAll()
}

// Update updates event information
func (s *EventService) Update(event model.Event) (*model.Event, error) {
	// Verify event exists
	existing, err := s.repo.GetByID(event.ID)
	if err != nil {
		return nil, err
	}

	// Update fields
	existing.Title = event.Title
	existing.Description = event.Description
	existing.Location = event.Location
	existing.Date = event.Date
	existing.Capacity = event.Capacity
	existing.UpdatedAt = time.Now()

	// Save changes
	err = s.repo.Update(existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete removes an event
func (s *EventService) Delete(id int) error {
	// Check if event exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// RegisterUser registers a user for an event
func (s *EventService) RegisterUser(eventID, userID int) error {
	// Check if event exists and has capacity
	event, err := s.repo.GetByID(eventID)
	if err != nil {
		return err
	}

	// Check if user is already registered
	isRegistered, err := s.repo.IsUserRegistered(eventID, userID)
	if err != nil {
		return err
	}
	if isRegistered {
		return errors.New("user is already registered for this event")
	}

	// Check event capacity
	registeredCount, err := s.repo.GetRegisteredUsersCount(eventID)
	if err != nil {
		return err
	}
	if registeredCount >= event.Capacity {
		return errors.New("event is at full capacity")
	}

	// Register user
	return s.repo.RegisterUser(eventID, userID)
}

// UnregisterUser removes a user's registration for an event
func (s *EventService) UnregisterUser(eventID, userID int) error {
	// Check if user is registered
	isRegistered, err := s.repo.IsUserRegistered(eventID, userID)
	if err != nil {
		return err
	}
	if !isRegistered {
		return errors.New("user is not registered for this event")
	}

	return s.repo.UnregisterUser(eventID, userID)
}

// GetRegisteredUsers retrieves all users registered for an event
func (s *EventService) GetRegisteredUsers(eventID int) ([]*model.User, error) {
	// Check if event exists
	_, err := s.repo.GetByID(eventID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetRegisteredUsers(eventID)
}
