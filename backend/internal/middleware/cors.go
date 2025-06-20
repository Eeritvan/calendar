package middleware

import (
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/rs/cors"
)

const DEFAULT_URL = "http://localhost:5173"

func CorsMiddleware(srv *handler.Server) http.Handler {
	frontendUrl := os.Getenv("FRONTEND_ORIGIN")
	if frontendUrl == "" {
		frontendUrl = DEFAULT_URL
	}

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{frontendUrl},
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(srv)
	return handler
}
