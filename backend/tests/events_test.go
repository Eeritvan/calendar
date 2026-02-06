package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/eeritvan/calendar/internal/models"
	"github.com/eeritvan/calendar/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEvents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "getEvents")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	server, queries := setupTestServer(t, ctx, connURI)

	timeNow := time.Now().UTC().Truncate(time.Microsecond)
	timePlusHour := time.Now().UTC().Add(time.Hour).Truncate(time.Microsecond)
	timeMinusHour := time.Now().UTC().Add(-1 * time.Hour).Truncate(time.Microsecond)

	userId1 := seedUser(t, ctx, queries, "user1", "password1")
	calendarId1 := seedCalendar(t, ctx, queries, "meetings", userId1)
	calendarId2 := seedCalendar(t, ctx, queries, "daily", userId1)
	eventId1 := seedEvent2(t, ctx, queries, userId1, models.AddEvent{
		CalendarId: calendarId1,
		Name:       "team meeting",
		StartTime:  timeMinusHour,
		EndTime:    timeNow,
		Location: models.LocationInput{
			Name:      "toimisto",
			Address:   "toimistokatu 1",
			Latitude:  60.248947411912596,
			Longitude: 24.978291441099014,
		},
	})
	eventId2 := seedEvent2(t, ctx, queries, userId1, models.AddEvent{
		CalendarId: calendarId2,
		Name:       "monday standup",
		StartTime:  timeNow,
		EndTime:    timePlusHour,
		Location: models.LocationInput{
			Name:      "toimisto",
			Address:   "toimistokatu 1",
			Latitude:  60.248947411912596,
			Longitude: 24.978291441099014,
		},
	})
	userId2 := seedUser(t, ctx, queries, "user2", "password2")
	calendarId3 := seedCalendar(t, ctx, queries, "video games", userId2)
	eventId3 := seedEvent2(t, ctx, queries, userId2, models.AddEvent{
		CalendarId: calendarId3,
		Name:       "my winter car",
		StartTime:  timeMinusHour,
		EndTime:    timePlusHour,
		Location: models.LocationInput{
			Name:      "koti",
			Address:   "kotikatu 1",
			Latitude:  61.248947411912596,
			Longitude: 23.978291441099014,
		},
	})

	tests := []struct {
		name             string
		userId           uuid.UUID
		queryStartTime   string
		queryEndTime     string
		expectedStatus   int
		expectedRespData []models.Event
	}{
		{
			name:           "fetching events from multiple calendars work",
			userId:         userId1,
			queryStartTime: timeMinusHour.Format(time.RFC3339),
			queryEndTime:   timePlusHour.Format(time.RFC3339),
			expectedStatus: http.StatusOK,
			expectedRespData: []models.Event{
				{
					Id:         eventId1,
					Name:       "team meeting",
					CalendarId: calendarId1,
					StartTime:  timeMinusHour,
					EndTime:    timeNow,
					Location: &models.Location{
						Name:      "toimisto",
						Address:   "toimistokatu 1",
						Latitude:  60.248947411912596,
						Longitude: 24.978291441099014,
					},
				},
				{
					Id:         eventId2,
					Name:       "monday standup",
					CalendarId: calendarId2,
					StartTime:  timeNow,
					EndTime:    timePlusHour,
					Location: &models.Location{
						Name:      "toimisto",
						Address:   "toimistokatu 1",
						Latitude:  60.248947411912596,
						Longitude: 24.978291441099014,
					},
				},
			},
		},
		{
			name:           "second user cannot see first user events",
			userId:         userId2,
			queryStartTime: timeMinusHour.Format(time.RFC3339),
			queryEndTime:   timePlusHour.Format(time.RFC3339),
			expectedStatus: http.StatusOK,
			expectedRespData: []models.Event{
				{
					Id:         eventId3,
					Name:       "my winter car",
					CalendarId: calendarId3,
					StartTime:  timeMinusHour,
					EndTime:    timePlusHour,
					Location: &models.Location{
						Name:      "koti",
						Address:   "kotikatu 1",
						Latitude:  61.248947411912596,
						Longitude: 23.978291441099014,
					},
				},
			},
		},
		{
			name:           "only events overlapping the query startTime and endTime are shown",
			userId:         userId1,
			queryStartTime: timeMinusHour.Format(time.RFC3339),
			queryEndTime:   timeNow.Format(time.RFC3339),
			expectedStatus: http.StatusOK,
			expectedRespData: []models.Event{
				{
					Id:         eventId1,
					Name:       "team meeting",
					CalendarId: calendarId1,
					StartTime:  timeMinusHour,
					EndTime:    timeNow,
					Location: &models.Location{
						Name:      "toimisto",
						Address:   "toimistokatu 1",
						Latitude:  60.248947411912596,
						Longitude: 24.978291441099014,
					},
				},
			},
		},
		{
			name:           "return error if startTime is after endTime",
			userId:         userId1,
			queryStartTime: timePlusHour.Format(time.RFC3339),
			queryEndTime:   timeMinusHour.Format(time.RFC3339),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, rec := echotest.ContextConfig{
				QueryValues: url.Values{
					"startTime": {tc.queryStartTime},
					"endTime":   {tc.queryEndTime},
				},
				Headers: map[string][]string{
					echo.HeaderContentType: {echo.MIMEApplicationJSON},
				},
			}.ToContextRecorder(t)
			c.Echo().Validator = &utils.CustomValidator{
				Validator: validator.New(validator.WithRequiredStructEnabled()),
			}
			c.Set("userId", tc.userId)

			_ = server.GetEvents(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var got []models.Event
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)

			if len(got) != len(tc.expectedRespData) {
				t.Fatalf("%v, expected %d events, got %d", tc.name, len(tc.expectedRespData), len(got))
			}

			for i := range tc.expectedRespData {
				tc.expectedRespData[i].Id = got[i].Id
			}

			assert.ElementsMatch(t, tc.expectedRespData, got)
		})
	}
}

func TestSearchEvents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "searchEvents")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	server, queries := setupTestServer(t, ctx, connURI)

	startTime := time.Now().UTC().Truncate(time.Microsecond)
	endTime := time.Now().UTC().Add(time.Hour).Truncate(time.Microsecond)

	userId1 := seedUser(t, ctx, queries, "user1", "password1")
	calendarId1 := seedCalendar(t, ctx, queries, "meetings", userId1)
	calendarId2 := seedCalendar(t, ctx, queries, "daily", userId1)
	eventId1 := seedEvent2(t, ctx, queries, userId1, models.AddEvent{
		CalendarId: calendarId1,
		Name:       "team meeting",
		StartTime:  startTime,
		EndTime:    endTime,
		Location: models.LocationInput{
			Name:      "toimisto",
			Address:   "toimistokatu 1",
			Latitude:  60.248947411912596,
			Longitude: 24.978291441099014,
		},
	})
	eventId2 := seedEvent2(t, ctx, queries, userId1, models.AddEvent{
		CalendarId: calendarId2,
		Name:       "daily meeting",
		StartTime:  startTime,
		EndTime:    endTime,
		Location: models.LocationInput{
			Name:      "toimisto",
			Address:   "toimistokatu 1",
			Latitude:  60.248947411912596,
			Longitude: 24.978291441099014,
		},
	})

	userId2 := seedUser(t, ctx, queries, "user2", "password2")
	calendarId3 := seedCalendar(t, ctx, queries, "project meetings", userId2)
	eventId3 := seedEvent2(t, ctx, queries, userId2, models.AddEvent{
		CalendarId: calendarId3,
		Name:       "weekly standup",
		StartTime:  startTime,
		EndTime:    endTime,
		Location: models.LocationInput{
			Name:      "kirjasto",
			Address:   "helsingin kirjasto 1",
			Latitude:  62.248947411912596,
			Longitude: 25.978291441099014,
		},
	})

	tests := []struct {
		name             string
		userId           uuid.UUID
		searchName       string
		expectedStatus   int
		expectedRespData []models.Event
	}{
		{
			name:           "fetching with keyword finds all (2) results",
			userId:         userId1,
			searchName:     "meeting",
			expectedStatus: http.StatusOK,
			expectedRespData: []models.Event{
				{
					Id:         eventId1,
					Name:       "team meeting",
					CalendarId: calendarId1,
					StartTime:  startTime,
					EndTime:    endTime,
					Location: &models.Location{
						Name:      "toimisto",
						Address:   "toimistokatu 1",
						Latitude:  60.248947411912596,
						Longitude: 24.978291441099014,
					},
				},
				{
					Id:         eventId2,
					Name:       "daily meeting",
					CalendarId: calendarId2,
					StartTime:  startTime,
					EndTime:    endTime,
					Location: &models.Location{
						Name:      "toimisto",
						Address:   "toimistokatu 1",
						Latitude:  60.248947411912596,
						Longitude: 24.978291441099014,
					},
				},
			},
		},
		{
			name:           "second user cannot search other user events",
			userId:         userId2,
			searchName:     "weekly",
			expectedStatus: http.StatusOK,
			expectedRespData: []models.Event{
				{
					Id:         eventId3,
					Name:       "weekly standup",
					CalendarId: calendarId3,
					StartTime:  startTime,
					EndTime:    endTime,
					Location: &models.Location{
						Name:      "kirjasto",
						Address:   "helsingin kirjasto 1",
						Latitude:  62.248947411912596,
						Longitude: 25.978291441099014,
					},
				},
			},
		},
		{
			name:             "no events are returned when search doesn't yield any results",
			userId:           userId1,
			searchName:       "abcxyz",
			expectedStatus:   http.StatusOK,
			expectedRespData: []models.Event{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, rec := echotest.ContextConfig{
				QueryValues: url.Values{
					"name": {tc.searchName},
				},
				Headers: map[string][]string{
					echo.HeaderContentType: {echo.MIMEApplicationJSON},
				},
			}.ToContextRecorder(t)
			c.Echo().Validator = &utils.CustomValidator{
				Validator: validator.New(validator.WithRequiredStructEnabled()),
			}
			c.Set("userId", tc.userId)

			_ = server.SearchEvents(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var got []models.Event
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)

			for i := range tc.expectedRespData {
				tc.expectedRespData[i].Id = got[i].Id
			}

			assert.ElementsMatch(t, tc.expectedRespData, got)
		})
	}
}

func TestAddEvent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "addEvents")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	server, queries := setupTestServer(t, ctx, connURI)

	userId := seedUser(t, ctx, queries, "eventUser", "password")
	calendarId := seedCalendar(t, ctx, queries, "meetings", userId)

	startTime := time.Now().UTC().Truncate(time.Microsecond)
	endTime := time.Now().UTC().Add(time.Hour).Truncate(time.Microsecond)

	randomUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	tests := []struct {
		name             string
		body             models.AddEvent
		expectedStatus   int
		expectedRespData models.Event
	}{
		{
			name: "adding event works",
			body: models.AddEvent{
				CalendarId: calendarId,
				Name:       "team meeting",
				StartTime:  startTime,
				EndTime:    endTime,
				Location: models.LocationInput{
					Name:      "kirjasto",
					Address:   "helsingin kirjasto 1",
					Latitude:  62.248947411912596,
					Longitude: 25.978291441099014,
				},
			},
			expectedStatus: http.StatusOK,
			expectedRespData: models.Event{
				// id is unknown beforehand
				Name:       "team meeting",
				CalendarId: calendarId,
				StartTime:  startTime,
				EndTime:    endTime,
				Location: &models.Location{
					Name:      "kirjasto",
					Address:   "helsingin kirjasto 1",
					Latitude:  62.248947411912596,
					Longitude: 25.978291441099014,
				},
			},
		},
		{
			name: "adding event to non-existent calendar fails",
			body: models.AddEvent{
				CalendarId: randomUUID,
				Name:       "team meeting",
				StartTime:  startTime,
				EndTime:    endTime,
				Location: models.LocationInput{
					Name:      "kirjasto",
					Address:   "helsingin kirjasto 1",
					Latitude:  62.248947411912596,
					Longitude: 25.978291441099014,
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "event with endTime before startTime fails",
			body: models.AddEvent{
				CalendarId: calendarId,
				Name:       "team meeting",
				StartTime:  endTime,
				EndTime:    startTime,
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

			_ = server.AddEvent(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var got models.Event
			err = json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)

			assert.NotNil(t, got.Id)
			assert.Equal(t, tc.expectedRespData.Name, got.Name)
			assert.Equal(t, tc.expectedRespData.CalendarId, got.CalendarId)
			assert.Equal(t, tc.expectedRespData.StartTime, got.StartTime)
			assert.Equal(t, tc.expectedRespData.EndTime, got.EndTime)
		})
	}
}

func TestEditEvent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "editEvents")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	server, queries := setupTestServer(t, ctx, connURI)

	startTime := time.Now().UTC().Truncate(time.Microsecond)
	endTime := time.Now().UTC().Add(time.Hour).Truncate(time.Microsecond)

	editedStartTime := time.Now().UTC().Add(time.Hour * 24).Truncate(time.Microsecond)
	editedEndTime := time.Now().UTC().Add(time.Hour * 48).Truncate(time.Microsecond)

	userId := seedUser(t, ctx, queries, "editEventUser", "password")
	calendarId := seedCalendar(t, ctx, queries, "meetings", userId)
	eventId := seedEvent(t, ctx, queries, "team meeting", userId, calendarId, startTime, endTime)

	userId2 := seedUser(t, ctx, queries, "editEventUser2", "password")
	calendarId2 := seedCalendar(t, ctx, queries, "meetings", userId2)

	editCalendarId := seedCalendar(t, ctx, queries, "meetings2", userId)

	randomUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	tests := []struct {
		name             string
		eventId          uuid.UUID
		body             models.EventEdit
		expectedStatus   int
		expectedRespData models.Event
	}{
		{
			name:    "editing events works",
			eventId: eventId,
			body: models.EventEdit{
				CalendarId: Ptr(calendarId),
				Name:       Ptr("daily"),
				StartTime:  Ptr(editedStartTime),
				EndTime:    Ptr(editedEndTime),
			},
			expectedStatus: http.StatusOK,
			expectedRespData: models.Event{
				Id:         eventId,
				CalendarId: calendarId,
				Name:       "daily",
				StartTime:  editedStartTime,
				EndTime:    editedEndTime,
			},
		},
		{
			name:    "calendar can be changed",
			eventId: eventId,
			body: models.EventEdit{
				CalendarId: Ptr(editCalendarId),
				Name:       Ptr("daily"),
			},
			expectedStatus: http.StatusOK,
			expectedRespData: models.Event{
				Id:         eventId,
				CalendarId: editCalendarId,
				Name:       "daily",
				StartTime:  editedStartTime,
				EndTime:    editedEndTime,
			},
		},
		{
			name:    "calendar can't be added to non-existent calendar",
			eventId: eventId,
			body: models.EventEdit{
				CalendarId: Ptr(randomUUID),
				Name:       Ptr("daily"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:    "calendar can't be added to other users calendar",
			eventId: eventId,
			body: models.EventEdit{
				CalendarId: Ptr(calendarId2),
				Name:       Ptr("daily"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()

			bodyJSON, err := json.Marshal(tc.body)
			require.NoError(t, err)

			c, rec := echotest.ContextConfig{
				PathValues: echo.PathValues{
					{Name: "eventId", Value: tc.eventId.String()},
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

			_ = server.EditEvent(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var got models.Event
			err = json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedRespData, got)
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "calendar3")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	server, queries := setupTestServer(t, ctx, connURI)

	startTime := time.Now().UTC().Truncate(time.Microsecond)
	endTime := time.Now().UTC().Add(time.Hour).Truncate(time.Microsecond)

	userId := seedUser(t, ctx, queries, "deleteCalendarUser", "password")
	calendarId := seedCalendar(t, ctx, queries, "meetings", userId)
	eventId := seedEvent(t, ctx, queries, "team meeting", userId, calendarId, startTime, endTime)

	randomUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	tests := []struct {
		name           string
		eventId        uuid.UUID
		expectedStatus int
	}{
		{
			name:           "deleting event works",
			eventId:        eventId,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "deleting event that does not exist does not fail",
			eventId:        randomUUID,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, rec := echotest.ContextConfig{
				PathValues: echo.PathValues{
					{Name: "eventId", Value: eventId.String()},
				},
				Headers: map[string][]string{
					echo.HeaderContentType: {echo.MIMEApplicationJSON},
				},
			}.ToContextRecorder(t)
			c.Set("userId", userId)

			_ = server.DeleteEvent(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
