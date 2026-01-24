package models

import "github.com/google/uuid"

type Calendar struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	OwnerId uuid.UUID `json:"ownerId"`
}

type AddCalendar struct {
	Name string `json:"name" validate:"required,max=100"`
}

type EditCalendar struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
}
