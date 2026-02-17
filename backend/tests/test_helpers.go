package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/eeritvan/calendar/internal/api"
	"github.com/eeritvan/calendar/internal/models"
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/eeritvan/calendar/internal/stream"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

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

func buildInitScripts(t *testing.T, dir string) []string {
	t.Helper()

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)

	var migrationFiles []string
	for _, e := range entries {
		migrationFiles = append(migrationFiles, filepath.Join(dir, e.Name()))
	}
	sort.Strings(migrationFiles)

	tmpScripts := make([]string, 0, len(migrationFiles))
	cleanup := func() {
		for _, p := range tmpScripts {
			_ = os.Remove(p)
		}
	}

	for i, path := range migrationFiles {
		b, err := os.ReadFile(path)
		if err != nil {
			cleanup()
			require.NoError(t, err)
		}

		sql := string(b)
		if idx := strings.Index(sql, "-- +goose Down"); idx >= 0 {
			sql = sql[:idx]
		}

		f, err := os.CreateTemp("", fmt.Sprintf("init-%03d-*.sql", i))
		if err != nil {
			cleanup()
			require.NoError(t, err)
		}

		if _, err := f.WriteString(sql); err != nil {
			_ = f.Close()
			_ = os.Remove(f.Name())
			cleanup()
			require.NoError(t, err)
		}
		if err := f.Close(); err != nil {
			_ = os.Remove(f.Name())
			cleanup()
			require.NoError(t, err)
		}

		tmpScripts = append(tmpScripts, f.Name())
	}

	t.Cleanup(cleanup)
	return tmpScripts
}

func spawnPostgresContainer(t *testing.T, reuseName string) string {
	t.Helper()

	ctx := context.Background()

	migrationsDir := filepath.Join("../migrations")
	filteredScripts := buildInitScripts(t, migrationsDir)

	postgresContainer, err := postgres.Run(context.Background(),
		"postgres:alpine",
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithOrderedInitScripts(filteredScripts...),
		testcontainers.WithReuseByName(reuseName),
		testcontainers.WithWaitStrategy(
			wait.
				ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	require.NoError(t, err)

	connURI, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	return connURI
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

func seedEvent(t *testing.T, ctx context.Context, queries *sqlc.Queries, ownerID uuid.UUID, body models.AddEvent) uuid.UUID {
	t.Helper()

	var locationName string
	var locationAddress *string
	var lat *float64
	var lng *float64

	if body.Location != nil {
		locationName = body.Location.Name
		locationAddress = body.Location.Address
		lat = body.Location.Latitude
		lng = body.Location.Longitude
	}

	event, err := queries.AddEvent(ctx, sqlc.AddEventParams{
		CalendarID:   body.CalendarId,
		Name:         body.Name,
		OwnerID:      ownerID,
		StartTime:    body.StartTime,
		EndTime:      body.EndTime,
		LocationName: locationName,
		Address:      locationAddress,
		Longitude:    lng,
		Latitude:     lat,
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
