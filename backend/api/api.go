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
	queryResp, err := s.queries.AddEvent(ctx, sqlc.AddEventParams{
		Name:        body.Name,
		Tstzrange:   body.StartTime,
		Tstzrange_2: body.EndTime,
	})

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := Event{
		Id:        queryResp.ID,
		Name:      queryResp.Name,
		StartTime: queryResp.Time.Lower.Time,
		EndTime:   queryResp.Time.Upper.Time,
	}

	return c.JSON(http.StatusOK, resp)
}

// TODO: probably unsafe long-term and should be deleted
// (GET /allEvents)
func (s *Server) GetAllEvents(c echo.Context) error {
	ctx := c.Request().Context()

	queryResp, err := s.queries.AllEvents(ctx)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := make([]Event, len(queryResp))
	for i, event := range queryResp {
		resp[i] = Event{
			Id:        event.ID,
			Name:      event.Name,
			StartTime: event.Time.Lower.Time,
			EndTime:   event.Time.Upper.Time,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// (GET /getEvents)
func (s *Server) GetGetEvents(c echo.Context, params GetGetEventsParams) error {
	ctx := c.Request().Context()

	fmt.Println(params.StartTime, params.EndTime)

	queryResp, err := s.queries.GetEvents(ctx, sqlc.GetEventsParams{
		Tstzrange:   params.StartTime,
		Tstzrange_2: params.EndTime,
	})

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := make([]Event, len(queryResp))
	for i, event := range queryResp {
		resp[i] = Event{
			Id:        event.ID,
			Name:      event.Name,
			StartTime: event.Time.Lower.Time,
			EndTime:   event.Time.Upper.Time,
		}
	}

	return c.JSON(http.StatusOK, resp)
}
