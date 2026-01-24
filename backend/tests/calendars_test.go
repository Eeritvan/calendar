package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/eeritvan/calendar/internal/api"
	"github.com/eeritvan/calendar/internal/models"
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCalendars(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "calendars0")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	queries := sqlc.New(pool)
	server := api.NewServer(queries, pool, nil)

	userId1 := seedUser(t, ctx, queries, "calendarUser", "password")
	seedCalendar(t, ctx, queries, "meetings", userId1)
	seedCalendar(t, ctx, queries, "daily", userId1)

	userId2 := seedUser(t, ctx, queries, "secondUser", "password")
	seedCalendar(t, ctx, queries, "video games", userId2)

	tests := []struct {
		name             string
		userId           uuid.UUID
		expectedStatus   int
		expectedRespData []models.Calendar
	}{
		{
			name:           "fetching calendars work",
			userId:         userId1,
			expectedStatus: http.StatusOK,
			expectedRespData: []models.Calendar{
				{Name: "meetings", OwnerId: userId1},
				{Name: "daily", OwnerId: userId1},
			},
		},
		{
			name:           "second user cannot see first user calendars",
			userId:         userId2,
			expectedStatus: http.StatusOK,
			expectedRespData: []models.Calendar{
				{Name: "video games", OwnerId: userId2},
			},
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
			c.Set("userId", tc.userId)

			_ = server.GetCalendars(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var got []models.Calendar
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)

			for i := range tc.expectedRespData {
				tc.expectedRespData[i].Id = got[i].Id
			}

			assert.ElementsMatch(t, tc.expectedRespData, got)
		})
	}
}

func TestAddCalendar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "calendars")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	queries := sqlc.New(pool)
	server := api.NewServer(queries, pool, nil)

	userId := seedUser(t, ctx, queries, "calendarUser", "password")

	tests := []struct {
		name             string
		body             models.AddCalendar
		expectedStatus   int
		expectedRespData models.Calendar
	}{
		{
			name: "adding calendar works",
			body: models.AddCalendar{
				Name: "meetings",
			},
			expectedStatus: http.StatusOK,
			expectedRespData: models.Calendar{
				// id is unknown beforehand
				Name:    "meetings",
				OwnerId: userId,
			},
		},
		// TODO
		// {
		// 	name:           "adding calendars fails with empty body",
		// 	body:           models.AddCalendar{},
		// 	expectedStatus: http.StatusBadRequest,
		// },
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
			c.Set("userId", userId)

			_ = server.AddCalendar(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var got models.Calendar
			err = json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)

			assert.NotNil(t, got.Id)
			assert.Equal(t, tc.expectedRespData.Name, got.Name)
			assert.Equal(t, tc.expectedRespData.OwnerId, got.OwnerId)
		})
	}
}

func TestEditCalendar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "calendars2")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	queries := sqlc.New(pool)
	server := api.NewServer(queries, pool, nil)

	userId := seedUser(t, ctx, queries, "calendarUser", "password")
	calendarId := seedCalendar(t, ctx, queries, "meetings", userId)

	userId2 := seedUser(t, ctx, queries, "calendarUser2", "password")
	calendarId2 := seedCalendar(t, ctx, queries, "meetings", userId2)

	randomUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	tests := []struct {
		name             string
		calendarId       uuid.UUID
		body             models.CalendarEdit
		expectedStatus   int
		expectedRespData models.Calendar
	}{
		{
			name:       "editing calendar works",
			calendarId: calendarId,
			body: models.CalendarEdit{
				Name: Ptr("daily"),
			},
			expectedStatus: http.StatusOK,
			expectedRespData: models.Calendar{
				Id:      calendarId,
				Name:    "daily",
				OwnerId: userId,
			},
		},
		{
			name:       "editing non-existent calendars fails",
			calendarId: randomUUID,
			body: models.CalendarEdit{
				Name: Ptr("daily"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:       "editing other users calendars fails",
			calendarId: calendarId2,
			body: models.CalendarEdit{
				Name: Ptr("daily"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bodyJSON, err := json.Marshal(tc.body)
			require.NoError(t, err)

			c, rec := echotest.ContextConfig{
				PathValues: echo.PathValues{
					{Name: "calendarId", Value: tc.calendarId.String()},
				},
				Headers: map[string][]string{
					echo.HeaderContentType: {echo.MIMEApplicationJSON},
				},
				JSONBody: bodyJSON,
			}.ToContextRecorder(t)
			c.Set("userId", userId)

			_ = server.EditCalendar(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var got models.Calendar
			err = json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedRespData, got)
		})
	}
}

func TestDeleteCalendar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "deleteCalendar")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	queries := sqlc.New(pool)
	server := api.NewServer(queries, pool, nil)

	userId := seedUser(t, ctx, queries, "calendarUser", "password")
	calendarId := seedCalendar(t, ctx, queries, "meetings", userId)

	randomUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	tests := []struct {
		name           string
		calendarId     uuid.UUID
		expectedStatus int
	}{
		{
			name:           "deleting calendar works",
			calendarId:     calendarId,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "deleting calendar that does not exist does not fail",
			calendarId:     randomUUID,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, rec := echotest.ContextConfig{
				PathValues: echo.PathValues{
					{Name: "calendarId", Value: calendarId.String()},
				},
				Headers: map[string][]string{
					echo.HeaderContentType: {echo.MIMEApplicationJSON},
				},
			}.ToContextRecorder(t)
			c.Set("userId", userId)

			_ = server.DeleteCalendar(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
