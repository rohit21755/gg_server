package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rohit21755/gg_server.git/graph"
	"github.com/rohit21755/gg_server.git/internal/db"
	"github.com/rohit21755/gg_server.git/internal/env"
	"github.com/rohit21755/gg_server.git/ws"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	env.Load()
	database := db.Connect()

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// WebSockets
	hub := ws.NewHub()
	go hub.Run()

	// GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: database}}))
	router.Handle("/graphql", srv)

	// REST API
	setupREST(router, database)

	// WebSocket endpoint
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})

	log.Println("Server running on :" + os.Getenv("SERVER_PORT"))
	http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), router)
}
