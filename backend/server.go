package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/eeritvan/calendar/graph"
	"github.com/eeritvan/calendar/internal/db"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/eeritvan/calendar/internal/middleware"
)

const defaultPort = "8081"

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		// todo: better error handling
		log.Print("health check failed")
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	ctx := context.Background()

	dbService, err := db.ConnectToDB(ctx)
	if err != nil {
		log.Fatal("failed to initialize database service")
	}
	defer func() {
		if err := dbService.Conn.Close(ctx); err != nil {
			log.Printf("failed to close DB connection: %v", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			DB: dbService,
		},
	}))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				return origin == "http://localhost:5173" || origin == "ws://localhost:5173"
			},
		},
	})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	handler := middleware.CorsMiddleware(srv)

	http.Handle("/healthz", http.HandlerFunc(healthCheck))
	http.Handle("/api", handler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
