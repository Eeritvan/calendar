package api

import (
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewServer(queries *sqlc.Queries, pool *pgxpool.Pool) *Server {
	return &Server{
		queries: queries,
		pool:    pool,
	}
}
