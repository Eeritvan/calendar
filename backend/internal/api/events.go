package api

import (
	"fmt"
	"net/http"

	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// (GET /getEvents)
func (s *Server) GetGetEvents(c echo.Context, params GetGetEventsParams) error {
	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.GetEvents(ctx, sqlc.GetEventsParams{
		OwnerID:   userId,
		StartTime: params.StartTime,
		EndTime:   params.EndTime,
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

// (GET /searchEvents)
func (s *Server) GetSearchEvents(c echo.Context, params GetSearchEventsParams) error {
	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.SearchEvents(ctx, sqlc.SearchEventsParams{
		OwnerID: userId,
		Name:    params.Name,
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

// (POST /addEvent)
func (s *Server) PostAddEvent(c echo.Context) error {
	body := new(AddEvent)
	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.AddEvent(ctx, sqlc.AddEventParams{
		CalendarID: body.CalendarId,
		Name:       body.Name,
		OwnerID:    userId,
		StartTime:  body.StartTime,
		EndTime:    body.EndTime,
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

	s.emit(userId, "event/post", resp)
	return c.JSON(http.StatusOK, resp)
}

// (PATCH /event/edit/{event_id})
// TODO: this crashes if the any field is missing (CalendarID and Name).
func (s *Server) PatchEventEditEventId(c echo.Context, eventId uuid.UUID) error {
	body := new(EventEdit)
	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	editedEvent, err := s.queries.EditEvent(ctx, sqlc.EditEventParams{
		ID:         eventId,
		OwnerID:    userId,
		CalendarID: *body.CalendarId,
		Name:       *body.Name,
		StartTime:  body.StartTime,
		EndTime:    body.EndTime,
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

	s.emit(userId, "event/edit", resp)
	return c.JSON(http.StatusOK, resp)
}

// (DELETE /event/delete/{event_id})
func (s *Server) DeleteEventDeleteEventId(c echo.Context, eventId uuid.UUID) error {
	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	if err := s.queries.DeleteEvent(ctx, sqlc.DeleteEventParams{
		ID:      eventId,
		OwnerID: userId,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, false)
	}

	s.emit(userId, "event/delete", eventId)
	return c.JSON(http.StatusOK, true)
}
