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
	eventId1 := seedEvent(t, ctx, queries, userId1, models.AddEvent{
		CalendarId: calendarId1,
		Name:       "team meeting",
		StartTime:  timeMinusHour,
		EndTime:    timeNow,
	})
	eventId2 := seedEvent(t, ctx, queries, userId1, models.AddEvent{
		CalendarId: calendarId2,
		Name:       "monday standup",
		StartTime:  timeNow,
		EndTime:    timePlusHour,
		Location: &models.LocationInput{
			Name:      "toimisto",
			Address:   utils.Ptr("toimistokatu 1"),
			Latitude:  utils.Ptr(60.248947411912596),
			Longitude: utils.Ptr(24.978291441099014),
		},
	})
	userId2 := seedUser(t, ctx, queries, "user2", "password2")
	calendarId3 := seedCalendar(t, ctx, queries, "video games", userId2)
	eventId3 := seedEvent(t, ctx, queries, userId2, models.AddEvent{
		CalendarId: calendarId3,
		Name:       "my winter car",
		StartTime:  timeMinusHour,
		EndTime:    timePlusHour,
	})
	userId3 := seedUser(t, ctx, queries, "user3", "password3")
	calendarId4 := seedCalendar(t, ctx, queries, "testing calendar", userId3)
	eventId4 := seedEvent(t, ctx, queries, userId3, models.AddEvent{
		CalendarId: calendarId4,
		Name:       "testing event",
		StartTime:  timeMinusHour,
		EndTime:    timePlusHour,
		Location: &models.LocationInput{
			Name:    "testing",
			Address: utils.Ptr("testing"),
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
				},
				{
					Id:         eventId2,
					Name:       "monday standup",
					CalendarId: calendarId2,
					StartTime:  timeNow,
					EndTime:    timePlusHour,
					Location: &models.Location{
						Name:      "toimisto",
						Address:   utils.Ptr("toimistokatu 1"),
						Latitude:  utils.Ptr(60.248947411912596),
						Longitude: utils.Ptr(24.978291441099014),
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
				},
			},
		},
		{
			name:           "event with missing lat/lng works",
			userId:         userId3,
			queryStartTime: timeMinusHour.Format(time.RFC3339),
			queryEndTime:   timeNow.Format(time.RFC3339),
			expectedStatus: http.StatusOK,
			expectedRespData: []models.Event{
				{
					Id:         eventId4,
					Name:       "testing event",
					CalendarId: calendarId4,
					StartTime:  timeMinusHour,
					EndTime:    timePlusHour,
					Location: &models.Location{
						Name:    "testing",
						Address: utils.Ptr("testing"),
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
	eventId1 := seedEvent(t, ctx, queries, userId1, models.AddEvent{
		CalendarId: calendarId1,
		Name:       "team meeting",
		StartTime:  startTime,
		EndTime:    endTime,
	})
	eventId2 := seedEvent(t, ctx, queries, userId1, models.AddEvent{
		CalendarId: calendarId2,
		Name:       "daily meeting",
		StartTime:  startTime,
		EndTime:    endTime,
		Location: &models.LocationInput{
			Name:      "toimisto",
			Address:   utils.Ptr("toimistokatu 1"),
			Latitude:  utils.Ptr(60.248947411912596),
			Longitude: utils.Ptr(24.978291441099014),
		},
	})

	userId2 := seedUser(t, ctx, queries, "user2", "password2")
	calendarId3 := seedCalendar(t, ctx, queries, "project meetings", userId2)
	eventId3 := seedEvent(t, ctx, queries, userId2, models.AddEvent{
		CalendarId: calendarId3,
		Name:       "weekly standup",
		StartTime:  startTime,
		EndTime:    endTime,
	})

	userId3 := seedUser(t, ctx, queries, "user3", "password3")
	calendarId4 := seedCalendar(t, ctx, queries, "testing calendar", userId3)
	eventId4 := seedEvent(t, ctx, queries, userId3, models.AddEvent{
		CalendarId: calendarId4,
		Name:       "testing event",
		StartTime:  startTime,
		EndTime:    endTime,
		Location: &models.LocationInput{
			Name:    "testing",
			Address: utils.Ptr("testing"),
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
					Location:   nil,
				},
				{
					Id:         eventId2,
					Name:       "daily meeting",
					CalendarId: calendarId2,
					StartTime:  startTime,
					EndTime:    endTime,
					Location: &models.Location{
						Name:      "toimisto",
						Address:   utils.Ptr("toimistokatu 1"),
						Latitude:  utils.Ptr(60.248947411912596),
						Longitude: utils.Ptr(24.978291441099014),
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
				},
			},
		},
		{
			name:           "event with missing lat/lng works",
			userId:         userId3,
			searchName:     "testing",
			expectedStatus: http.StatusOK,
			expectedRespData: []models.Event{
				{
					Id:         eventId4,
					Name:       "testing event",
					CalendarId: calendarId4,
					StartTime:  startTime,
					EndTime:    endTime,
					Location: &models.Location{
						Name:    "testing",
						Address: utils.Ptr("testing"),
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
			name: "adding event with location works",
			body: models.AddEvent{
				CalendarId: calendarId,
				Name:       "team meeting",
				StartTime:  startTime,
				EndTime:    endTime,
				Location: &models.LocationInput{
					Name:      "kirjasto",
					Address:   utils.Ptr("helsingin kirjasto 1"),
					Latitude:  utils.Ptr(62.248947411912596),
					Longitude: utils.Ptr(25.978291441099014),
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
					Address:   utils.Ptr("helsingin kirjasto 1"),
					Latitude:  utils.Ptr(62.248947411912596),
					Longitude: utils.Ptr(25.978291441099014),
				},
			},
		},
		{
			name: "adding event without location works",
			body: models.AddEvent{
				CalendarId: calendarId,
				Name:       "online team meeting",
				StartTime:  startTime,
				EndTime:    endTime,
			},
			expectedStatus: http.StatusOK,
			expectedRespData: models.Event{
				// id is unknown beforehand
				Name:       "online team meeting",
				CalendarId: calendarId,
				StartTime:  startTime,
				EndTime:    endTime,
				Location:   nil,
			},
		},
		{
			name: "adding event with only location name works",
			body: models.AddEvent{
				CalendarId: calendarId,
				Name:       "online team meeting",
				StartTime:  startTime,
				EndTime:    endTime,
				Location: &models.LocationInput{
					Name: "home",
				},
			},
			expectedStatus: http.StatusOK,
			expectedRespData: models.Event{
				// id is unknown beforehand
				Name:       "online team meeting",
				CalendarId: calendarId,
				StartTime:  startTime,
				EndTime:    endTime,
				Location: &models.Location{
					Name: "home",
				},
			},
		},
		{
			name: "adding event with location name and lat/lng works",
			body: models.AddEvent{
				CalendarId: calendarId,
				Name:       "online team meeting",
				StartTime:  startTime,
				EndTime:    endTime,
				Location: &models.LocationInput{
					Name:      "home",
					Latitude:  utils.Ptr(33.4321),
					Longitude: utils.Ptr(22.1234),
				},
			},
			expectedStatus: http.StatusOK,
			expectedRespData: models.Event{
				// id is unknown beforehand
				Name:       "online team meeting",
				CalendarId: calendarId,
				StartTime:  startTime,
				EndTime:    endTime,
				Location: &models.Location{
					Name:      "home",
					Latitude:  utils.Ptr(33.4321),
					Longitude: utils.Ptr(22.1234),
				},
			},
		},
		{
			name: "adding event with location name and address works",
			body: models.AddEvent{
				CalendarId: calendarId,
				Name:       "online team meeting",
				StartTime:  startTime,
				EndTime:    endTime,
				Location: &models.LocationInput{
					Name:    "home",
					Address: utils.Ptr("testroad 1"),
				},
			},
			expectedStatus: http.StatusOK,
			expectedRespData: models.Event{
				// id is unknown beforehand
				Name:       "online team meeting",
				CalendarId: calendarId,
				StartTime:  startTime,
				EndTime:    endTime,
				Location: &models.Location{
					Name:    "home",
					Address: utils.Ptr("testroad 1"),
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
				Location: &models.LocationInput{
					Name:      "kirjasto",
					Address:   utils.Ptr("helsingin kirjasto 1"),
					Latitude:  utils.Ptr(62.248947411912596),
					Longitude: utils.Ptr(25.978291441099014),
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
		{
			name: "event with added latitude but missing longitude fails",
			body: models.AddEvent{
				CalendarId: calendarId,
				Name:       "team meeting",
				StartTime:  endTime,
				EndTime:    startTime,
				Location: &models.LocationInput{
					Name:     "office",
					Latitude: utils.Ptr(22.12345),
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "event with added longitude but missing latitude fails",
			body: models.AddEvent{
				CalendarId: calendarId,
				Name:       "team meeting",
				StartTime:  endTime,
				EndTime:    startTime,
				Location: &models.LocationInput{
					Name:      "office",
					Longitude: utils.Ptr(22.12345),
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "event with too high latitude fails",
			body: models.AddEvent{
				CalendarId: calendarId,
				Name:       "online team meeting",
				StartTime:  startTime,
				EndTime:    endTime,
				Location: &models.LocationInput{
					Name:      "office",
					Longitude: utils.Ptr(22.1234),
					Latitude:  utils.Ptr(181.0),
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "event with too low longitude fails",
			body: models.AddEvent{
				CalendarId: calendarId,
				Name:       "online team meeting",
				StartTime:  startTime,
				EndTime:    endTime,
				Location: &models.LocationInput{
					Name:      "office",
					Longitude: utils.Ptr(22.123),
					Latitude:  utils.Ptr(-91.0),
				},
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
			assert.Equal(t, tc.expectedRespData.Location, got.Location)
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
	eventId := seedEvent(t, ctx, queries, userId, models.AddEvent{
		CalendarId: calendarId,
		Name:       "team meeting",
		StartTime:  startTime,
		EndTime:    endTime,
	})
	eventId2 := seedEvent(t, ctx, queries, userId, models.AddEvent{
		CalendarId: calendarId,
		Name:       "team meeting",
		StartTime:  startTime,
		EndTime:    endTime,
		Location: &models.LocationInput{
			Name:      "office",
			Address:   utils.Ptr("officeRoad 1"),
			Latitude:  utils.Ptr(12.3456),
			Longitude: utils.Ptr(65.4321),
		},
	})

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
				CalendarId: utils.Ptr(calendarId),
				Name:       utils.Ptr("daily"),
				StartTime:  utils.Ptr(editedStartTime),
				EndTime:    utils.Ptr(editedEndTime),
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
				CalendarId: utils.Ptr(editCalendarId),
				Name:       utils.Ptr("daily"),
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
			name:    "location can be changed",
			eventId: eventId,
			body: models.EventEdit{
				CalendarId: utils.Ptr(calendarId),
				Name:       utils.Ptr("testing"),
				Location: &models.LocationEdit{
					Name:      utils.Ptr("office"),
					Address:   utils.Ptr("road 1"),
					Latitude:  utils.Ptr(33.33),
					Longitude: utils.Ptr(22.22),
				},
			},
			expectedStatus: http.StatusOK,
			expectedRespData: models.Event{
				Id:         eventId,
				CalendarId: calendarId,
				Name:       "testing",
				StartTime:  editedStartTime,
				EndTime:    editedEndTime,
				Location: &models.Location{
					Name:      "office",
					Address:   utils.Ptr("road 1"),
					Latitude:  utils.Ptr(33.33),
					Longitude: utils.Ptr(22.22),
				},
			},
		},
		{
			name:    "location can be removed",
			eventId: eventId2,
			body: models.EventEdit{
				CalendarId: utils.Ptr(calendarId),
				Name:       utils.Ptr("testing"),
				Location:   nil,
			},
			expectedStatus: http.StatusOK,
			expectedRespData: models.Event{
				Id:         eventId2,
				CalendarId: calendarId,
				Name:       "testing",
				StartTime:  startTime,
				EndTime:    endTime,
				Location:   nil,
			},
		},
		{
			name:    "event can't be added to non-existent calendar",
			eventId: eventId,
			body: models.EventEdit{
				CalendarId: utils.Ptr(randomUUID),
				Name:       utils.Ptr("daily"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:    "event latitude and longitude must be edited together",
			eventId: eventId,
			body: models.EventEdit{
				CalendarId: utils.Ptr(calendarId),
				Name:       utils.Ptr("daily"),
				Location: &models.LocationEdit{
					Name:     utils.Ptr("testing"),
					Latitude: utils.Ptr(23.22),
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "event can't be added to other users calendar",
			eventId: eventId,
			body: models.EventEdit{
				CalendarId: utils.Ptr(calendarId2),
				Name:       utils.Ptr("daily"),
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
	eventId := seedEvent(t, ctx, queries, userId, models.AddEvent{
		CalendarId: calendarId,
		Name:       "team meeting",
		StartTime:  startTime,
		EndTime:    endTime,
	})

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

func TestBatchDeleteEvents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	connURI, err := spawnPostgresContainer(t, "calendar5")
	require.NoError(t, err)

	err = runMigrations(t, connURI)
	require.NoError(t, err)

	server, queries := setupTestServer(t, ctx, connURI)

	startTime := time.Now().UTC().Truncate(time.Microsecond)
	endTime := time.Now().UTC().Add(time.Hour).Truncate(time.Microsecond)

	userId := seedUser(t, ctx, queries, "batchDeleteCalendarUser", "password")
	calendarId := seedCalendar(t, ctx, queries, "meetings", userId)
	eventId := seedEvent(t, ctx, queries, userId, models.AddEvent{
		CalendarId: calendarId,
		Name:       "team meeting",
		StartTime:  startTime,
		EndTime:    endTime,
	})
	eventId2 := seedEvent(t, ctx, queries, userId, models.AddEvent{
		CalendarId: calendarId,
		Name:       "team meeting2",
		StartTime:  startTime,
		EndTime:    endTime,
	})
	eventId3 := seedEvent(t, ctx, queries, userId, models.AddEvent{
		CalendarId: calendarId,
		Name:       "team meeting3",
		StartTime:  startTime,
		EndTime:    endTime,
	})

	randomUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	tests := []struct {
		name           string
		body           models.BatchDeleteEvents
		expectedStatus int
	}{
		{
			name: "deleting event works",
			body: models.BatchDeleteEvents{
				Ids: []uuid.UUID{
					eventId,
					eventId2,
					eventId3,
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "deleting events that does not exist does not fail",
			body: models.BatchDeleteEvents{
				Ids: []uuid.UUID{
					randomUUID,
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "request with no events will fail",
			body:           models.BatchDeleteEvents{},
			expectedStatus: http.StatusBadRequest,
		}}

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

			_ = server.BatchDeleteEvents(c)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
