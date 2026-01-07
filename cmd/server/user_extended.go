package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// Search users by name/email/college
func searchUsersHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			writeJSONError(w, http.StatusBadRequest, "query parameter 'q' is required")
			return
		}

		var users []store.User
		searchPattern := "%" + query + "%"
		if err := db.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
			searchPattern, searchPattern, searchPattern).
			Where("is_active = ?", true).
			Limit(50).
			Find(&users).Error; err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to search users")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"users": users,
		})
	}
}

// Block/unblock user (admin)
func blockUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := chi.URLParam(r, "id")
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid user ID")
			return
		}

		var req struct {
			Blocked bool `json:"blocked"`
		}
		if err := readJSON(w, r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		user, err := store.GetUserByID(db, uint(userID))
		if err != nil {
			writeJSONError(w, http.StatusNotFound, "user not found")
			return
		}

		user.IsActive = !req.Blocked
		if err := store.UpdateUser(db, user); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to update user")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"message": "user status updated",
			"user":    user,
		})
	}
}

// Get detailed user statistics
func getUserStatsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := chi.URLParam(r, "id")
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid user ID")
			return
		}

		user, err := store.GetUserByID(db, uint(userID))
		if err != nil {
			writeJSONError(w, http.StatusNotFound, "user not found")
			return
		}

		// Get badge count
		badgeCount, _ := store.GetUserBadgeCount(db, uint(userID))

		// Get wallet balance
		wallet, _ := store.GetUserWallet(db, uint(userID))

		// Get submission stats
		var totalSubmissions, approvedSubmissions int64
		db.Model(&store.Submission{}).Where("user_id = ?", userID).Count(&totalSubmissions)
		db.Model(&store.Submission{}).Where("user_id = ? AND status = ?", userID, "approved").Count(&approvedSubmissions)

		stats := map[string]interface{}{
			"user":                 user,
			"badge_count":          badgeCount,
			"wallet":               wallet,
			"total_submissions":    totalSubmissions,
			"approved_submissions": approvedSubmissions,
			"approval_rate":        0.0,
		}

		if totalSubmissions > 0 {
			stats["approval_rate"] = float64(approvedSubmissions) / float64(totalSubmissions) * 100
		}

		writeJSON(w, http.StatusOK, stats)
	}
}

// Get global leaderboard
func getGlobalLeaderboardHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 100
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
				limit = l
			}
		}

		var users []store.User
		if err := db.Where("is_active = ?", true).
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

// Get user activity feed
func getUserActivityHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		limit := 50
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		activities, err := store.GetUserActivityLogs(db, uint(user.ID), limit)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch activities")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"activities": activities,
		})
	}
}

// Get user dashboard stats
func getUserDashboardStatsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		// Get various stats
		badgeCount, _ := store.GetUserBadgeCount(db, uint(user.ID))
		wallet, _ := store.GetUserWallet(db, uint(user.ID))

		var pendingTasks, completedTasks int64
		userIDInt := int(user.ID)
		db.Model(&store.TaskAssignment{}).Where("assignee_id = ? AND status = ?", userIDInt, "assigned").Count(&pendingTasks)
		db.Model(&store.TaskAssignment{}).Where("assignee_id = ? AND status = ?", userIDInt, "completed").Count(&completedTasks)

		var unreadNotifications int64
		db.Model(&store.Notification{}).Where("user_id = ? AND is_read = ?", user.ID, false).Count(&unreadNotifications)

		stats := map[string]interface{}{
			"xp":                   user.XP,
			"level":                user.LevelID,
			"streak":               user.StreakCount,
			"badge_count":          badgeCount,
			"wallet":               wallet,
			"pending_tasks":        pendingTasks,
			"completed_tasks":      completedTasks,
			"unread_notifications": unreadNotifications,
			"total_submissions":    user.TotalSubmissions,
			"approved_submissions": user.ApprovedSubmissions,
		}

		writeJSON(w, http.StatusOK, stats)
	}
}
