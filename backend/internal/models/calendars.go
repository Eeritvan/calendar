package models

import "github.com/google/uuid"

type Calendar struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	OwnerId uuid.UUID `json:"ownerId"`
}

type AddCalendar struct {
	Name string `json:"name"`
}

type CalendarEdit struct {
	Name *string `json:"name,omitempty"`
}
