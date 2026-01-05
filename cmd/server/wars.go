package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// Get Active Wars
func getActiveWarsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		status := r.URL.Query().Get("status")
		warType := r.URL.Query().Get("type")

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit

		// Build query
		query := db.Model(&store.CampusWar{})

		if status != "" {
			query = query.Where("status = ?", status)
		} else {
			// Default to active and upcoming wars
			query = query.Where("status IN (?, ?)", "active", "upcoming")
		}

		if warType != "" {
			query = query.Where("war_type = ?", warType)
		}

		// Get total count
		var totalCount int64
		query.Count(&totalCount)

		// Get wars
		var wars []store.CampusWar
		result := query.
			Order("start_date DESC").
			Offset(offset).
			Limit(limit).
			Find(&wars)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Format response
		var responseWars []map[string]interface{}
		for _, war := range wars {
			warData := map[string]interface{}{
				"id":          war.ID,
				"name":        war.Name,
				"description": war.Description,
				"war_type":    war.WarType,
				"start_date":  war.StartDate,
				"end_date":    war.EndDate,
				"status":      war.Status,
				"created_at":  war.CreatedAt,
			}

			// Parse metrics JSON
			if war.Metrics != nil {
				var metrics map[string]interface{}
				if err := json.Unmarshal([]byte(*war.Metrics), &metrics); err == nil {
					warData["metrics"] = metrics
				}
			}

			// Parse rewards JSON
			if war.Rewards != nil {
				var rewards map[string]interface{}
				if err := json.Unmarshal([]byte(*war.Rewards), &rewards); err == nil {
					warData["rewards"] = rewards
				}
			}

			// Get participant count
			var participantCount int64
			db.Model(&store.WarParticipant{}).Where("war_id = ?", war.ID).Count(&participantCount)
			warData["participant_count"] = participantCount

			responseWars = append(responseWars, warData)
		}

		response := map[string]interface{}{
			"wars": responseWars,
			"pagination": map[string]interface{}{
				"page":        page,
				"limit":       limit,
				"total":       totalCount,
				"total_pages": (int(totalCount) + limit - 1) / limit,
			},
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get War Details
func getWarHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warIDStr := chi.URLParam(r, "id")
		warID, err := strconv.ParseUint(warIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid war ID"))
			return
		}

		war, err := store.GetCampusWarByID(db, uint(warID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("war not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Get participant count
		var participantCount int64
		db.Model(&store.WarParticipant{}).Where("war_id = ?", war.ID).Count(&participantCount)

		// Format response
		response := map[string]interface{}{
			"id":                war.ID,
			"name":              war.Name,
			"description":       war.Description,
			"war_type":          war.WarType,
			"start_date":        war.StartDate,
			"end_date":          war.EndDate,
			"status":            war.Status,
			"created_at":        war.CreatedAt,
			"participant_count": participantCount,
		}

		// Parse metrics JSON
		if war.Metrics != nil {
			var metrics map[string]interface{}
			if err := json.Unmarshal([]byte(*war.Metrics), &metrics); err == nil {
				response["metrics"] = metrics
			}
		}

		// Parse rewards JSON
		if war.Rewards != nil {
			var rewards map[string]interface{}
			if err := json.Unmarshal([]byte(*war.Rewards), &rewards); err == nil {
				response["rewards"] = rewards
			}
		}

		// Calculate time remaining if active
		if war.Status == "active" {
			now := time.Now()
			if war.EndDate.After(now) {
				response["time_remaining"] = war.EndDate.Sub(now).String()
			}
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get War Participants
func getWarParticipantsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warIDStr := chi.URLParam(r, "id")
		warID, err := strconv.ParseUint(warIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid war ID"))
			return
		}

		// Verify war exists
		war, err := store.GetCampusWarByID(db, uint(warID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("war not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Parse query parameters
		entityType := r.URL.Query().Get("entity_type")

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 || limit > 100 {
			limit = 50
		}
		offset := (page - 1) * limit

		// Build query
		query := db.Model(&store.WarParticipant{}).Where("war_id = ?", warID)

		if entityType != "" {
			query = query.Where("entity_type = ?", entityType)
		}

		// Get total count
		var totalCount int64
		query.Count(&totalCount)

		// Get participants
		var participants []store.WarParticipant
		result := query.
			Order("rank ASC NULLS LAST, total_xp DESC, total_submissions DESC").
			Offset(offset).
			Limit(limit).
			Find(&participants)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Format response with entity details
		var responseParticipants []map[string]interface{}
		for _, participant := range participants {
			participantData := map[string]interface{}{
				"id":                participant.ID,
				"entity_type":       participant.EntityType,
				"entity_id":         participant.EntityID,
				"total_xp":          participant.TotalXP,
				"total_submissions": participant.TotalSubmissions,
				"total_referrals":   participant.TotalReferrals,
				"rank":              participant.Rank,
				"created_at":        participant.CreatedAt,
			}

			// Get entity details based on type
			if participant.EntityType == "college" {
				var college store.College
				if err := db.Preload("State").First(&college, participant.EntityID).Error; err == nil {
					participantData["entity"] = map[string]interface{}{
						"id":    college.ID,
						"name":  college.Name,
						"code":  college.Code,
						"state": nil,
					}
					if college.State != nil {
						participantData["entity"].(map[string]interface{})["state"] = map[string]interface{}{
							"id":   college.State.ID,
							"name": college.State.Name,
							"code": college.State.Code,
						}
					}
				}
			} else if participant.EntityType == "state" {
				var state store.State
				if err := db.First(&state, participant.EntityID).Error; err == nil {
					participantData["entity"] = map[string]interface{}{
						"id":   state.ID,
						"name": state.Name,
						"code": state.Code,
					}
				}
			}

			responseParticipants = append(responseParticipants, participantData)
		}

		response := map[string]interface{}{
			"war": map[string]interface{}{
				"id":     war.ID,
				"name":   war.Name,
				"status": war.Status,
			},
			"participants": responseParticipants,
			"pagination": map[string]interface{}{
				"page":        page,
				"limit":       limit,
				"total":       totalCount,
				"total_pages": (int(totalCount) + limit - 1) / limit,
			},
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get War Leaderboard
func getWarLeaderboardHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warIDStr := chi.URLParam(r, "id")
		warID, err := strconv.ParseUint(warIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid war ID"))
			return
		}

		// Verify war exists
		war, err := store.GetCampusWarByID(db, uint(warID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("war not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Parse query parameters
		entityType := r.URL.Query().Get("entity_type")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 || limit > 100 {
			limit = 50
		}

		// Build query for leaderboard (sorted by rank, then by metrics)
		query := db.Model(&store.WarParticipant{}).Where("war_id = ?", warID)

		if entityType != "" {
			query = query.Where("entity_type = ?", entityType)
		}

		// Get participants ordered by rank and metrics
		var participants []store.WarParticipant
		result := query.
			Order("rank ASC NULLS LAST, total_xp DESC, total_submissions DESC, total_referrals DESC").
			Limit(limit).
			Find(&participants)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Format leaderboard response
		var leaderboard []map[string]interface{}
		for idx, participant := range participants {
			leaderboardData := map[string]interface{}{
				"rank":              participant.Rank,
				"position":          idx + 1, // Display position (1-based)
				"entity_type":       participant.EntityType,
				"entity_id":         participant.EntityID,
				"total_xp":          participant.TotalXP,
				"total_submissions": participant.TotalSubmissions,
				"total_referrals":   participant.TotalReferrals,
				"total_score":       participant.TotalXP + (participant.TotalSubmissions * 10) + (participant.TotalReferrals * 50), // Simple scoring
			}

			// Get entity details
			if participant.EntityType == "college" {
				var college store.College
				if err := db.Preload("State").First(&college, participant.EntityID).Error; err == nil {
					leaderboardData["entity_name"] = college.Name
					leaderboardData["entity_code"] = college.Code
					if college.State != nil {
						leaderboardData["state"] = map[string]interface{}{
							"id":   college.State.ID,
							"name": college.State.Name,
							"code": college.State.Code,
						}
					}
				}
			} else if participant.EntityType == "state" {
				var state store.State
				if err := db.First(&state, participant.EntityID).Error; err == nil {
					leaderboardData["entity_name"] = state.Name
					leaderboardData["entity_code"] = state.Code
				}
			}

			leaderboard = append(leaderboard, leaderboardData)
		}

		// Parse metrics to determine scoring
		var metrics map[string]interface{}
		if war.Metrics != nil {
			json.Unmarshal([]byte(*war.Metrics), &metrics)
		}

		response := map[string]interface{}{
			"war": map[string]interface{}{
				"id":         war.ID,
				"name":       war.Name,
				"status":     war.Status,
				"start_date": war.StartDate,
				"end_date":   war.EndDate,
				"metrics":    metrics,
			},
			"leaderboard":        leaderboard,
			"total_participants": len(participants),
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}
