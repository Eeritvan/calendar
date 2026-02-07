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
	params := new(models.GetEventsParams)
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(params); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
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
		var location *models.Location
		if event.LocationID.Valid {
			location = &models.Location{
				Name:      event.LocationName,
				Address:   event.Address,
				Latitude:  event.Point.P.Y,
				Longitude: event.Point.P.X,
			}
		}

		resp[i] = models.Event{
			Id:         event.ID,
			CalendarId: event.CalendarID,
			Name:       event.Name,
			StartTime:  event.Time.Lower.Time.UTC(),
			EndTime:    event.Time.Upper.Time.UTC(),
			Location:   location,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// (GET /searchEvents?name=<NAME>)
func (s *Server) SearchEvents(c *echo.Context) error {
	params := new(models.SearchEventsParams)
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(params); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
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
		var location *models.Location
		if event.LocationID.Valid {
			location = &models.Location{
				Name:      event.LocationName,
				Address:   event.Address,
				Latitude:  event.Point.P.Y,
				Longitude: event.Point.P.X,
			}
		}

		resp[i] = models.Event{
			Id:         event.ID,
			CalendarId: event.CalendarID,
			Name:       event.Name,
			StartTime:  event.Time.Lower.Time.UTC(),
			EndTime:    event.Time.Upper.Time.UTC(),
			Location:   location,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// (POST /addEvent)
func (s *Server) AddEvent(c *echo.Context) error {
	body := new(models.AddEvent)
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.AddEvent(ctx, sqlc.AddEventParams{
		CalendarID:   body.CalendarId,
		Name:         body.Name,
		OwnerID:      userId,
		StartTime:    body.StartTime,
		EndTime:      body.EndTime,
		LocationName: body.Location.Name,
		Address:      body.Location.Address,
		Latitude:     body.Location.Latitude,
		Longitude:    body.Location.Longitude,
	})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	var location *models.Location
	if queryResp.LocationID.Valid {
		location = &models.Location{
			Name:      queryResp.LocationName,
			Address:   queryResp.Address,
			Latitude:  queryResp.Point.P.Y,
			Longitude: queryResp.Point.P.X,
		}
	}

	resp := models.Event{
		Id:         queryResp.ID,
		CalendarId: queryResp.CalendarID,
		Name:       queryResp.Name,
		StartTime:  queryResp.Time.Lower.Time.UTC(),
		EndTime:    queryResp.Time.Upper.Time.UTC(),
		Location:   location,
	}

	s.sse.Emit(userId, "event/post", resp)
	return c.JSON(http.StatusOK, resp)
}

// (PATCH /event/edit/:eventId)
// TODO: this crashes if the any field is missing (CalendarID and Name).
func (s *Server) EditEvent(c *echo.Context) error {
	eventId, err := echo.PathParam[uuid.UUID](c, "eventId")
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	body := new(models.EventEdit)
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
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
		StartTime:  editedEvent.Time.Lower.Time.UTC(),
		EndTime:    editedEvent.Time.Upper.Time.UTC(),
	}

	s.sse.Emit(userId, "event/edit", resp)
	return c.JSON(http.StatusOK, resp)
}

// (DELETE /event/delete/:eventId)
func (s *Server) DeleteEvent(c *echo.Context) error {
	eventId, err := echo.PathParam[uuid.UUID](c, "eventId")
	if err != nil {
		return c.JSON(http.StatusBadRequest, false)
	}
	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	if err := s.queries.DeleteEvent(ctx, sqlc.DeleteEventParams{
		ID:      eventId,
		OwnerID: userId,
	}); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	s.sse.Emit(userId, "event/delete", eventId)
	return c.JSON(http.StatusOK, nil)
}

// (POST /event/delete/batch)
func (s *Server) BatchDeleteEvents(c *echo.Context) error {
	body := new(models.BatchDeleteEvents)
	if err := c.Bind(&body); err != nil {
		fmt.Println("1", err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		fmt.Println("2", err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	batchParams := make([]sqlc.DeleteManyEventsParams, len(body.Ids))
	for i, id := range body.Ids {
		batchParams[i] = sqlc.DeleteManyEventsParams{
			ID:      id,
			OwnerID: userId,
		}
	}

	ctx := c.Request().Context()
	batchResults := s.queries.DeleteManyEvents(ctx, batchParams)

	var batchErr error
	batchResults.Exec(func(i int, err error) {
		if err != nil {
			fmt.Println(i, err)
			batchErr = err
		}
	})

	if batchErr != nil {
		fmt.Println(batchErr)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}
