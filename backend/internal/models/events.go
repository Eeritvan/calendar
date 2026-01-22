package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	CalendarId uuid.UUID `json:"calendarId"`
	EndTime    time.Time `json:"endTime"`
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	StartTime  time.Time `json:"startTime"`
}

type AddEvent struct {
	CalendarId uuid.UUID `json:"calendarId"`
	EndTime    time.Time `json:"endTime"`
	Name       string    `json:"name"`
	StartTime  time.Time `json:"startTime"`
}

type EventEdit struct {
	CalendarId *uuid.UUID `json:"calendarId,omitempty"`
	EndTime    *time.Time `json:"endTime,omitempty"`
	Name       *string    `json:"name,omitempty"`
	StartTime  *time.Time `json:"startTime,omitempty"`
}

type GetGetEventsParams struct {
	StartTime time.Time `query:"startTime"`
	EndTime   time.Time `query:"endTime"`
}

type GetSearchEventsParams struct {
	Name string `query:"name"`
}
