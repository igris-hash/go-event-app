package repository

import (
	"errors"

	"github.com/igris-hash/go-event-app/internal/database"
	"github.com/igris-hash/go-event-app/internal/model"
)

type RegistrationRepository struct {
	db database.Database
}

func NewRegistrationRepository(db database.Database) *RegistrationRepository {
	return &RegistrationRepository{db: db}
}

func (r *RegistrationRepository) InitTable() error {
	_, err := r.db.Exec(database.CreateRegistrationsTable)
	return err
}

func (r *RegistrationRepository) Create(eventID, userID int) error {
	result, err := r.db.Exec(database.RegisterForEvent, eventID, userID)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	return err
}

func (r *RegistrationRepository) UpdateStatus(eventID, userID int, status model.RegistrationStatus) error {
	result, err := r.db.Exec(database.UpdateRegistrationStatus, status, eventID, userID)
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

func (r *RegistrationRepository) GetEventRegistrations(eventID int) ([]*model.Registration, error) {
	rows, err := r.db.Query(database.GetEventRegistrations, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registrations []*model.Registration
	for rows.Next() {
		reg := &model.Registration{}
		err := rows.Scan(
			&reg.ID,
			&reg.Status,
			&reg.CreatedAt,
			&reg.UserID,
			&reg.Username,
			&reg.Email,
		)
		if err != nil {
			return nil, err
		}
		reg.EventID = eventID
		registrations = append(registrations, reg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return registrations, nil
}

func (r *RegistrationRepository) GetUserRegistrations(userID int) ([]*model.Registration, error) {
	rows, err := r.db.Query(database.GetUserRegistrations, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registrations []*model.Registration
	for rows.Next() {
		reg := &model.Registration{}
		err := rows.Scan(
			&reg.ID,
			&reg.Status,
			&reg.CreatedAt,
			&reg.EventID,
			&reg.EventTitle,
			&reg.EventDate,
			&reg.Location,
		)
		if err != nil {
			return nil, err
		}
		reg.UserID = userID
		registrations = append(registrations, reg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return registrations, nil
}

func (r *RegistrationRepository) GetRegistrationCount(eventID int) (int, error) {
	var count int
	err := r.db.QueryRow(database.GetRegistrationCount, eventID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
