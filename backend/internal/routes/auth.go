package routes

import (
	"github.com/eeritvan/calendar/internal/api"
	"github.com/labstack/echo/v5"
)

func authRoutes(e *echo.Group, s *api.Server) {
	g := e.Group("/auth")

	g.POST("/signup", s.PostSignup)
	g.POST("/login", s.PostLogin)

	g.POST("/totp/enable", s.PostTotpEnable)
	g.PATCH("/totp/enable/verify", s.PatchTotpEnableVerify)
	g.PATCH("/totp/disable", s.PatchTotpDisable)
	g.POST("/totp/authenticate", s.PostTotpAuthenticate)
	g.POST("/totp/recovery", s.PostTotpRecovery)
}
