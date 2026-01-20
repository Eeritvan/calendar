package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eeritvan/calendar/internal/api"
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/eeritvan/calendar/internal/stream"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/r3labs/sse/v2"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		fmt.Println("Error loading .env file")
	}

	dbUrl := os.Getenv("DB_URL")
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		// TODO: error handling
	}

	queries := sqlc.New(pool)

	sseServer := sse.New()
	sseServer.AutoReplay = false

	server := api.NewServer(queries, pool, sseServer)

	e := echo.New()

	e.Use(middleware.BodyLimit("500KB"))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	JWTkey := os.Getenv("JWT_KEY")
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(JWTkey),
		TokenLookup: "cookie:access_token",
		SuccessHandler: func(c echo.Context) {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)

			userIdStr := claims["userId"].(string)
			uid, err := uuid.Parse(userIdStr)
			if err != nil {
				// TODO: errro handling
				return
			}

			c.Set("userId", uid)
		},
		Skipper: func(c echo.Context) bool {
			switch c.Path() {
			case "/api/signup", "/api/login", "/api/totp/authenticate", "/api/totp/recovery":
				return true
			}
			return false
		},
	}))

	sseHandler := &stream.SSEHandler{
		SSEServer: sseServer,
	}

	basePath := e.Group("/api")
	basePath.GET("/sse", sseHandler.HandleSSE)
	api.RegisterHandlers(basePath, server)

	port := os.Getenv("PORT")
	log.Fatal(e.Start("0.0.0.0:" + port))
}
