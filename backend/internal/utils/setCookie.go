package utils

import (
	"net/http"
	"time"
)

func CreateCookie(jwtToken string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "access_token"
	cookie.Value = jwtToken
	cookie.Expires = time.Now().Add(24 * time.Hour)

	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteLaxMode

	return cookie
}
