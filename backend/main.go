package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eeritvan/calendar/api"
	"github.com/eeritvan/calendar/internal/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	dbUrl := os.Getenv("DB_URL")
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		// TODO: error handling
	}

	queries := sqlc.New(pool)

	server := api.NewServer(queries)

	e := echo.New()

	api.RegisterHandlers(e, server)

	log.Fatal(e.Start("0.0.0.0:8080"))
}
