package routes

import (
	"github.com/eeritvan/calendar/internal/api"
	"github.com/labstack/echo/v5"
)

func eventRoutes(e *echo.Group, s *api.Server) {
	g := e.Group("/event")

	g.GET("/getEvents", s.GetGetEvents)
	g.GET("/searchEvents", s.GetSearchEvents)
	g.POST("/addEvent", s.PostAddEvent)
	g.PATCH("/edit/:eventId", s.PatchEventEditEventId)
	g.DELETE("/delete/:eventId", s.DeleteEventDeleteEventId)
}
