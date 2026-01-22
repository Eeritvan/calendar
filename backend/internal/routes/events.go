package routes

import (
	"github.com/eeritvan/calendar/internal/api"
	"github.com/labstack/echo/v5"
)

func eventRoutes(e *echo.Group, s *api.Server) {
	g := e.Group("/event")

	g.GET("/getEvents", s.GetEvents)
	g.GET("/searchEvents", s.SearchEvents)
	g.POST("/addEvent", s.AddEvent)
	g.PATCH("/edit/:eventId", s.EventEdit)
	g.DELETE("/delete/:eventId", s.EventDelete)
}
