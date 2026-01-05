package api

import (
	"fmt"
	"net/http"

	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/google/uuid"
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
			Id:         event.ID,
			CalendarId: event.CalendarID,
			Name:       event.Name,
			StartTime:  event.Time.Lower.Time,
			EndTime:    event.Time.Upper.Time,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// (PATCH /event/edit/{event_id})
// TODO: this crashes if the any field is missing.
// TODO: validate that only owner can edit event
// TODO: change calendar
func (s *Server) PatchEventEditEventId(c echo.Context, eventId uuid.UUID) error {
	body := new(EventEdit)

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	fmt.Println(*body.EndTime)
	fmt.Println(*body.StartTime)

	ctx := c.Request().Context()
	editedEvent, err := s.queries.EditEvent(ctx, sqlc.EditEventParams{
		Name:    *body.Name,
		Column2: body.StartTime,
		Column3: body.EndTime,
		ID:      eventId,
	})

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := Event{
		Id:         editedEvent.ID,
		CalendarId: editedEvent.CalendarID,
		Name:       editedEvent.Name,
		StartTime:  editedEvent.Time.Lower.Time,
		EndTime:    editedEvent.Time.Upper.Time,
	}

	return c.JSON(http.StatusOK, resp)
}

// (DELETE /event/delete/{event_id})
// TODO: validate that only owner can delete event
func (s *Server) DeleteEventDeleteEventId(c echo.Context, eventId uuid.UUID) error {
	ctx := c.Request().Context()
	if err := s.queries.DeleteEvent(ctx, eventId); err != nil {
		return c.JSON(http.StatusInternalServerError, false)
	}

	return c.JSON(http.StatusOK, true)
}
