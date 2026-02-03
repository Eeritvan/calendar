package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/eeritvan/calendar/internal/api"
	"github.com/eeritvan/calendar/internal/routes"
	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/eeritvan/calendar/internal/stream"
	"github.com/eeritvan/calendar/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/r3labs/sse/v2"
)

//go:embed "dist"
var dist embed.FS

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

	sseHandler := &stream.SSEHandler{
		SSEServer:   sseServer,
		UserClients: make(map[uuid.UUID]map[string]struct{}),
	}

	server := api.NewServer(queries, pool, sseHandler)

	e := echo.New()
	e.Validator = &utils.CustomValidator{
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	e.Use(middleware.BodyLimit(524_288)) // 500kb

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	JWTkey := os.Getenv("JWT_KEY")
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(JWTkey),
		TokenLookup: "cookie:access_token",
		SuccessHandler: func(c *echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			userIdStr := claims["userId"].(string)

			uid, err := uuid.Parse(userIdStr)
			if err != nil {
				// TODO: errro handling
				return nil
			}
			c.Set("userId", uid)

			return nil
		},
		Skipper: func(c *echo.Context) bool {
			switch c.Path() {
			case "/*", "/api/auth/signup", "/api/auth/login", "/api/auth/totp/authenticate", "/api/auth/totp/recovery":
				return true
			}
			return false
		},
	}))

	e.GET("/api/sse", sseHandler.HandleSSE)
	routes.RegisterRoutes(e, server)

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "dist",
		Filesystem: dist,
	}))
	e.StaticFS("/", echo.MustSubFS(dist, "dist"))

	port := os.Getenv("PORT")
	log.Fatal(e.Start("0.0.0.0:" + port))
}
