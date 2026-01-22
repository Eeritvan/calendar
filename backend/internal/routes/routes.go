package routes

import (
	"github.com/eeritvan/calendar/internal/api"
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo, s *api.Server) {
	g := e.Group("/api")

	authRoutes(g, s)
	calendarRoutes(g, s)
	eventRoutes(g, s)
}
