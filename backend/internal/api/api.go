package api

import (
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/eeritvan/calendar/internal/stream"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
	sse     *stream.SSEHandler
}

func NewServer(queries *sqlc.Queries, pool *pgxpool.Pool, sse *stream.SSEHandler) *Server {
	return &Server{
		queries: queries,
		pool:    pool,
		sse:     sse,
	}
}
