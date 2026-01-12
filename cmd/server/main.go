package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rohit21755/gg_server.git/internal/db"
	"github.com/rohit21755/gg_server.git/internal/env"
	"github.com/rohit21755/gg_server.git/ws"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	env.Load()
	log.Println("Loading environment variables...")
	
	database := db.Connect()
	if database == nil {
		log.Fatal("Failed to connect to database")
	}
	log.Println("Database connected successfully")

	router := chi.NewRouter()
	log.Println("Router created")
	router.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// WebSockets
	hub := ws.NewHub()
	go hub.Run()

	// REST API
	setupREST(router, database)

	// WebSocket endpoint
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})

	log.Println("Server running on :" + os.Getenv("SERVER_PORT"))
	http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), router)
}
