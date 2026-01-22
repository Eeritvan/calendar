package routes

import (
	"github.com/eeritvan/calendar/internal/api"
	"github.com/labstack/echo/v5"
)

func calendarRoutes(e *echo.Group, s *api.Server) {
	g := e.Group("/calendar")

	g.GET("/getCalendars", s.GetGetCalendars)
	g.POST("/addCalendar", s.PostAddCalendar)
	g.PATCH("/edit/:calendarId", s.PatchCalendarEditCalendarId)
	g.DELETE("/delete/:calendarId", s.DeleteCalendarDeleteCalendarId)
}
