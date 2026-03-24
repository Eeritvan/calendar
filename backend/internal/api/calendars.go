package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
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
		fmt.Println("test", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := make([]models.Calendar, len(queryResp))
	for i, calendar := range queryResp {
		item := models.Calendar{
			Id:         calendar.ID,
			Name:       calendar.Name,
			OwnerId:    calendar.OwnerID,
			Visibility: calendar.Visibility,
			Permission: calendar.Permission,
			IsOwner:    calendar.IsOwner,
			Folder:     nil,
		}

		if calendar.FolderID != uuid.Nil {
			item.Folder = &models.Folder{
				Id:   calendar.FolderID,
				Name: *calendar.FolderName,
			}
		}

		resp[i] = item
	}

	return c.JSON(http.StatusOK, resp)
}

// (POST /addCalendar)
// -- TODO: add to folder also
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
		Id:         queryResp.ID,
		Name:       queryResp.Name,
		OwnerId:    queryResp.OwnerID,
		Visibility: queryResp.Visibility,
	}

	s.sse.Emit(userId, "calendar/post", resp)
	return c.JSON(http.StatusCreated, resp)
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
		Id:         editedCalendar.ID,
		Name:       editedCalendar.Name,
		OwnerId:    editedCalendar.OwnerID,
		Visibility: editedCalendar.Visibility,
	}

	s.sse.Emit(userId, "calendar/edit", resp)
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
		return c.JSON(http.StatusInternalServerError, nil)
	}

	s.sse.Emit(userId, "calendar/delete", calendarId)
	return c.NoContent(http.StatusNoContent)
}

func parseDate(dateStr string, layouts []string) (time.Time, error) {
	for _, layout := range layouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return t, nil
		}
		fmt.Printf("errr parsing time %q as %q: %v\n", dateStr, layout, err)
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %s", dateStr)
}

// (POST /calendar/:calendarId/event/import)
func (s *Server) ImportEvents(c *echo.Context) error {
	userId := c.Get("userId").(uuid.UUID)

	layouts := []string{
		"20060102T150405",
		"20060102T150405Z",
		"20060102T1504",
		"20060102",
		"2006-01-02 15:04:05",
	}

	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if !strings.HasPrefix(contentType, "text/calendar") {
		return c.JSON(http.StatusUnsupportedMediaType, nil)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	icsContent := string(body)

	cal, err := ics.ParseCalendar(strings.NewReader(icsContent))
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	calendarId, err := echo.PathParam[uuid.UUID](c, "calendarId")
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	batchParams := make([]sqlc.ImportCalendarEventsParams, len(cal.Events()))
	for i, ev := range cal.Events() {
		startStr := ev.GetProperty(ics.ComponentPropertyDtStart).Value
		parsedStart, err := parseDate(startStr, layouts)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, nil)
		}

		endStr := ev.GetProperty(ics.ComponentPropertyDtEnd).Value
		parsedEnd, err := parseDate(endStr, layouts)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, nil)
		}

		batchParams[i] = sqlc.ImportCalendarEventsParams{
			CalendarID: calendarId,
			OwnerID:    userId,
			Name:       ev.GetProperty(ics.ComponentPropertySummary).Value,
			StartTime:  parsedStart,
			EndTime:    parsedEnd,
		}
	}

	ctx := c.Request().Context()
	batchResults := s.queries.ImportCalendarEvents(ctx, batchParams)

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

	return c.NoContent(http.StatusNoContent)
}

// (POST /calendar/:calendarId/event/export)
func (s *Server) ExportEvents(c *echo.Context) error {
	userId := c.Get("userId").(uuid.UUID)
	calendarId, err := echo.PathParam[uuid.UUID](c, "calendarId")
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	ctx := c.Request().Context()
	queryResp, err := s.queries.ExportCalendarEvents(ctx, sqlc.ExportCalendarEventsParams{
		CalendarID: calendarId,
		OwnerID:    userId,
	})

	cal := ics.NewCalendar()
	for _, event := range queryResp {
		iscEvent := cal.AddEvent(fmt.Sprintf("%s@eeritvan.dev", uuid.New().String()))

		iscEvent.SetSummary(event.Name)
		iscEvent.SetStartAt(event.Time.Lower.Time.UTC())
		iscEvent.SetEndAt(event.Time.Upper.Time.UTC())
	}

	icsData := []byte(cal.Serialize())

	c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="events.ics"`)
	return c.Blob(http.StatusOK, "text/calendar; charset=utf-8", icsData)
}

// (POST /calendar/:calendarId/share)
//
// TODO: set calendar status to shared
func (s *Server) ShareCalendar(c *echo.Context) error {
	userId := c.Get("userId").(uuid.UUID)
	calendarId, err := echo.PathParam[uuid.UUID](c, "calendarId")
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	body := new(models.ShareCalendar)
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	ctx := c.Request().Context()
	if err = s.queries.ShareCalendar(ctx, sqlc.ShareCalendarParams{
		CalendarID: calendarId,
		SharedWith: body.UserId,
		Permission: body.Permission,
		OwnerID:    userId,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}

// (POST /calendar/:calendarId/share/batch)
//
// TODO: set calendar status to shared
func (s *Server) BatchShareCalendar(c *echo.Context) error {
	userId := c.Get("userId").(uuid.UUID)
	calendarId, err := echo.PathParam[uuid.UUID](c, "calendarId")
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	body := new(models.BatchShareCalendar)
	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	batchParams := make([]sqlc.BatchShareCalendarParams, len(body.Items))
	for i, body := range body.Items {
		batchParams[i] = sqlc.BatchShareCalendarParams{
			CalendarID: calendarId,
			SharedWith: body.UserId,
			Permission: body.Permission,
			OwnerID:    userId,
		}
	}

	ctx := c.Request().Context()
	batchResults := s.queries.BatchShareCalendar(ctx, batchParams)

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

// (PATCH /calendar/:calendarId/share/private)
func (s *Server) ShareCalendarPrivate(c *echo.Context) error {
	userId := c.Get("userId").(uuid.UUID)
	calendarId, err := echo.PathParam[uuid.UUID](c, "calendarId")
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	ctx := c.Request().Context()
	if err = s.queries.SetCalendarPrivate(ctx, sqlc.SetCalendarPrivateParams{
		ID:      calendarId,
		OwnerID: userId,
	}); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if err := s.queries.DeleteAllCalendarShares(ctx, sqlc.DeleteAllCalendarSharesParams{
		OwnerID:    userId,
		CalendarID: calendarId,
	}); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}

// (PATCH /calendar/:calendarId/share/edit)
func (s *Server) CalendarShareEdit(c *echo.Context) error {
	userId := c.Get("userId").(uuid.UUID)
	calendarId, err := echo.PathParam[uuid.UUID](c, "calendarId")
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	body := new(models.ShareCalendar)
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	ctx := c.Request().Context()
	if err = s.queries.EditSharedCalendarPermissions(ctx, sqlc.EditSharedCalendarPermissionsParams{
		OwnerID:    userId,
		CalendarID: calendarId,
		Permission: body.Permission,
		SharedWith: body.UserId,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}

// (PATCH /calendar/:calendarId/share/edit/batch)
func (s *Server) BatchCalendarShareEdit(c *echo.Context) error {
	userId := c.Get("userId").(uuid.UUID)
	calendarId, err := echo.PathParam[uuid.UUID](c, "calendarId")
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	body := new(models.BatchShareCalendar)
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	batchParams := make([]sqlc.BatchEditSharedCalendarPermissionsParams, len(body.Items))
	for i, body := range body.Items {
		batchParams[i] = sqlc.BatchEditSharedCalendarPermissionsParams{
			OwnerID:    userId,
			CalendarID: calendarId,
			SharedWith: body.UserId,
			Permission: body.Permission,
		}
	}

	ctx := c.Request().Context()
	batchResults := s.queries.BatchEditSharedCalendarPermissions(ctx, batchParams)

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

// (DELETE /:calendarId/share/remove/self)
func (s *Server) RemoveUserCalendar(c *echo.Context) error {
	userId := c.Get("userId").(uuid.UUID)
	calendarId, err := echo.PathParam[uuid.UUID](c, "calendarId")
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	ctx := c.Request().Context()
	if err := s.queries.RemoveCalendarShareSelf(ctx, sqlc.RemoveCalendarShareSelfParams{
		CalendarID: calendarId,
		SharedWith: userId,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}

// (POST /:calendarId/share/remove/batch)
func (s *Server) BatchRemoveUserCalendar(c *echo.Context) error {
	userId := c.Get("userId").(uuid.UUID)
	calendarId, err := echo.PathParam[uuid.UUID](c, "calendarId")
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	body := new(models.BatchRemoveUserShare)
	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	ctx := c.Request().Context()
	if err := s.queries.RemoveCalendarShareMany(ctx, sqlc.RemoveCalendarShareManyParams{
		CalendarID:    calendarId,
		SharedWithIds: body.Ids,
		OwnerID:       userId,
	}); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}
