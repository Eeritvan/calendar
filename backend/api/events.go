package api

import (
	"fmt"
	"net/http"

	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/labstack/echo/v4"
)

// (POST /addEvent)
// -- TODO: users can add events to calendars they dont own
func (s *Server) PostAddEvent(c echo.Context) error {
	body := new(AddEvent)

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	ctx := c.Request().Context()
	queryResp, err := s.queries.AddEvent(ctx, sqlc.AddEventParams{
		CalendarID:  body.CalendarId,
		Name:        body.Name,
		Tstzrange:   body.StartTime,
		Tstzrange_2: body.EndTime,
	})

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := Event{
		Id:         queryResp.ID,
		CalendarId: queryResp.CalendarID,
		Name:       queryResp.Name,
		StartTime:  queryResp.Time.Lower.Time,
		EndTime:    queryResp.Time.Upper.Time,
	}

	return c.JSON(http.StatusOK, resp)
}

// (GET /getEvents)
// TODO: get only own events
func (s *Server) GetGetEvents(c echo.Context, params GetGetEventsParams) error {
	ctx := c.Request().Context()

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
