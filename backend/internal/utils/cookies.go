package utils

import (
	"net/http"
	"time"
)

func CreateJWTCookie(jwtToken string) *http.Cookie {
	return &http.Cookie{
		Name:     "access_token",
		Value:    jwtToken,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

func DeleteJWTCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
