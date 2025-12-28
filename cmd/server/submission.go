package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// Submission Response
type SubmissionResponse struct {
	ID             uint       `json:"id"`
	TaskID         uint       `json:"task_id"`
	TaskTitle      string     `json:"task_title"`
	CampaignID     uint       `json:"campaign_id"`
	CampaignName   string     `json:"campaign_name"`
	ProofType      string     `json:"proof_type"`
	ProofURL       string     `json:"proof_url"`
	Status         string     `json:"status"`
	Score          *float64   `json:"score"`
	XPAwarded      int        `json:"xp_awarded"`
	ReviewComments string     `json:"review_comments"`
	SubmittedAt    time.Time  `json:"submitted_at"`
	ReviewedAt     *time.Time `json:"reviewed_at"`
}

// Get User Submissions
func getSubmissionsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Parse query parameters
		status := r.URL.Query().Get("status")
		campaignID, _ := strconv.ParseUint(r.URL.Query().Get("campaign_id"), 10, 32)
		taskID, _ := strconv.ParseUint(r.URL.Query().Get("task_id"), 10, 32)
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
		query := db.Model(&store.Submission{}).Where("user_id = ?", user.ID)

		if status != "" {
			query = query.Where("status = ?", status)
		}
		if campaignID > 0 {
			query = query.Where("campaign_id = ?", campaignID)
		}
		if taskID > 0 {
			query = query.Where("task_id = ?", taskID)
		}

		// Get total count
		var totalCount int64
		query.Count(&totalCount)

		// Get submissions
		var submissions []store.Submission
		result := query.
			Preload("Task").
			Preload("Campaign").
			Order("submitted_at DESC").
			Offset(offset).
			Limit(limit).
			Find(&submissions)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Build response
		submissionResponses := make([]SubmissionResponse, len(submissions))
		for i, sub := range submissions {
			resp := SubmissionResponse{
				ID:             sub.ID,
				TaskID:         uint(*sub.TaskID),
				CampaignID:     uint(*sub.CampaignID),
				ProofType:      sub.ProofType,
				ProofURL:       sub.ProofURL,
				Status:         sub.Status,
				Score:          sub.Score,
				XPAwarded:      sub.XPAwarded,
				ReviewComments: *sub.ReviewComments,
				SubmittedAt:    sub.SubmittedAt,
				ReviewedAt:     sub.ReviewedAt,
			}

			if sub.Task != nil {
				resp.TaskTitle = sub.Task.Title
			}
			if sub.Campaign != nil {
				resp.CampaignName = sub.Campaign.Title
			}

			submissionResponses[i] = resp
		}

		response := map[string]interface{}{
			"submissions": submissionResponses,
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

// Create Submission
func createSubmissionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		var req SubmissionRequest
		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Validate request
		if err := Validate.Struct(req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Get task
		task, err := store.GetTaskByID(db, req.TaskID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("task not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Check if task is active
		if !task.IsActive {
			badRequestResponse(w, r, errors.New("task is not active"))
			return
		}

		// Check if user has already submitted
		existingSubmissions, err := store.GetSubmissionsByUserAndTask(db, user.ID, req.TaskID)
		if err == nil && len(existingSubmissions) > 0 {
			// Check if any submission is pending or approved
			for _, sub := range existingSubmissions {
				if sub.Status == "pending" || sub.Status == "under_review" || sub.Status == "approved" {
					conflictResponse(w, r, errors.New("you have already submitted for this task"))
					return
				}
			}
		}

		// Check if task has duration limit
		if task.DurationHours != nil && *task.DurationHours > 0 {
			deadline := task.CreatedAt.Add(time.Duration(*task.DurationHours) * time.Hour)
			if time.Now().After(deadline) {
				badRequestResponse(w, r, errors.New("task submission deadline has passed"))
				return
			}
		}

		// Create submission
		taskID := int(req.TaskID)  // convert uint -> int
		userID := int(user.ID)     // convert uint -> int
		proofText := req.ProofText // already string

		submission := &store.Submission{
			TaskID:      &taskID,         // now *int
			UserID:      &userID,         // now *int (if struct expects pointer)
			CampaignID:  task.CampaignID, // check if needs conversion too
			ProofType:   req.ProofType,
			ProofURL:    req.ProofURL,
			ProofText:   &proofText, // now *string
			Status:      "pending",
			SubmittedAt: time.Now(),
		}

		if err := store.CreateSubmission(db, submission); err != nil {
			internalServerError(w, r, err)
			return
		}

		// Update user submission count
		user.TotalSubmissions++
		if err := store.UpdateUser(db, user); err != nil {
			internalServerError(w, r, err)
			return
		}

		response := map[string]interface{}{
			"submission": submission,
			"message":    "Submission created successfully",
		}

		if err := jsonResponse(w, http.StatusCreated, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Submission Details
func getSubmissionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		submissionIDStr := chi.URLParam(r, "id")
		submissionID, err := strconv.ParseUint(submissionIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid submission ID"))
			return
		}

		submission, err := store.GetSubmissionByID(db, uint(submissionID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("submission not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Check if user owns the submission or is admin
		if submission.UserID != intPtr(int(user.ID)) && user.Role != "admin" && user.Role != "state_lead" {
			unauthorizedResponse(w, r, errors.New("not authorized to view this submission"))
			return
		}

		// Load related data
		if err := db.Preload("Task").Preload("Campaign").First(submission).Error; err != nil {
			internalServerError(w, r, err)
			return
		}

		if err := jsonResponse(w, http.StatusOK, submission); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Update Submission
func updateSubmissionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		submissionIDStr := chi.URLParam(r, "id")
		submissionID, err := strconv.ParseUint(submissionIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid submission ID"))
			return
		}

		submission, err := store.GetSubmissionByID(db, uint(submissionID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("submission not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Check if user owns the submission
		if submission.UserID != intPtr(int(user.ID)) {
			unauthorizedResponse(w, r, errors.New("not authorized to update this submission"))
			return
		}

		// Check if submission can be updated
		if submission.Status != "draft" && submission.Status != "needs_revision" {
			badRequestResponse(w, r, errors.New("submission cannot be updated in its current state"))
			return
		}

		var req struct {
			ProofURL  string `json:"proof_url"`
			ProofText string `json:"proof_text"`
		}

		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Update fields
		if req.ProofURL != "" {
			submission.ProofURL = req.ProofURL
		}
		if req.ProofText != "" {
			submission.ProofText = &req.ProofText
		}

		submission.Status = "pending"
		submission.SubmissionStage = "resubmission"
		submission.SubmittedAt = time.Now()
		submission.UpdatedAt = time.Now()

		if err := store.UpdateSubmission(db, submission); err != nil {
			internalServerError(w, r, err)
			return
		}

		response := map[string]interface{}{
			"submission": submission,
			"message":    "Submission updated successfully",
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Delete Submission
func deleteSubmissionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		submissionIDStr := chi.URLParam(r, "id")
		submissionID, err := strconv.ParseUint(submissionIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid submission ID"))
			return
		}

		submission, err := store.GetSubmissionByID(db, uint(submissionID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("submission not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Check if user owns the submission
		if submission.UserID != intPtr(int(user.ID)) {
			unauthorizedResponse(w, r, errors.New("not authorized to delete this submission"))
			return
		}

		// Check if submission can be deleted
		if submission.Status != "draft" {
			badRequestResponse(w, r, errors.New("only draft submissions can be deleted"))
			return
		}

		if err := store.DeleteSubmission(db, submission.ID); err != nil {
			internalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Get Submission Proof
func getSubmissionProofHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		submissionIDStr := chi.URLParam(r, "id")
		submissionID, err := strconv.ParseUint(submissionIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid submission ID"))
			return
		}

		submission, err := store.GetSubmissionByID(db, uint(submissionID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("submission not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Check if user owns the submission or is admin/reviewer
		if submission.UserID != intPtr(int(user.ID)) && user.Role != "admin" && user.Role != "state_lead" {
			unauthorizedResponse(w, r, errors.New("not authorized to view this proof"))
			return
		}

		response := map[string]interface{}{
			"proof_url":  submission.ProofURL,
			"proof_text": submission.ProofText,
			"proof_type": submission.ProofType,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}
