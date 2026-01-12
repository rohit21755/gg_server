package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// Helper function to check if string contains substring (case-insensitive)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Get all colleges
func getCollegesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var colleges []store.College
		if err := db.Where("is_active = ?", true).Find(&colleges).Error; err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch colleges")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"colleges": colleges,
		})
	}
}

// Get college by ID
func getCollegeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collegeIDStr := chi.URLParam(r, "id")
		collegeID, err := strconv.ParseUint(collegeIDStr, 10, 32)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid college ID")
			return
		}

		college, err := store.GetCollegeByID(db, uint(collegeID))
		if err != nil {
			writeJSONError(w, http.StatusNotFound, "college not found")
			return
		}

		writeJSON(w, http.StatusOK, college)
	}
}

// Get college leaderboard
func getCollegeLeaderboardHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collegeIDStr := chi.URLParam(r, "id")
		collegeID, err := strconv.ParseUint(collegeIDStr, 10, 32)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid college ID")
			return
		}

		limit := 100
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
				limit = l
			}
		}

		var users []store.User
		if err := db.Where("college_id = ? AND is_active = ?", collegeID, true).
			Order("xp DESC").
			Limit(limit).
			Find(&users).Error; err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch leaderboard")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"leaderboard": users,
		})
	}
}

// Get all states
func getStatesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var states []store.State
		if err := db.Find(&states).Error; err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch states")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"states": states,
		})
	}
}

// Get state by ID
func getStateHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stateIDStr := chi.URLParam(r, "id")
		stateID, err := strconv.ParseUint(stateIDStr, 10, 32)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid state ID")
			return
		}

		state, err := store.GetStateByID(db, uint(stateID))
		if err != nil {
			writeJSONError(w, http.StatusNotFound, "state not found")
			return
		}

		writeJSON(w, http.StatusOK, state)
	}
}

// Get state leaderboard
func getStateLeaderboardHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stateIDStr := chi.URLParam(r, "id")
		stateID, err := strconv.ParseUint(stateIDStr, 10, 32)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid state ID")
			return
		}

		limit := 100
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
				limit = l
			}
		}

		var users []store.User
		if err := db.Where("state_id = ? AND is_active = ?", stateID, true).
			Order("xp DESC").
			Limit(limit).
			Find(&users).Error; err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch leaderboard")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"leaderboard": users,
		})
	}
}

// Create or get college - public endpoint (no auth required)
func createOrGetCollegeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var req struct {
			Name    string  `json:"name" validate:"required,min=1,max=200"`
			Code    *string `json:"code,omitempty" validate:"omitempty,max=50"`
			StateID *int    `json:"state_id,omitempty"`
		}

		if err := readJSON(w, r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Validate request
		if err := Validate.Struct(req); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Check if college already exists by name or code
		existingCollege, err := store.GetCollegeByNameOrCode(db, req.Name, req.Code)
		if err == nil {
			// College exists, return it
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"college": existingCollege,
				"message": "College already exists",
			})
			return
		}

		// Validate state_id if provided
		if req.StateID != nil {
			_, err := store.GetStateByID(db, uint(*req.StateID))
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, "invalid state_id")
				return
			}
		}

		// Create new college
		college := &store.College{
			Name:    req.Name,
			Code:    req.Code,
			StateID: req.StateID,
			IsActive: true,
		}

		if err := store.CreateCollege(db, college); err != nil {
			// Check if it's a unique constraint violation (PostgreSQL)
			errStr := err.Error()
			if contains(errStr, "duplicate key") || contains(errStr, "UNIQUE constraint") || contains(errStr, "colleges_code_key") {
				writeJSONError(w, http.StatusConflict, "college with this name or code already exists")
				return
			}
			writeJSONError(w, http.StatusInternalServerError, "failed to create college")
			return
		}

		// Load the created college with relations
		createdCollege, err := store.GetCollegeByID(db, college.ID)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch created college")
			return
		}

		writeJSON(w, http.StatusCreated, map[string]interface{}{
			"college": createdCollege,
			"message": "College created successfully",
		})
	}
}
