package api_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
)

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

func seedUser(t *testing.T, ctx context.Context, queries *sqlc.Queries, name, password string) {
	t.Helper()

	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	require.NoError(t, err)

	_, err = queries.Signup(ctx, sqlc.SignupParams{
		Name:         name,
		PasswordHash: string(hashedPW),
	})
	require.NoError(t, err)
}

func clearDatabase(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name != 'goose_db_version';`

	rows, err := pool.Query(context.Background(), query)
	require.NoError(t, err)
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		require.NoError(t, err)
		tables = append(tables, table)
	}

	if len(tables) > 0 {
		truncateQuery := fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", strings.Join(tables, ", "))
		_, err := pool.Exec(context.Background(), truncateQuery)
		require.NoError(t, err)
	}
}
