package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eeritvan/calendar/internal/api"
	"github.com/eeritvan/calendar/internal/sqlc"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func wsHandler(c echo.Context) error {
	conn, err := websocket.Accept(c.Response(), c.Request(), &websocket.AcceptOptions{
		InsecureSkipVerify: true, // TODO
	})
	if err != nil {
		return err
	}
	defer conn.CloseNow()

	ctx := c.Request().Context()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err = wsjson.Write(ctx, conn, "testing")
			if err != nil {
				// TODO: error handling
				fmt.Println(err)
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

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

	server := api.NewServer(queries, pool)

	e := echo.New()

	basePath := e.Group("/api")

	e.Use(middleware.BodyLimit("500KB"))

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
			case "/api/signup", "/api/login", "/api/totp/authenticate", "/api/totp/recovery", "/ws":
				return true
			}
			return false
		},
	}))

	e.GET("/ws", wsHandler)
	api.RegisterHandlers(basePath, server)

	port := os.Getenv("PORT")
	log.Fatal(e.Start("0.0.0.0:" + port))
}
