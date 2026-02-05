package models

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Event struct {
	Id         uuid.UUID `json:"id"`
	CalendarId uuid.UUID `json:"calendarId"`
	Name       string    `json:"name"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Location   *Location `json:"location"`
}

type LocationInput struct {
	Name      string  `json:"name" validate:"max=100"`
	Address   string  `json:"address" validate:"max=100"`
	Latitude  float64 `json:"latitude" validate:"latitude"`
	Longitude float64 `json:"longitude" validate:"longitude"`
}

type AddEvent struct {
	CalendarId uuid.UUID     `json:"calendarId" validate:"required,uuid"`
	Name       string        `json:"name" validate:"required,max=100"`
	StartTime  time.Time     `json:"startTime" validate:"required"`
	EndTime    time.Time     `json:"endTime" validate:"required,gtfield=StartTime"`
	Location   LocationInput `json:"location" validate:"required"`
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
