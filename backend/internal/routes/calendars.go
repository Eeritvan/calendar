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
	g.GET("/:calendarId/event/export", s.ExportEvents)

	// sharing
	g.POST("/:calendarId/share", s.ShareCalendar)
	g.POST("/:calendarId/share/batch", s.BatchShareCalendar)
	g.PATCH("/:calendarId/share/public", s.ShareCalendarPublic)
	g.PATCH("/:calendarId/share/private", s.ShareCalendarPrivate)
	g.PATCH("/:calendarId/share/edit", s.CalendarShareEdit)
	g.DELETE("/:calendarId/share/remove/:userId", s.RemoveUserCalendar)
	g.POST("/:calendarId/share/remove/batch", s.BatchRemoveUserCalendar)
}
