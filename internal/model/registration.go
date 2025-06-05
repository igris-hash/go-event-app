package model

import "time"

type RegistrationStatus string

const (
	StatusPending   RegistrationStatus = "pending"
	StatusConfirmed RegistrationStatus = "confirmed"
	StatusCancelled RegistrationStatus = "cancelled"
)

type Registration struct {
	ID        int                `json:"id"`
	EventID   int                `json:"event_id"`
	UserID    int                `json:"user_id"`
	Status    RegistrationStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`

	// Joined fields
	Username   string    `json:"username,omitempty"`
	Email      string    `json:"email,omitempty"`
	EventTitle string    `json:"event_title,omitempty"`
	EventDate  time.Time `json:"event_date,omitempty"`
	Location   string    `json:"location,omitempty"`
}

type RegistrationUpdate struct {
	Status RegistrationStatus `json:"status" binding:"required,oneof=pending confirmed cancelled"`
}
