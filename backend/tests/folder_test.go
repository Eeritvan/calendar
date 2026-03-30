package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/eeritvan/calendar/internal/models"
	"github.com/eeritvan/calendar/internal/utils"
	"github.com/google/uuid"

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

	userId := seedUser(t, ctx, queries, "addFolderUser", "password1")

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
		{
			name: "adding new folder with no name returns bad request",
			body: models.AddFolder{
				Name: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "adding new folder with too long name (>100 chars) returns bad request",
			body: models.AddFolder{
				Name: strings.Repeat("a", 101),
			},
			expectedStatus: http.StatusBadRequest,
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

func TestAddCalendarToFolder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI := spawnPostgresContainer(t, "folders")
	server, queries := setupTestServer(t, ctx, connURI)

	userId := seedUser(t, ctx, queries, "addCalendarToFolderUser1", "password1")
	calendarId := seedCalendar(t, ctx, queries, "calendar 1", userId)
	folderId := seedFolder(t, ctx, queries, "folder 1", userId)

	tests := []struct {
		name             string
		folderId         uuid.UUID
		calendarId       uuid.UUID
		expectedStatus   int
		expectedRespData models.Calendar
	}{
		{
			name:           "adding calendar to folder works",
			folderId:       folderId,
			calendarId:     calendarId,
			expectedStatus: http.StatusOK,
			expectedRespData: models.Calendar{
				Id:         calendarId,
				Name:       "calendar 1",
				OwnerId:    userId,
				Visibility: models.VisibilityPrivate,
				// Permission: models.PermissionWrite,
				IsOwner: true,
				Folder: &models.Folder{
					Id:   folderId,
					Name: "folder 1",
				},
			},
		},
		// {
		// 	name:           "calendars that are not own cannot be added to folders",
		// 	folderId:       folderId,
		// 	calendarId:     calendarId,
		// 	expectedStatus: http.StatusOK,
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, rec := echotest.ContextConfig{
				PathValues: echo.PathValues{
					{Name: "calendarId", Value: calendarId.String()},
					{Name: "folderId", Value: folderId.String()},
				},
				Headers: map[string][]string{
					echo.HeaderContentType: {echo.MIMEApplicationJSON},
				},
			}.ToContextRecorder(t)
			c.Echo().Validator = &utils.CustomValidator{
				Validator: validator.New(validator.WithRequiredStructEnabled()),
			}
			c.Set("userId", userId)

			_ = server.AddCalendarToFolder(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var got models.Calendar
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedRespData, got)
		})
	}
}

func TestEditFolder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI := spawnPostgresContainer(t, "folders")
	server, queries := setupTestServer(t, ctx, connURI)

	userId := seedUser(t, ctx, queries, "editFolderUser", "password1")
	folderId := seedFolder(t, ctx, queries, "folder 1", userId)

	tests := []struct {
		name             string
		folderId         uuid.UUID
		body             models.FolderEdit
		expectedStatus   int
		expectedRespData models.Folder
	}{
		{
			name:     "editing folder works",
			folderId: folderId,
			body: models.FolderEdit{
				Name: new("edited folder"),
			},
			expectedStatus: http.StatusOK,
			expectedRespData: models.Folder{
				Id:   folderId,
				Name: "edited folder",
			},
		},
		// {
		// 	name:           "other users folders cannot be edited",
		// 	folderId:       folderId,
		// 	calendarId:     calendarId,
		// 	expectedStatus: http.StatusOK,
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bodyJSON, err := json.Marshal(tc.body)
			require.NoError(t, err)

			c, rec := echotest.ContextConfig{
				PathValues: echo.PathValues{
					{Name: "folderId", Value: folderId.String()},
				},
				Headers: map[string][]string{
					echo.HeaderContentType: {echo.MIMEApplicationJSON},
				},
				JSONBody: bodyJSON,
			}.ToContextRecorder(t)
			c.Echo().Validator = &utils.CustomValidator{
				Validator: validator.New(validator.WithRequiredStructEnabled()),
			}
			c.Set("userId", userId)

			_ = server.EditFolder(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var got models.Folder
			err = json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedRespData, got)
		})
	}
}

// // TOOD: adding new calendar with folderId already required before this test can be created
// func TestRemoveCalendarFromFolder(t *testing.T) {
// 	t.Parallel()

// 	ctx := context.Background()
// 	connURI := spawnPostgresContainer(t, "folders")
// 	server, queries := setupTestServer(t, ctx, connURI)

// 	userId := seedUser(t, ctx, queries, "removeCalendarFromFolder", "password")
// 	folderId := seedFolder(t, ctx, queries, "folder 1", userId)

// 	tests := []struct {
// 		name           string
// 		folderId       uuid.UUID
// 		expectedStatus int
// 	}{
// 		{
// 			name:           "",
// 			folderId:       folderId,
// 			expectedStatus: http.StatusNoContent,
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()

// 			c, rec := echotest.ContextConfig{
// 				PathValues: echo.PathValues{
// 					{Name: "folderId", Value: folderId.String()},
// 				},
// 				Headers: map[string][]string{
// 					echo.HeaderContentType: {echo.MIMEApplicationJSON},
// 				},
// 			}.ToContextRecorder(t)
// 			c.Set("userId", userId)

// 			_ = server.DeleteFolder(c)

// 			assert.Equal(t, tc.expectedStatus, rec.Code)
// 		})
// 	}
// }

func TestDeleteFolder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI := spawnPostgresContainer(t, "folders")
	server, queries := setupTestServer(t, ctx, connURI)

	userId := seedUser(t, ctx, queries, "deleteFolderUser", "password")
	folderId := seedFolder(t, ctx, queries, "folder 1", userId)

	randomUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	tests := []struct {
		name           string
		folderId       uuid.UUID
		expectedStatus int
	}{
		{
			name:           "deleting folder works",
			folderId:       folderId,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "deleting folder that does not exist does not fail",
			folderId:       randomUUID,
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, rec := echotest.ContextConfig{
				PathValues: echo.PathValues{
					{Name: "folderId", Value: folderId.String()},
				},
				Headers: map[string][]string{
					echo.HeaderContentType: {echo.MIMEApplicationJSON},
				},
			}.ToContextRecorder(t)
			c.Set("userId", userId)

			_ = server.DeleteFolder(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
