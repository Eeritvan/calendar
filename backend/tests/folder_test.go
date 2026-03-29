package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/eeritvan/calendar/internal/models"
	"github.com/eeritvan/calendar/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddFolder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI := spawnPostgresContainer(t, "folders")
	server, queries := setupTestServer(t, ctx, connURI)

	userId := seedUser(t, ctx, queries, "addFolderUser1", "password1")

	tests := []struct {
		name             string
		body             models.AddFolder
		expectedStatus   int
		expectedRespData models.Folder
	}{
		{
			name: "adding new folder works",
			body: models.AddFolder{
				Name: "work",
			},
			expectedStatus: http.StatusCreated,
			expectedRespData: models.Folder{
				// id is unknown beforehand
				Name: "work",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bodyJSON, err := json.Marshal(tc.body)
			require.NoError(t, err)

			c, rec := echotest.ContextConfig{
				Headers: map[string][]string{
					echo.HeaderContentType: {echo.MIMEApplicationJSON},
				},
				JSONBody: bodyJSON,
			}.ToContextRecorder(t)
			c.Echo().Validator = &utils.CustomValidator{
				Validator: validator.New(validator.WithRequiredStructEnabled()),
			}
			c.Set("userId", userId)

			_ = server.NewFolder(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var got models.Folder
			err = json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)

			assert.NotNil(t, got.Id)
			assert.Equal(t, tc.expectedRespData.Name, got.Name)
		})
	}
}
