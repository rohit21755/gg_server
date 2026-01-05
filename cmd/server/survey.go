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

// Get Available Surveys
func getAvailableSurveysHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Parse query parameters
		surveyType := r.URL.Query().Get("type")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit

		now := time.Now()

		// Build query for available surveys
		// Available surveys: is_active = true, (start_date IS NULL OR start_date <= now), (end_date IS NULL OR end_date >= now)
		query := db.Model(&store.Survey{}).
			Where("is_active = ?", true).
			Where("(start_date IS NULL OR start_date <= ?)", now).
			Where("(end_date IS NULL OR end_date >= ?)", now)

		if surveyType != "" {
			query = query.Where("survey_type = ?", surveyType)
		}

		// Exclude surveys the user has already submitted
		subQuery := db.Model(&store.SurveyResponse{}).
			Select("survey_id").
			Where("user_id = ?", user.ID)

		query = query.Where("id NOT IN (?)", subQuery)

		// Get total count
		var totalCount int64
		query.Count(&totalCount)

		// Get surveys
		var surveys []store.Survey
		result := query.
			Order("created_at DESC").
			Offset(offset).
			Limit(limit).
			Find(&surveys)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Format response (exclude questions for list view)
		var responseSurveys []map[string]interface{}
		for _, survey := range surveys {
			surveyData := map[string]interface{}{
				"id":          survey.ID,
				"title":       survey.Title,
				"description": survey.Description,
				"survey_type": survey.SurveyType,
				"xp_reward":   survey.XPReward,
				"start_date":  survey.StartDate,
				"end_date":    survey.EndDate,
				"created_at":  survey.CreatedAt,
			}

			responseSurveys = append(responseSurveys, surveyData)
		}

		response := map[string]interface{}{
			"surveys": responseSurveys,
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

// Get Survey Details
func getSurveyHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		surveyIDStr := chi.URLParam(r, "id")
		surveyID, err := strconv.ParseUint(surveyIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid survey ID"))
			return
		}

		survey, err := store.GetSurveyByID(db, uint(surveyID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("survey not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Check if survey is available
		now := time.Now()
		if !survey.IsActive {
			badRequestResponse(w, r, errors.New("survey is not active"))
			return
		}

		if survey.StartDate != nil && survey.StartDate.After(now) {
			badRequestResponse(w, r, errors.New("survey has not started yet"))
			return
		}

		if survey.EndDate != nil && survey.EndDate.Before(now) {
			badRequestResponse(w, r, errors.New("survey has ended"))
			return
		}

		// Check if user has already submitted
		var existingResponse store.SurveyResponse
		result := db.Where("survey_id = ? AND user_id = ?", survey.ID, user.ID).First(&existingResponse)
		if result.Error == nil {
			// User has already submitted
			conflictResponse(w, r, errors.New("you have already submitted this survey"))
			return
		}
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			internalServerError(w, r, result.Error)
			return
		}

		// Parse questions JSON
		var questions interface{}
		if err := json.Unmarshal([]byte(survey.Questions), &questions); err != nil {
			internalServerError(w, r, err)
			return
		}

		response := map[string]interface{}{
			"id":          survey.ID,
			"title":       survey.Title,
			"description": survey.Description,
			"survey_type": survey.SurveyType,
			"questions":   questions,
			"xp_reward":   survey.XPReward,
			"start_date":  survey.StartDate,
			"end_date":    survey.EndDate,
			"created_at":  survey.CreatedAt,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Submit Survey
func submitSurveyHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		surveyIDStr := chi.URLParam(r, "id")
		surveyID, err := strconv.ParseUint(surveyIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid survey ID"))
			return
		}

		// Get survey
		survey, err := store.GetSurveyByID(db, uint(surveyID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("survey not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Validate survey is available
		now := time.Now()
		if !survey.IsActive {
			badRequestResponse(w, r, errors.New("survey is not active"))
			return
		}

		if survey.StartDate != nil && survey.StartDate.After(now) {
			badRequestResponse(w, r, errors.New("survey has not started yet"))
			return
		}

		if survey.EndDate != nil && survey.EndDate.Before(now) {
			badRequestResponse(w, r, errors.New("survey has ended"))
			return
		}

		// Check if user has already submitted
		var existingResponse store.SurveyResponse
		result := db.Where("survey_id = ? AND user_id = ?", survey.ID, user.ID).First(&existingResponse)
		if result.Error == nil {
			conflictResponse(w, r, errors.New("you have already submitted this survey"))
			return
		}
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			internalServerError(w, r, result.Error)
			return
		}

		// Parse request body
		var req struct {
			Responses map[string]interface{} `json:"responses" validate:"required"`
		}

		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		if len(req.Responses) == 0 {
			badRequestResponse(w, r, errors.New("responses are required"))
			return
		}

		// Validate responses against questions (basic validation)
		var questions map[string]interface{}
		if err := json.Unmarshal([]byte(survey.Questions), &questions); err != nil {
			internalServerError(w, r, err)
			return
		}

		// Marshal responses to JSON string
		responsesJSON, err := json.Marshal(req.Responses)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		// Calculate completion percentage (simplified - assume 100% if responses provided)
		completionPercentage := 100

		// Start transaction
		tx := db.Begin()

		// Create survey response
		surveyResponse := &store.SurveyResponse{
			SurveyID:             intPtr(int(survey.ID)),
			UserID:               intPtr(int(user.ID)),
			Responses:            string(responsesJSON),
			CompletionPercentage: completionPercentage,
			XPAwarded:            survey.XPReward,
		}

		if err := store.CreateSurveyResponse(tx, surveyResponse); err != nil {
			tx.Rollback()
			internalServerError(w, r, err)
			return
		}

		// Award XP to user if reward > 0
		if survey.XPReward > 0 {
			dbUser, err := store.GetUserByID(tx, user.ID)
			if err != nil {
				tx.Rollback()
				internalServerError(w, r, err)
				return
			}

			dbUser.XP += survey.XPReward
			if err := tx.Save(dbUser).Error; err != nil {
				tx.Rollback()
				internalServerError(w, r, err)
				return
			}

			// Create XP transaction
			metadataJSON, err := json.Marshal(map[string]interface{}{
				"survey_id":          survey.ID,
				"survey_title":       survey.Title,
				"survey_response_id": surveyResponse.ID,
			})
			if err != nil {
				tx.Rollback()
				internalServerError(w, r, err)
				return
			}

			xpTransaction := &store.XPTransaction{
				UserID:          intPtr(int(user.ID)),
				TransactionType: "quiz", // Using quiz type for survey XP
				Amount:          survey.XPReward,
				SourceType:      stringPtr("survey"),
				SourceID:        intPtr(int(surveyResponse.ID)),
				Description:     stringPtr("Survey completion: " + survey.Title),
				Metadata:        stringPtr(string(metadataJSON)),
			}

			if err := store.CreateXPTransaction(tx, xpTransaction); err != nil {
				tx.Rollback()
				internalServerError(w, r, err)
				return
			}
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			internalServerError(w, r, err)
			return
		}

		// Get updated user for response
		dbUser, _ := store.GetUserByID(db, user.ID)

		response := map[string]interface{}{
			"message": "Survey submitted successfully",
			"response": map[string]interface{}{
				"id":                    surveyResponse.ID,
				"survey_id":             survey.ID,
				"completion_percentage": surveyResponse.CompletionPercentage,
				"xp_awarded":            surveyResponse.XPAwarded,
				"submitted_at":          surveyResponse.SubmittedAt,
			},
			"user": map[string]interface{}{
				"xp": dbUser.XP,
			},
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Survey Responses
func getSurveyResponsesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Parse query parameters
		surveyID, _ := strconv.ParseUint(r.URL.Query().Get("survey_id"), 10, 32)

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
		query := db.Model(&store.SurveyResponse{}).Where("user_id = ?", user.ID)

		if surveyID > 0 {
			query = query.Where("survey_id = ?", surveyID)
		}

		// Get total count
		var totalCount int64
		query.Count(&totalCount)

		// Get responses
		var responses []store.SurveyResponse
		result := query.
			Preload("Survey").
			Order("submitted_at DESC").
			Offset(offset).
			Limit(limit).
			Find(&responses)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Format response
		var responseResponses []map[string]interface{}
		for _, response := range responses {
			responseData := map[string]interface{}{
				"id":                    response.ID,
				"survey_id":             response.SurveyID,
				"completion_percentage": response.CompletionPercentage,
				"xp_awarded":            response.XPAwarded,
				"submitted_at":          response.SubmittedAt,
			}

			// Parse responses JSON
			var responsesData interface{}
			if err := json.Unmarshal([]byte(response.Responses), &responsesData); err == nil {
				responseData["responses"] = responsesData
			}

			// Add survey details
			if response.Survey != nil {
				responseData["survey"] = map[string]interface{}{
					"id":          response.Survey.ID,
					"title":       response.Survey.Title,
					"description": response.Survey.Description,
					"survey_type": response.Survey.SurveyType,
				}
			}

			responseResponses = append(responseResponses, responseData)
		}

		// Get statistics
		var stats struct {
			TotalResponses    int64   `gorm:"column:total_responses"`
			TotalXPEarned     int64   `gorm:"column:total_xp_earned"`
			AverageCompletion float64 `gorm:"column:average_completion"`
		}

		db.Raw(`
			SELECT 
				COUNT(*) as total_responses,
				COALESCE(SUM(xp_awarded), 0) as total_xp_earned,
				COALESCE(AVG(completion_percentage), 0) as average_completion
			FROM survey_responses 
			WHERE user_id = ?
		`, user.ID).Scan(&stats)

		response := map[string]interface{}{
			"responses": responseResponses,
			"stats": map[string]interface{}{
				"total_responses":    stats.TotalResponses,
				"total_xp_earned":    stats.TotalXPEarned,
				"average_completion": stats.AverageCompletion,
			},
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
