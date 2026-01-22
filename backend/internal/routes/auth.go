package routes

import (
	"github.com/eeritvan/calendar/internal/api"
	"github.com/labstack/echo/v5"
)

func authRoutes(e *echo.Group, s *api.Server) {
	g := e.Group("/auth")

	g.POST("/signup", s.Signup)
	g.POST("/login", s.Login)

	g.POST("/totp/enable", s.TotpEnable)
	g.PATCH("/totp/enable/verify", s.TotpEnableVerify)
	g.PATCH("/totp/disable", s.TotpDisable)
	g.POST("/totp/authenticate", s.TotpAuthenticate)
	g.POST("/totp/recovery", s.TotpRecovery)
}
