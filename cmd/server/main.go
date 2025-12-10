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
	"github.com/go-chi/chi/v5"
)

func main() {
	env.Load()
	database := db.Connect()

	router := chi.NewRouter()

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
