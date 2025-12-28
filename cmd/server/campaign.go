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

// Get Campaigns
func getCampaignsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		status := r.URL.Query().Get("status")
		campaignType := r.URL.Query().Get("type")
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
		query := db.Model(&store.Campaign{})

		if status != "" {
			query = query.Where("status = ?", status)
		}
		if campaignType != "" {
			query = query.Where("campaign_type = ?", campaignType)
		}

		// Get total count
		var totalCount int64
		query.Count(&totalCount)

		// Get campaigns
		var campaigns []store.Campaign
		result := query.
			Order("start_date DESC").
			Offset(offset).
			Limit(limit).
			Find(&campaigns)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		response := map[string]interface{}{
			"campaigns": campaigns,
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

// Get Campaign Details
func getCampaignHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		campaignIDStr := chi.URLParam(r, "id")
		campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid campaign ID"))
			return
		}

		campaign, err := store.GetCampaignByID(db, uint(campaignID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("campaign not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Get tasks count
		var tasksCount int64
		db.Model(&store.Task{}).Where("campaign_id = ?", campaign.ID).Count(&tasksCount)

		// Get participant count
		var participantsCount int64
		db.Model(&store.Submission{}).
			Where("campaign_id = ?", campaign.ID).
			Group("user_id").
			Count(&participantsCount)

		response := map[string]interface{}{
			"campaign": campaign,
			"stats": map[string]interface{}{
				"tasks_count":          tasksCount,
				"participants_count":   participantsCount,
				"current_participants": campaign.CurrentParticipants,
			},
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Join Campaign
func joinCampaignHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		campaignIDStr := chi.URLParam(r, "id")
		campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid campaign ID"))
			return
		}

		campaign, err := store.GetCampaignByID(db, uint(campaignID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("campaign not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Check if campaign is active
		if campaign.Status != "active" {
			badRequestResponse(w, r, errors.New("campaign is not active"))
			return
		}

		// Check if campaign has started
		if time.Now().Before(campaign.StartDate) {
			badRequestResponse(w, r, errors.New("campaign has not started yet"))
			return
		}

		// Check if campaign has ended
		if time.Now().After(campaign.EndDate) {
			badRequestResponse(w, r, errors.New("campaign has ended"))
			return
		}

		// Check max participants
		if campaign.MaxParticipants != nil && *campaign.MaxParticipants > 0 && campaign.CurrentParticipants >= *campaign.MaxParticipants {
			conflictResponse(w, r, errors.New("campaign has reached maximum participants"))
			return
		}

		// Check if user has already joined (has submissions in this campaign)
		var existingSubmissions int64
		db.Model(&store.Submission{}).
			Where("user_id = ? AND campaign_id = ?", user.ID, campaign.ID).
			Count(&existingSubmissions)

		if existingSubmissions > 0 {
			conflictResponse(w, r, errors.New("you have already joined this campaign"))
			return
		}

		// Increment participant count
		campaign.CurrentParticipants++
		if err := store.UpdateCampaign(db, campaign); err != nil {
			internalServerError(w, r, err)
			return
		}

		// Create activity log
		userIDInt := int(user.ID)
		activityDataJSON, _ := json.Marshal(map[string]interface{}{
			"campaign_id":   campaign.ID,
			"campaign_name": campaign.Title,
		})
		activityLog := &store.ActivityLog{
			UserID:       &userIDInt,
			ActivityType: "campaign_joined",
			ActivityData: stringPtr(string(activityDataJSON)),
		}
		store.CreateActivityLog(db, activityLog)

		response := map[string]string{
			"message": "Successfully joined campaign",
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Campaign Tasks
func getCampaignTasksHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		campaignIDStr := chi.URLParam(r, "id")
		campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid campaign ID"))
			return
		}

		// Get tasks for campaign
		var tasks []store.Task
		result := db.Where("campaign_id = ? AND is_active = ?", campaignID, true).
			Order("priority DESC, created_at DESC").
			Find(&tasks)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		if err := jsonResponse(w, http.StatusOK, tasks); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Campaign Leaderboard
func getCampaignLeaderboardHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		campaignIDStr := chi.URLParam(r, "id")
		campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid campaign ID"))
			return
		}

		// Get top performers in campaign
		var leaderboard []struct {
			UserID      uint   `json:"user_id"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			College     string `json:"college"`
			TotalXP     int    `json:"total_xp"`
			Submissions int    `json:"submissions"`
			Rank        int    `json:"rank"`
		}

		result := db.Model(&store.Submission{}).
			Select(`
				users.id as user_id,
				users.first_name,
				users.last_name,
				colleges.name as college,
				SUM(submissions.xp_awarded) as total_xp,
				COUNT(submissions.id) as submissions
			`).
			Joins("JOIN users ON users.id = submissions.user_id").
			Joins("LEFT JOIN colleges ON colleges.id = users.college_id").
			Where("submissions.campaign_id = ? AND submissions.status = 'approved'", campaignID).
			Group("users.id, colleges.name").
			Order("total_xp DESC").
			Limit(20).
			Scan(&leaderboard)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Add ranks
		for i := range leaderboard {
			leaderboard[i].Rank = i + 1
		}

		if err := jsonResponse(w, http.StatusOK, leaderboard); err != nil {
			internalServerError(w, r, err)
		}
	}
}
