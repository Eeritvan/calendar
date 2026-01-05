package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eeritvan/calendar/api"
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/golang-jwt/jwt/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
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

	JWTkey := os.Getenv("JWT_KEY")
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(JWTkey),
		SuccessHandler: func(c echo.Context) {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)

			userId := claims["userId"].(string)
			username := claims["username"].(string)

			c.Set("userId", userId)
			c.Set("username", username)
		},
		Skipper: func(c echo.Context) bool {
			switch c.Path() {
			case "/signup", "/login":
				return true
			}
			return false
		},
	}))

	api.RegisterHandlers(e, server)

	log.Fatal(e.Start("0.0.0.0:8080"))
}
