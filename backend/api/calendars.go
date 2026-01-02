package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// (GET /getCalendars)
func (s *Server) GetGetCalendars(c echo.Context) error {
	ctx := c.Request().Context()

	queryResp, err := s.queries.GetCalendars(ctx)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := make([]Calendar, len(queryResp))
	for i, event := range queryResp {
		resp[i] = Calendar{
			Id:   event.ID,
			Name: event.Name,
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

	ctx := c.Request().Context()
	queryResp, err := s.queries.AddCalendar(ctx, body.Name)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := Calendar{
		Id:   queryResp.ID,
		Name: queryResp.Name,
	}

	return c.JSON(http.StatusOK, resp)
}
