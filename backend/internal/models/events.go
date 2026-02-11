package models

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	Name      string   `json:"name"`
	Address   *string  `json:"address,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

type Event struct {
	Id         uuid.UUID `json:"id"`
	CalendarId uuid.UUID `json:"calendarId"`
	Name       string    `json:"name"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Location   *Location `json:"location,omitempty"`
}

type LocationInput struct {
	Name      string   `json:"name" validate:"required,max=100"`
	Address   *string  `json:"address" validate:"omitempty,max=100"`
	Latitude  *float64 `json:"latitude" validate:"omitempty,latitude"`
	Longitude *float64 `json:"longitude" validate:"omitempty,longitude"`
}

type AddEvent struct {
	CalendarId uuid.UUID      `json:"calendarId" validate:"required,uuid"`
	Name       string         `json:"name" validate:"required,max=100"`
	StartTime  time.Time      `json:"startTime" validate:"required"`
	EndTime    time.Time      `json:"endTime" validate:"required,gtfield=StartTime"`
	Location   *LocationInput `json:"location,omitempty" validate:"omitempty"`
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
