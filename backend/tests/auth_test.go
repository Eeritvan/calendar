package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/eeritvan/calendar/internal/api"
	"github.com/eeritvan/calendar/internal/models"
	"github.com/eeritvan/calendar/internal/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignup(t *testing.T) {
	t.Setenv("JWT_KEY", "test_secret_key")

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "users")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	queries := sqlc.New(pool)
	server := api.NewServer(queries, pool, nil)

	seedUser(t, ctx, queries, "user 3", "password 1")

	tests := []struct {
		name           string
		signup         models.Signup
		expectedStatus int
	}{
		{
			name: "user signup works",
			signup: models.Signup{
				Name:                 "user 1",
				Password:             "password1",
				PasswordConfirmation: "password1",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "signup fails with mismatched passwords",
			signup: models.Signup{
				Name:                 "user 2",
				Password:             "password1",
				PasswordConfirmation: "wrong",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "signup fails with when name is already in use",
			signup: models.Signup{
				Name:                 "user 3",
				Password:             "password1",
				PasswordConfirmation: "password1",
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userJSON, err := json.Marshal(tc.signup)
			require.NoError(t, err)

			c, rec := echotest.ContextConfig{
				Headers: map[string][]string{
					echo.HeaderContentType: {echo.MIMEApplicationJSON},
				},
				JSONBody: userJSON,
			}.ToContextRecorder(t)

			_ = server.Signup(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedStatus == http.StatusOK {
				cookies := rec.Result().Cookies()
				assert.NotEmpty(t, cookies)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	t.Setenv("JWT_KEY", "test_secret_key")

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "users2")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	queries := sqlc.New(pool)
	server := api.NewServer(queries, pool, nil)

	seedUser(t, ctx, queries, "user 1", "password 1")

	tests := []struct {
		name           string
		login          models.Login
		expectedStatus int
	}{
		{
			name: "user login works",
			login: models.Login{
				Name:     "user 1",
				Password: "password 1",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "login fails with incorrect password",
			login: models.Login{
				Name:     "user 1",
				Password: "password 2",
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userJSON, err := json.Marshal(tc.login)
			require.NoError(t, err)

			c, rec := echotest.ContextConfig{
				Headers: map[string][]string{
					echo.HeaderContentType: {echo.MIMEApplicationJSON},
				},
				JSONBody: userJSON,
			}.ToContextRecorder(t)

			_ = server.Login(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedStatus == http.StatusOK {
				cookies := rec.Result().Cookies()
				assert.NotEmpty(t, cookies)
			}
		})
	}
}
