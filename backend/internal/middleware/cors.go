package middleware

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/rs/cors"
)

func CorsMiddleware(srv *handler.Server) http.Handler {
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(srv)
	return handler
}
