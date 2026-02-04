package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/eeritvan/calendar/internal/models"
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// (GET /getEvents?startTime=<END_TIME>&endTime=<START_TIME>)
func (s *Server) GetEvents(c *echo.Context) error {
	params := new(models.GetEventsParams)
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(params); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.GetEvents(ctx, sqlc.GetEventsParams{
		OwnerID:   userId,
		StartTime: params.StartTime,
		EndTime:   params.EndTime,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := make([]models.Event, len(queryResp))
	for i, event := range queryResp {
		resp[i] = models.Event{
			Id:         event.ID,
			CalendarId: event.CalendarID,
			Name:       event.Name,
			StartTime:  event.Time.Lower.Time.UTC(),
			EndTime:    event.Time.Upper.Time.UTC(),
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// (GET /searchEvents?name=<NAME>)
func (s *Server) SearchEvents(c *echo.Context) error {
	params := new(models.SearchEventsParams)
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(params); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.SearchEvents(ctx, sqlc.SearchEventsParams{
		OwnerID: userId,
		Name:    params.Name,
	})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := make([]models.Event, len(queryResp))
	for i, event := range queryResp {
		resp[i] = models.Event{
			Id:         event.ID,
			CalendarId: event.CalendarID,
			Name:       event.Name,
			StartTime:  event.Time.Lower.Time.UTC(),
			EndTime:    event.Time.Upper.Time.UTC(),
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// (POST /addEvent)
func (s *Server) AddEvent(c *echo.Context) error {
	body := new(models.AddEvent)
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.AddEvent(ctx, sqlc.AddEventParams{
		CalendarID: body.CalendarId,
		Name:       body.Name,
		OwnerID:    userId,
		StartTime:  body.StartTime,
		EndTime:    body.EndTime,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := models.Event{
		Id:         queryResp.ID,
		CalendarId: queryResp.CalendarID,
		Name:       queryResp.Name,
		StartTime:  queryResp.Time.Lower.Time.UTC(),
		EndTime:    queryResp.Time.Upper.Time.UTC(),
	}

	s.sse.Emit(userId, "event/post", resp)
	return c.JSON(http.StatusOK, resp)
}

// (PATCH /event/edit/:eventId)
// TODO: this crashes if the any field is missing (CalendarID and Name).
func (s *Server) EditEvent(c *echo.Context) error {
	eventId, err := echo.PathParam[uuid.UUID](c, "eventId")
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	body := new(models.EventEdit)
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	editedEvent, err := s.queries.EditEvent(ctx, sqlc.EditEventParams{
		ID:         eventId,
		OwnerID:    userId,
		CalendarID: *body.CalendarId,
		Name:       *body.Name,
		StartTime:  body.StartTime,
		EndTime:    body.EndTime,
	})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := models.Event{
		Id:         editedEvent.ID,
		CalendarId: editedEvent.CalendarID,
		Name:       editedEvent.Name,
		StartTime:  editedEvent.Time.Lower.Time.UTC(),
		EndTime:    editedEvent.Time.Upper.Time.UTC(),
	}

	s.sse.Emit(userId, "event/edit", resp)
	return c.JSON(http.StatusOK, resp)
}

// (DELETE /event/delete/:eventId)
func (s *Server) DeleteEvent(c *echo.Context) error {
	eventId, err := echo.PathParam[uuid.UUID](c, "eventId")
	if err != nil {
		return c.JSON(http.StatusBadRequest, false)
	}
	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	if err := s.queries.DeleteEvent(ctx, sqlc.DeleteEventParams{
		ID:      eventId,
		OwnerID: userId,
	}); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	s.sse.Emit(userId, "event/delete", eventId)
	return c.JSON(http.StatusOK, nil)
}

func parseDate(dateStr string, layouts []string) (time.Time, error) {
	for _, layout := range layouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return t, nil
		}
		fmt.Printf("errr parsing time %q as %q: %v\n", dateStr, layout, err)
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %s", dateStr)
}

// (POST /event/import)
func (s *Server) ImportEvents(c *echo.Context) error {
	layouts := []string{
		"20060102T150405Z", // 20201217T111500Z
		"20060102T150405",  // 20260206T140000
		"20060102T1504",    // 20240118T0930
		"20060102",         // 20210920
		"2006-01-02 15:04:05",
	}

	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if !strings.HasPrefix(contentType, "text/calendar") {
		return c.JSON(http.StatusUnsupportedMediaType, nil)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	icsContent := string(body)

	cal, err := ics.ParseCalendar(strings.NewReader(icsContent))
	ruuid, _ := uuid.NewRandom()

	events := make([]models.AddEvent, len(cal.Events()))
	for i, ev := range cal.Events() {
		startStr := ev.GetProperty(ics.ComponentPropertyDtStart).Value
		parsedStart, err := parseDate(startStr, layouts)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, nil)
		}

		endStr := ev.GetProperty(ics.ComponentPropertyDtEnd).Value
		parsedEnd, err := parseDate(endStr, layouts)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, nil)
		}

		event := models.AddEvent{
			CalendarId: ruuid,
			Name:       ev.GetProperty(ics.ComponentPropertySummary).Value,
			StartTime:  parsedStart,
			EndTime:    parsedEnd,
		}

		events[i] = event
	}

	fmt.Println(events)

	return c.JSON(http.StatusOK, nil)
}
