package api

import (
	"fmt"
	"net/http"

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

	s.emit(userId, "calendar/post", resp)
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

	s.emit(userId, "calendar/edit", resp)
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
		return c.JSON(http.StatusInternalServerError, false)
	}

	s.emit(userId, "calendar/delete", calendarId)
	return c.JSON(http.StatusOK, true)
}
