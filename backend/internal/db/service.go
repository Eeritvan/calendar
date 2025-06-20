package db

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	sqlc "github.com/eeritvan/calendar/internal/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

type EventNotification struct {
	Type   string `json:"type"`
	Table  string `json:"table"`
	Action string `json:"action"`
	Data   Event  `json:"data"`
}

type DBService struct {
	Queries *sqlc.Queries
	Pool    *pgxpool.Pool
}

func initSchema(ctx context.Context, pool *pgxpool.Pool) error {
	sqlFile, err := os.ReadFile("./schema.sql")
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, string(sqlFile))
	return err
}

func ConnectToDB(ctx context.Context) (*DBService, error) {
	dbUrl := os.Getenv("DB_URL")

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}

	if err := initSchema(ctx, pool); err != nil {
		pool.Close()
		return nil, err
	}

	return &DBService{
		Queries: sqlc.New(pool),
		Pool:    pool,
	}, nil
}

// https://github.com/jackc/pgx/issues/1121 << useful stuff
func (s *DBService) Listen(ctx context.Context, channel string, callback func(EventNotification)) error {
	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, "LISTEN "+channel)
	if err != nil {
		conn.Release()
		return err
	}

	go func() {
		defer conn.Release()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				notification, err := conn.Conn().WaitForNotification(ctx)
				if err != nil {
					log.Printf("Error waiting for notification: %v", err)
					return
				}

				var eventNotification EventNotification
				if err := json.Unmarshal([]byte(notification.Payload), &eventNotification); err != nil {
					log.Printf("Error unmarshaling notification payload: %v", err)
					continue
				}

				eventNotification.Data.StartTime = eventNotification.Data.StartTime.Local()
				eventNotification.Data.EndTime = eventNotification.Data.EndTime.Local()

				callback(eventNotification)
			}
		}
	}()

	return nil
}
