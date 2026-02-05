package routes

import (
	"github.com/eeritvan/calendar/internal/api"
	"github.com/labstack/echo/v5"
)

func calendarRoutes(e *echo.Group, s *api.Server) {
	g := e.Group("/calendar")

	g.GET("/getCalendars", s.GetCalendars)
	g.POST("/addCalendar", s.AddCalendar)
	g.PATCH("/edit/:calendarId", s.EditCalendar)
	g.DELETE("/delete/:calendarId", s.DeleteCalendar)
	g.POST("/:calendarId/event/import", s.ImportEvents)
}
