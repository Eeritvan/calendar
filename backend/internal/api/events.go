package api

import (
	"fmt"
	"net/http"

	"github.com/eeritvan/calendar/internal/models"
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// (GET /getEvents?startTime=<END_TIME>&endTime=<START_TIME>)
func (s *Server) GetEvents(c *echo.Context) error {
	params := new(models.GetGetEventsParams)
	if err := c.Bind(params); err != nil {
		// TODO: error handling
		return nil
	}

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

	resp := make([]models.Event, len(queryResp))
	for i, event := range queryResp {
		resp[i] = models.Event{
			Id:         event.ID,
			CalendarId: event.CalendarID,
			Name:       event.Name,
			StartTime:  event.Time.Lower.Time,
			EndTime:    event.Time.Upper.Time,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// (GET /searchEvents?name=<NAME>)
func (s *Server) SearchEvents(c *echo.Context) error {
	params := new(models.GetSearchEventsParams)
	if err := c.Bind(params); err != nil {
		// TODO: error handling
		return nil
	}

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

	resp := make([]models.Event, len(queryResp))
	for i, event := range queryResp {
		resp[i] = models.Event{
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
func (s *Server) AddEvent(c *echo.Context) error {
	body := new(models.AddEvent)
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

	resp := models.Event{
		Id:         queryResp.ID,
		CalendarId: queryResp.CalendarID,
		Name:       queryResp.Name,
		StartTime:  queryResp.Time.Lower.Time,
		EndTime:    queryResp.Time.Upper.Time,
	}

	s.emit(userId, "event/post", resp)
	return c.JSON(http.StatusOK, resp)
}

// (PATCH /event/edit/:eventId)
// TODO: this crashes if the any field is missing (CalendarID and Name).
func (s *Server) EventEdit(c *echo.Context) error {
	eventId, _ := echo.PathParam[uuid.UUID](c, "eventID")
	body := new(models.EventEdit)
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

	resp := models.Event{
		Id:         editedEvent.ID,
		CalendarId: editedEvent.CalendarID,
		Name:       editedEvent.Name,
		StartTime:  editedEvent.Time.Lower.Time,
		EndTime:    editedEvent.Time.Upper.Time,
	}

	s.emit(userId, "event/edit", resp)
	return c.JSON(http.StatusOK, resp)
}

// (DELETE /event/delete/:eventId)
func (s *Server) EventDelete(c *echo.Context) error {
	eventId, _ := echo.PathParam[uuid.UUID](c, "eventID")
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
