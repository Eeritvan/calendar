package api

import "github.com/eeritvan/calendar/internal/sqlc"

type Server struct {
	queries *sqlc.Queries
}

func NewServer(queries *sqlc.Queries) *Server {
	return &Server{queries: queries}
}
