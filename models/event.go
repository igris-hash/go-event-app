package models

import "time"

type Event struct {
	ID          int
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int
}

var eventsBox = []Event{}

func (e Event) Save() {
	// later: add events into Database
	eventsBox = append(eventsBox, e)
}

func GetEvents() []Event {
	return eventsBox
}
