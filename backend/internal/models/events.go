package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id         uuid.UUID `json:"id"`
	CalendarId uuid.UUID `json:"calendarId"`
	Name       string    `json:"name"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
}

type AddEvent struct {
	Name       string    `json:"name" validate:"required,max=100"`
	CalendarId uuid.UUID `json:"calendarId" validate:"required,uuid"`
	StartTime  time.Time `json:"startTime" validate:"required"`
	EndTime    time.Time `json:"endTime" validate:"required,gtfield=StartTime"`
}

type EventEdit struct {
	CalendarId *uuid.UUID `json:"calendarId,omitempty" validate:"omitempty,uuid"`
	Name       *string    `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	StartTime  *time.Time `json:"startTime,omitempty" validate:"omitempty"`
	EndTime    *time.Time `json:"endTime,omitempty" validate:"omitempty,gtfield=StartTime"`
}

type GetEventsParams struct {
	StartTime time.Time `query:"startTime" validate:"required"`
	EndTime   time.Time `query:"endTime" validate:"required,gtfield=StartTime"`
}

type SearchEventsParams struct {
	Name string `query:"name" validate:"required,max=100"`
}

type BatchDeleteEvents struct {
	Ids []uuid.UUID `json:"ids" validate:"required,min=1,dive,required,uuid"`
}
