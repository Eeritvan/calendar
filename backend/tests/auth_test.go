package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/eeritvan/calendar/internal/models"
	"github.com/eeritvan/calendar/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignup(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI := spawnPostgresContainer(t, "users")
	server, queries := setupTestServer(t, ctx, connURI)

	seedUser(t, ctx, queries, "signupUser 3", "password 1")

	tests := []struct {
		name           string
		signup         models.Signup
		expectedStatus int
	}{
		{
			name: "user signup works",
			signup: models.Signup{
				Name:                 "signupUser 1",
				Password:             "password",
				PasswordConfirmation: "password",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "signup fails with mismatched passwords",
			signup: models.Signup{
				Name:                 "signupUser 2",
				Password:             "password",
				PasswordConfirmation: "wrong",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "signup fails with when name is already in use",
			signup: models.Signup{
				Name:                 "signupUser 3",
				Password:             "password",
				PasswordConfirmation: "password",
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "signup fails with too short password (< 8 characters)",
			signup: models.Signup{
				Name:                 "signupUser 4",
				Password:             "secret1",
				PasswordConfirmation: "secret1",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "signup fails with missing username",
			signup: models.Signup{
				Name:                 "",
				Password:             "password",
				PasswordConfirmation: "password",
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
	connURI := spawnPostgresContainer(t, "users")
	server, queries := setupTestServer(t, ctx, connURI)

	seedUser(t, ctx, queries, "loginUser1", "password 1")

	tests := []struct {
		name           string
		login          models.Login
		expectedStatus int
	}{
		{
			name: "user login works",
			login: models.Login{
				Name:     "loginUser1",
				Password: "password 1",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "login fails with incorrect password",
			login: models.Login{
				Name:     "loginUser1",
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

func TestLogout(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI := spawnPostgresContainer(t, "users")
	server, _ := setupTestServer(t, ctx, connURI)

	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "logout works and clears the cookie",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, rec := echotest.ContextConfig{
				Headers: map[string][]string{
					echo.HeaderContentType: {echo.MIMEApplicationJSON},
				},
			}.ToContextRecorder(t)
			c.Echo().Validator = &utils.CustomValidator{
				Validator: validator.New(validator.WithRequiredStructEnabled()),
			}

			_ = server.Logout(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedStatus == http.StatusOK {
				cookies := rec.Result().Cookies()
				assert.NotEmpty(t, cookies)

				var accessTokenCookie *http.Cookie
				for _, cookie := range cookies {
					if cookie.Name == "access_token" {
						accessTokenCookie = cookie
						break
					}
				}

				require.NotNil(t, accessTokenCookie, "access_token cookie should be present")
				assert.Equal(t, "", accessTokenCookie.Value, "cookie value should be empty")
				assert.Equal(t, -1, accessTokenCookie.MaxAge, "MaxAge should be -1")
				assert.True(t, accessTokenCookie.Expires.Before(time.Now().UTC()), "Expires should be in the past")
			}
		})
	}
}
