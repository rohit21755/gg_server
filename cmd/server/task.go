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

// Task List Response
type TaskResponse struct {
	ID           uint                `json:"id"`
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	CampaignID   uint                `json:"campaign_id"`
	CampaignName string              `json:"campaign_name"`
	TaskType     string              `json:"task_type"`
	ProofType    string              `json:"proof_type"`
	XPReward     int                 `json:"xp_reward"`
	CoinReward   int                 `json:"coin_reward"`
	Priority     string              `json:"priority"`
	Duration     int                 `json:"duration_hours"`
	IsActive     bool                `json:"is_active"`
	CreatedAt    time.Time           `json:"created_at"`
	Deadline     time.Time           `json:"deadline"`
	MySubmission *SubmissionResponse `json:"my_submission,omitempty"`
}

// Submission Request
type SubmissionRequest struct {
	TaskID    uint   `json:"task_id" validate:"required"`
	ProofType string `json:"proof_type" validate:"required"`
	ProofURL  string `json:"proof_url" validate:"required,url"`
	ProofText string `json:"proof_text"`
}

// Get Tasks with Pagination
func getTasksHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Parse query parameters
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit

		// Filter parameters
		campaignID, _ := strconv.ParseUint(r.URL.Query().Get("campaign_id"), 10, 32)
		taskType := r.URL.Query().Get("type")
		priority := r.URL.Query().Get("priority")
		isActive := r.URL.Query().Get("is_active")

		// Build query
		query := db.Model(&store.Task{})

		if campaignID > 0 {
			query = query.Where("campaign_id = ?", campaignID)
		}
		if taskType != "" {
			query = query.Where("task_type = ?", taskType)
		}
		if priority != "" {
			query = query.Where("priority = ?", priority)
		}
		if isActive != "" {
			query = query.Where("is_active = ?", isActive == "true")
		}

		// Get total count
		var totalCount int64
		query.Count(&totalCount)

		// Get tasks
		var tasks []store.Task
		result := query.
			Preload("Campaign").
			Order("created_at DESC").
			Offset(offset).
			Limit(limit).
			Find(&tasks)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Get user's submissions for these tasks
		taskIDs := make([]uint, len(tasks))
		for i, task := range tasks {
			taskIDs[i] = task.ID
		}

		submissions, _ := store.GetSubmissionsByUserAndTasks(db, user.ID, taskIDs)
		submissionMap := make(map[uint]store.Submission)
		for _, sub := range submissions {
			submissionMap[uint(*sub.TaskID)] = sub
		}

		// Build response
		taskResponses := make([]TaskResponse, len(tasks))
		for i, task := range tasks {
			taskResp := TaskResponse{
				ID:          task.ID,
				Title:       task.Title,
				Description: task.Description,
				CampaignID:  uint(*task.CampaignID),
				TaskType:    task.TaskType,
				ProofType:   task.ProofType,
				XPReward:    task.XPReward,
				CoinReward:  task.CoinReward,
				Priority:    task.Priority,
				Duration:    *task.DurationHours,
				IsActive:    task.IsActive,
				CreatedAt:   task.CreatedAt,
			}

			if task.Campaign != nil {
				taskResp.CampaignName = task.Campaign.Title
				// Calculate deadline
				taskResp.Deadline = task.CreatedAt.Add(time.Duration(*task.DurationHours) * time.Hour)
				if taskResp.Deadline.After(task.Campaign.EndDate) {
					taskResp.Deadline = task.Campaign.EndDate
				}
			}

			// Add user's submission if exists
			if sub, exists := submissionMap[task.ID]; exists {
				submissionResp := SubmissionResponse{
					ID:          sub.ID,
					Status:      sub.Status,
					Score:       sub.Score,
					SubmittedAt: sub.SubmittedAt,
				}
				taskResp.MySubmission = &submissionResp
			}

			taskResponses[i] = taskResp
		}

		response := map[string]interface{}{
			"tasks": taskResponses,
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

// Get Task Details
func getTaskHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskIDStr := chi.URLParam(r, "id")
		taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid task ID"))
			return
		}

		task, err := store.GetTaskByID(db, uint(taskID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("task not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Get campaign details
		var campaign store.Campaign
		if err := db.First(&campaign, task.CampaignID).Error; err == nil {
			task.Campaign = &campaign
		}

		// Get assignment details
		var assignments []store.TaskAssignment
		db.Where("task_id = ?", task.ID).Find(&assignments)

		if err := jsonResponse(w, http.StatusOK, task); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Assigned Tasks
func getAssignedTasksHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Get tasks assigned to user
		assignments, err := store.GetTaskAssignmentsByUser(db, user.ID)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		// Extract task IDs
		taskIDs := make([]uint, len(assignments))
		for i, assignment := range assignments {
			taskIDs[i] = uint(*assignment.TaskID)
		}

		// Get tasks
		var tasks []store.Task
		if err := db.Where("id IN ?", taskIDs).
			Preload("Campaign").
			Find(&tasks).Error; err != nil {
			internalServerError(w, r, err)
			return
		}

		// Get submissions
		submissions, _ := store.GetSubmissionsByUserAndTasks(db, user.ID, taskIDs)
		submissionMap := make(map[uint]store.Submission)
		for _, sub := range submissions {
			submissionMap[uint(*sub.TaskID)] = sub
		}

		// Build response
		taskResponses := make([]TaskResponse, len(tasks))
		for i, task := range tasks {
			taskResp := TaskResponse{
				ID:          task.ID,
				Title:       task.Title,
				Description: task.Description,
				CampaignID:  uint(*task.CampaignID),
				TaskType:    task.TaskType,
				ProofType:   task.ProofType,
				XPReward:    task.XPReward,
				CoinReward:  task.CoinReward,
				Priority:    task.Priority,
				Duration:    *task.DurationHours,
				IsActive:    task.IsActive,
				CreatedAt:   task.CreatedAt,
			}

			if task.Campaign != nil {
				taskResp.CampaignName = task.Campaign.Title
				taskResp.Deadline = task.CreatedAt.Add(time.Duration(*task.DurationHours) * time.Hour)
			}

			if sub, exists := submissionMap[task.ID]; exists {
				submissionResp := SubmissionResponse{
					ID:          sub.ID,
					Status:      sub.Status,
					Score:       sub.Score,
					SubmittedAt: sub.SubmittedAt,
				}
				taskResp.MySubmission = &submissionResp
			}

			taskResponses[i] = taskResp
		}

		if err := jsonResponse(w, http.StatusOK, taskResponses); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Available Tasks
func getAvailableTasksHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Get tasks that are active and not assigned to user
		var tasks []store.Task
		query := db.Model(&store.Task{}).
			Where("is_active = ?", true).
			Where("id NOT IN (SELECT task_id FROM task_assignments WHERE assignee_id = ? AND assignee_type = 'user')", user.ID).
			Preload("Campaign").
			Order("priority DESC, created_at DESC")

		// Filter by campaign if specified
		if campaignID := r.URL.Query().Get("campaign_id"); campaignID != "" {
			query = query.Where("campaign_id = ?", campaignID)
		}

		if err := query.Find(&tasks).Error; err != nil {
			internalServerError(w, r, err)
			return
		}

		if err := jsonResponse(w, http.StatusOK, tasks); err != nil {
			internalServerError(w, r, err)
		}
	}
}
