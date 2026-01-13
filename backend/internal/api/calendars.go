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
	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.GetCalendars(ctx, userId)
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

	resp := Calendar{
		Id:      queryResp.ID,
		Name:    queryResp.Name,
		OwnerId: queryResp.OwnerID,
	}

	s.emit(userId, "calendar/post", resp)
	return c.JSON(http.StatusOK, resp)
}

// (PATCH /calendar/edit/{calendar_id})
// TODO: this crashes if the any field is missing (Name).
func (s *Server) PatchCalendarEditCalendarId(c echo.Context, calendarId uuid.UUID) error {
	body := new(CalendarEdit)
	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
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

	resp := Calendar{
		Id:      editedCalendar.ID,
		Name:    editedCalendar.Name,
		OwnerId: editedCalendar.OwnerID,
	}

	s.emit(userId, "calendar/edit", resp)
	return c.JSON(http.StatusOK, resp)
}

// (DELETE /calendar/delete/{calendar_id})
func (s *Server) DeleteCalendarDeleteCalendarId(c echo.Context, calendarId uuid.UUID) error {
	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	if err := s.queries.DeleteCalendar(ctx, sqlc.DeleteCalendarParams{
		ID:      calendarId,
		OwnerID: userId,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, false)
	}

	s.emit(userId, "calendar/delete", calendarId)
	return c.JSON(http.StatusOK, true)
}
