package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eeritvan/calendar/internal/api"
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

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
			case "/signup", "/login", "/totp/authenticate", "/totp/recovery":
				return true
			}
			return false
		},
	}))

	api.RegisterHandlers(e, server)

	port := os.Getenv("PORT")
	log.Fatal(e.Start("0.0.0.0:" + port))
}
