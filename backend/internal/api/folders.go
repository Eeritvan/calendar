package api

import (
	"fmt"
	"net/http"

	"github.com/eeritvan/calendar/internal/models"
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// (POST /folders/new)
func (s *Server) NewFolder(c *echo.Context) error {
	body := new(models.AddFolder)
	if err := c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.AddFolder(ctx, sqlc.AddFolderParams{
		Name:   body.Name,
		UserID: userId,
	})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := models.Folder{
		Id:   queryResp.ID,
		Name: queryResp.Name,
	}

	return c.JSON(http.StatusCreated, resp)
}

// (POST /folders/add/:calendarId/:folderId)
func (s *Server) AddCalendarToFolder(c *echo.Context) error {
	calendarId, err := echo.PathParam[uuid.UUID](c, "calendarId")
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	folderId, err := echo.PathParam[uuid.UUID](c, "folderId")
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.AddCalendarToFolder(ctx, sqlc.AddCalendarToFolderParams{
		ID:       calendarId,
		FolderID: folderId,
		OwnerID:  userId,
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
		// TODO: permission
		IsOwner: queryResp.IsOwner,
		Folder: &models.Folder{
			Id:   queryResp.FolderID,
			Name: queryResp.FolderName,
		},
	}

	return c.JSON(http.StatusOK, resp)
}

// (PATCH /folders/edit/:folderId)
func (s *Server) EditFolder(c *echo.Context) error {
	folderId, err := echo.PathParam[uuid.UUID](c, "folderId")
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	body := new(models.FolderEdit)
	if err := c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.EditFolder(ctx, sqlc.EditFolderParams{
		ID:     folderId,
		Name:   *body.Name,
		UserID: userId,
	})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := models.Folder{
		Id:   queryResp.ID,
		Name: queryResp.Name,
	}

	return c.JSON(http.StatusOK, resp)
}

// (DELETE /folders/remove/:calendarId)
func (s *Server) RemoveCalendarFromFolder(c *echo.Context) error {
	calendarId, err := echo.PathParam[uuid.UUID](c, "calendarId")
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	if err := s.queries.RemoveCalendarFromFolder(ctx, sqlc.RemoveCalendarFromFolderParams{
		ID:      calendarId,
		OwnerID: userId,
	}); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}

// (DELETE /folders/delete/:folderId)
// -- TODO: delete all calendars inside the folder as well?
func (s *Server) DeleteFolder(c *echo.Context) error {
	folderId, err := echo.PathParam[uuid.UUID](c, "folderId")
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	if err := s.queries.DeleteFolder(ctx, sqlc.DeleteFolderParams{
		ID:     folderId,
		UserID: userId,
	}); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}
