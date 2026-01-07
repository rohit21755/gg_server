package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

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

