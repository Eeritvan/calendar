package api

import (
	"fmt"
	"net/http"

	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// (GET /getCalendars)
func (s *Server) GetGetCalendars(c echo.Context) error {
	ctx := c.Request().Context()

	userIdStr := c.Get("userId").(string)
	userUUID, err := uuid.Parse(userIdStr)
	if err != nil {
		fmt.Printf("Invalid UUID format: %v\n", err)
	}

	queryResp, err := s.queries.GetCalendars(ctx, userUUID)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := make([]Calendar, len(queryResp))
	for i, calendar := range queryResp {
		resp[i] = Calendar{
			Id:      calendar.ID,
			Name:    calendar.Name,
			OwnerId: calendar.OwnerID,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// (POST /addCalendar)
func (s *Server) PostAddCalendar(c echo.Context) error {
	body := new(CalendarNoId)

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	userIdStr := c.Get("userId").(string)
	userUUID, err := uuid.Parse(userIdStr)
	if err != nil {
		fmt.Printf("Invalid UUID format: %v\n", err)
	}

	ctx := c.Request().Context()
	queryResp, err := s.queries.AddCalendar(ctx, sqlc.AddCalendarParams{
		Name:    body.Name,
		OwnerID: userUUID,
	})

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := Calendar{
		Id:      queryResp.ID,
		Name:    queryResp.Name,
		OwnerId: queryResp.OwnerID,
	}

	return c.JSON(http.StatusOK, resp)
}

// (PATCH /calendar/edit/{calendar_id})
// TODO: this crashes if the any field is missing.
// TODO: validate that only owner can edit calendar
func (s *Server) PatchCalendarEditCalendarId(c echo.Context, calendarId uuid.UUID) error {
	body := new(CalendarEdit)

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	ctx := c.Request().Context()
	editedCalendar, err := s.queries.EditCalendar(ctx, sqlc.EditCalendarParams{
		Name: *body.Name,
		ID:   calendarId,
	})

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := Calendar{
		Id:      editedCalendar.ID,
		Name:    editedCalendar.Name,
		OwnerId: editedCalendar.OwnerID,
	}

	return c.JSON(http.StatusOK, resp)
}

// (DELETE /calendar/delete/{calendar_id})
// TODO: validate that only owner can delete calendar
func (s *Server) DeleteCalendarDeleteCalendarId(c echo.Context, calendarId uuid.UUID) error {
	ctx := c.Request().Context()
	if err := s.queries.DeleteCalendar(ctx, calendarId); err != nil {
		return c.JSON(http.StatusInternalServerError, false)
	}

	return c.JSON(http.StatusOK, true)
}
