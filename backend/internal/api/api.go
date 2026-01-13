package api

import (
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3labs/sse/v2"
)

type Server struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
	sse     *sse.Server
}

func NewServer(queries *sqlc.Queries, pool *pgxpool.Pool, sse *sse.Server) *Server {
	return &Server{
		queries: queries,
		pool:    pool,
		sse:     sse,
	}
}
