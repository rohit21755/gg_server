package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func setupREST(r chi.Router, db *gorm.DB) {
	log.Println("Setting up REST API")
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthHandler)
		r.Route("/users", func(r chi.Router) {
			r.Route("/auth", func(r chi.Router) {
				r.Post("/login", loginHandler(db))
				r.Post("/register", registerHandler(db))
				r.Post("/forgot-password", forgotPasswordHandler(db))
				r.Post("/forgot-username", forgotUsernameHandler(db))
			})
		})
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
