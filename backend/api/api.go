package api

import (
	"fmt"
	"net/http"

	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/labstack/echo/v4"
)

type Server struct {
	queries *sqlc.Queries
}

func NewServer(queries *sqlc.Queries) *Server {
	return &Server{queries: queries}
}

// (POST /addEvent)
func (s *Server) PostAddEvent(c echo.Context) error {
	body := new(EventNoId)

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	ctx := c.Request().Context()
	queryResp, err := s.queries.AddEvent(ctx, body.Name)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := Event{
		Id:   int(queryResp.ID),
		Name: queryResp.Name,
	}

	return c.JSON(http.StatusOK, resp)
}

// (get /allEvents)
func (s *Server) GetAllEvents(c echo.Context) error {
	ctx := c.Request().Context()

	queryResp, err := s.queries.GetEvents(ctx)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := make([]Event, len(queryResp))
	for i, event := range queryResp {
		resp[i] = Event{
			Id:   int(event.ID),
			Name: event.Name,
		}
	}

	return c.JSON(http.StatusOK, resp)
}
