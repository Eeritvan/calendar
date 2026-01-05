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
