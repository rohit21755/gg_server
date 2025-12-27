package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rohit21755/gg_server.git/internal/store"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterUserPayload struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Phone      string `json:"phone,omitempty"`
	CollegeID  *int   `json:"college_id,omitempty"`
	StateID    *int   `json:"state_id,omitempty"`
	ReferredBy *int   `json:"referred_by,omitempty"`
}
type LoggingInUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func loginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload LoggingInUserPayload
		if err := readJSON(w, r, &payload); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		if err := Validate.Struct(payload); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		user, err := store.GetUserByEmail(db, payload.Email)
		if err != nil {
			notFoundResponse(w, r, err)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
			unauthorizedResponse(w, r, err)
			return
		}

		userID := int(user.ID)
		session := store.UserSession{
			UserID:       &userID,
			SessionToken: uuid.New().String(),
			ExpiresAt:    time.Now().Add(24 * time.Hour),
			CreatedAt:    time.Now(),
		}
		if err := store.CreateSession(db, &session); err != nil {
			internalServerError(w, r, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"id": user.ID, "email": user.Email, "session_token": session.SessionToken})
	}
}

func registerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload RegisterUserPayload
		if err := readJSON(w, r, &payload); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		if err := Validate.Struct(payload); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		user := store.User{
			Email:        payload.Email,
			PasswordHash: string(hashed),
			FirstName:    payload.FirstName,
			LastName:     payload.LastName,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if payload.Phone != "" {
			user.Phone = &payload.Phone
		}
		user.CollegeID = payload.CollegeID
		user.StateID = payload.StateID
		user.ReferredBy = payload.ReferredBy

		if err := store.CreateUser(db, &user); err != nil {
			internalServerError(w, r, err)
			return
		}

		writeJSON(w, http.StatusCreated, map[string]any{"id": user.ID, "email": user.Email})
	}
}

func forgotPasswordHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSONError(w, http.StatusNotImplemented, "forgot-password handler not implemented")
	}
}

func forgotUsernameHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSONError(w, http.StatusNotImplemented, "forgot-username handler not implemented")
	}
}
