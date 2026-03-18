package routes

import (
	"github.com/eeritvan/calendar/internal/api"
	"github.com/labstack/echo/v5"
)

func folderRoutes(e *echo.Group, s *api.Server) {
	g := e.Group("/folders")

	g.POST("/new", s.NewFolder)
	g.POST("/add/:calendarId/:folderId", s.AddCalendarToFolder)
	g.DELETE("/remove/:calendarId", s.RemoveCalendarFromFolder)
	g.PATCH("/edit/:folderId", s.EditFolder)
	g.DELETE("/delete/:folderId", s.DeleteFolder)
}
