// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createEvent = `-- name: CreateEvent :one
INSERT INTO events (name, description, start_time, end_time, color)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, description, start_time, end_time, color
`

type CreateEventParams struct {
	Name        string
	Description *string
	StartTime   time.Time
	EndTime     time.Time
	Color       EventColor
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error) {
	row := q.db.QueryRow(ctx, createEvent,
		arg.Name,
		arg.Description,
		arg.StartTime,
		arg.EndTime,
		arg.Color,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.StartTime,
		&i.EndTime,
		&i.Color,
	)
	return i, err
}

const deleteEvent = `-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1
`

func (q *Queries) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteEvent, id)
	return err
}

const eventsByTimeRange = `-- name: EventsByTimeRange :many
SELECT id, name, description, start_time, end_time, color
FROM events
WHERE start_time < $1 AND end_time > $2
`

type EventsByTimeRangeParams struct {
	StartTime time.Time
	EndTime   time.Time
}

func (q *Queries) EventsByTimeRange(ctx context.Context, arg EventsByTimeRangeParams) ([]Event, error) {
	rows, err := q.db.Query(ctx, eventsByTimeRange, arg.StartTime, arg.EndTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.StartTime,
			&i.EndTime,
			&i.Color,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listEvents = `-- name: ListEvents :many
SELECT id, name, description, start_time, end_time, color
FROM events
`

func (q *Queries) ListEvents(ctx context.Context) ([]Event, error) {
	rows, err := q.db.Query(ctx, listEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.StartTime,
			&i.EndTime,
			&i.Color,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateEvent = `-- name: UpdateEvent :one
UPDATE events
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    start_time = COALESCE($4, start_time),
    end_time = COALESCE($5, end_time),
    color = COALESCE($6, color)
WHERE id = $1
RETURNING id, name, description, start_time, end_time, color
`

type UpdateEventParams struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	StartTime   *time.Time
	EndTime     *time.Time
	Color       *EventColor
}

func (q *Queries) UpdateEvent(ctx context.Context, arg UpdateEventParams) (Event, error) {
	row := q.db.QueryRow(ctx, updateEvent,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.StartTime,
		arg.EndTime,
		arg.Color,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.StartTime,
		&i.EndTime,
		&i.Color,
	)
	return i, err
}
