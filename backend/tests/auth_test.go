package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/eeritvan/calendar/internal/models"
	"github.com/eeritvan/calendar/internal/utils"

	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignup(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "users")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	server, queries := setupTestServer(t, ctx, connURI)

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
				Password:             "password",
				PasswordConfirmation: "password",
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
		{
			name: "signup fails with too short password (< 8 characters)",
			signup: models.Signup{
				Name:                 "user 4",
				Password:             "secret1",
				PasswordConfirmation: "secret1",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "signup fails with missing username",
			signup: models.Signup{
				Name:                 "",
				Password:             "password1",
				PasswordConfirmation: "password1",
			},
			expectedStatus: http.StatusBadRequest,
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
			c.Echo().Validator = &utils.CustomValidator{
				Validator: validator.New(validator.WithRequiredStructEnabled()),
			}

			_ = server.Signup(c)
			fmt.Println(rec.Body)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedStatus == http.StatusOK {
				cookies := rec.Result().Cookies()
				assert.NotEmpty(t, cookies)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "users2")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	server, queries := setupTestServer(t, ctx, connURI)

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
			c.Echo().Validator = &utils.CustomValidator{
				Validator: validator.New(validator.WithRequiredStructEnabled()),
			}

			_ = server.Login(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedStatus == http.StatusOK {
				cookies := rec.Result().Cookies()
				assert.NotEmpty(t, cookies)
			}
		})
	}
}
