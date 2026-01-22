package api_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/eeritvan/calendar/internal/api"
	"github.com/eeritvan/calendar/internal/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v5"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSignup(t *testing.T) {
	ctx := context.Background()
	postgresContainer, err := postgres.Run(context.Background(),
		"postgres:18-alpine",
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.
				ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	if err != nil {
		t.Fatal(err)
	}

	connURI, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("pgx", connURI)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatal(err)
	}

	if err := goose.Up(db, "../migrations"); err != nil {
		t.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, connURI)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	queries := sqlc.New(pool)
	server := api.NewServer(queries, pool, nil)

	// -----

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(`{"name":"Jon Snow","password":"1","password_confirmation":"1"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = server.PostSignup(c)
	fmt.Println(err)

	assert.Equal(t, http.StatusOK, rec.Code)
}
