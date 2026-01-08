package main

import (
	"log"

	"github.com/rohit21755/gg_server.git/internal/db"
	"github.com/rohit21755/gg_server.git/internal/env"
)

func main() {
	env.Load()
	database := db.Connect()

	if err := db.Seed(database); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}
