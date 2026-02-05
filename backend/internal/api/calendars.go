package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	ics "github.com/arran4/golang-ical"
	"github.com/eeritvan/calendar/internal/models"
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// (GET /getCalendars)
func (s *Server) GetCalendars(c *echo.Context) error {
	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.GetCalendars(ctx, userId)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := make([]models.Calendar, len(queryResp))
	for i, calendar := range queryResp {
		resp[i] = models.Calendar{
			Id:      calendar.ID,
			Name:    calendar.Name,
			OwnerId: calendar.OwnerID,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// (POST /addCalendar)
func (s *Server) AddCalendar(c *echo.Context) error {
	body := new(models.AddCalendar)
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.AddCalendar(ctx, sqlc.AddCalendarParams{
		Name:    body.Name,
		OwnerID: userId,
	})

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := models.Calendar{
		Id:      queryResp.ID,
		Name:    queryResp.Name,
		OwnerId: queryResp.OwnerID,
	}

	s.sse.Emit(userId, "calendar/post", resp)
	return c.JSON(http.StatusOK, resp)
}

// (PATCH /calendar/edit/:calendarId)
// TODO: this crashes if the any field is missing (Name).
func (s *Server) EditCalendar(c *echo.Context) error {
	calendarId, _ := echo.PathParam[uuid.UUID](c, "calendarId")
	body := new(models.EditCalendar)
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	editedCalendar, err := s.queries.EditCalendar(ctx, sqlc.EditCalendarParams{
		Name:    *body.Name,
		ID:      calendarId,
		OwnerID: userId,
	})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := models.Calendar{
		Id:      editedCalendar.ID,
		Name:    editedCalendar.Name,
		OwnerId: editedCalendar.OwnerID,
	}

	s.sse.Emit(userId, "calendar/edit", resp)
	return c.JSON(http.StatusOK, resp)
}

// (DELETE /calendar/delete/:calendarId)

func (s *Server) DeleteCalendar(c *echo.Context) error {
	calendarId, err := echo.PathParam[uuid.UUID](c, "calendarId")
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	if err := s.queries.DeleteCalendar(ctx, sqlc.DeleteCalendarParams{
		ID:      calendarId,
		OwnerID: userId,
	}); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	s.sse.Emit(userId, "calendar/delete", calendarId)
	return c.JSON(http.StatusOK, nil)
}

// (POST /calendar/:calendarId/event/import)
func (s *Server) ImportEvents(c *echo.Context) error {
	userId := c.Get("userId").(uuid.UUID)

	layouts := []string{
		"20060102T150405Z",
		"20060102T150405",
		"20060102T1504",
		"20060102",
		"2006-01-02 15:04:05",
	}

	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if !strings.HasPrefix(contentType, "text/calendar") {
		return c.JSON(http.StatusUnsupportedMediaType, nil)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	icsContent := string(body)

	cal, err := ics.ParseCalendar(strings.NewReader(icsContent))

	calendarId, err := echo.PathParam[uuid.UUID](c, "calendarId")
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	resp := make([]models.Event, len(cal.Events()))
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

		ctx := c.Request().Context()
		queryResp, err := s.queries.AddEvent(ctx, sqlc.AddEventParams{
			CalendarID: calendarId,
			OwnerID:    userId,
			Name:       ev.GetProperty(ics.ComponentPropertySummary).Value,
			StartTime:  parsedStart,
			EndTime:    parsedEnd,
		})

		event := models.Event{
			Id:         queryResp.ID,
			CalendarId: queryResp.CalendarID,
			Name:       queryResp.Name,
			StartTime:  queryResp.Time.Lower.Time.UTC(),
			EndTime:    queryResp.Time.Upper.Time.UTC(),
		}

		resp[i] = event
	}

	return c.JSON(http.StatusOK, resp)
}
