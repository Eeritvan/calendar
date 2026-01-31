package tests

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/eeritvan/calendar/internal/api"
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/eeritvan/calendar/internal/stream"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	os.Setenv("TZ", "UTC")
	os.Setenv("JWT_KEY", "test_secret")
	os.Exit(m.Run())
}

func Ptr[T any](v T) *T {
	return &v
}

func spawnPostgresContainer(t *testing.T, reuseName string) (string, error) {
	t.Helper()

	ctx := context.Background()
	postgresContainer, err := postgres.Run(context.Background(),
		"postgres:alpine",
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithReuseByName(reuseName),
		testcontainers.WithWaitStrategy(
			wait.
				ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	if err != nil {
		return "", err
	}

	connURI, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return "", err
	}
	return connURI, nil
}

func runMigrations(t *testing.T, connURI string) error {
	t.Helper()

	db, err := sql.Open("pgx", connURI)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "../migrations"); err != nil {
		return err
	}

	return nil
}

func seedUser(t *testing.T, ctx context.Context, queries *sqlc.Queries, name, password string) uuid.UUID {
	t.Helper()

	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	require.NoError(t, err)

	user, err := queries.Signup(ctx, sqlc.SignupParams{
		Name:         name,
		PasswordHash: string(hashedPW),
	})
	require.NoError(t, err)

	return user.ID
}

func seedCalendar(t *testing.T, ctx context.Context, queries *sqlc.Queries, name string, userId uuid.UUID) uuid.UUID {
	t.Helper()

	calendar, err := queries.AddCalendar(ctx, sqlc.AddCalendarParams{
		Name:    name,
		OwnerID: userId,
	})
	require.NoError(t, err)

	return calendar.ID
}

func seedEvent(t *testing.T, ctx context.Context, queries *sqlc.Queries, name string, ownerID, calendarId uuid.UUID, startTime, endTime time.Time) uuid.UUID {
	t.Helper()

	event, err := queries.AddEvent(ctx, sqlc.AddEventParams{
		Name:       name,
		OwnerID:    ownerID,
		CalendarID: calendarId,
		StartTime:  startTime,
		EndTime:    endTime,
	})
	require.NoError(t, err)

	return event.ID
}

func setupTestServer(t *testing.T, ctx context.Context, connURI string) (*api.Server, *sqlc.Queries) {
	t.Helper()

	pool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	queries := sqlc.New(pool)

	sseHandler := &stream.SSEHandler{
		SSEServer:   nil,
		UserClients: make(map[uuid.UUID]map[string]struct{}),
	}

	server := api.NewServer(queries, pool, sseHandler)

	return server, queries
}
