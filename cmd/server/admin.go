package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rohit21755/gg_server.git/internal/store"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Admin: Get all users
func adminGetUsersHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 50
		offset := 0

		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}
		if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
				offset = o
			}
		}

		var users []store.User
		if err := db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch users")
			return
		}

		var total int64
		db.Model(&store.User{}).Count(&total)

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"users": users,
			"total": total,
			"limit": limit,
			"offset": offset,
		})
	}
}

// Admin: Get user by ID
func adminGetUserHandler(db *gorm.DB) http.HandlerFunc {
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

		writeJSON(w, http.StatusOK, user)
	}
}

// Admin: Create user
func adminCreateUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Role      string `json:"role"`
			CollegeID *int   `json:"college_id"`
			StateID   *int   `json:"state_id"`
		}

		if err := readJSON(w, r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to hash password")
			return
		}

		// Generate unique referral code
		var referralCode string
		for {
			referralCode = generateReferralCode()
			_, err := store.GetUserByReferralCode(db, referralCode)
			if err != nil && err == gorm.ErrRecordNotFound {
				break
			}
		}

		user := &store.User{
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			Role:         req.Role,
			CollegeID:    req.CollegeID,
			StateID:      req.StateID,
			ReferralCode: referralCode,
		}

		if err := store.CreateUser(db, user); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to create user")
			return
		}

		writeJSON(w, http.StatusCreated, user)
	}
}

// Admin: Update user
func adminUpdateUserHandler(db *gorm.DB) http.HandlerFunc {
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

		var req struct {
			Email     *string `json:"email"`
			FirstName *string `json:"first_name"`
			LastName  *string `json:"last_name"`
			Role      *string `json:"role"`
			CollegeID *int   `json:"college_id"`
			StateID   *int   `json:"state_id"`
			IsActive  *bool  `json:"is_active"`
		}

		if err := readJSON(w, r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Email != nil {
			user.Email = *req.Email
		}
		if req.FirstName != nil {
			user.FirstName = *req.FirstName
		}
		if req.LastName != nil {
			user.LastName = *req.LastName
		}
		if req.Role != nil {
			user.Role = *req.Role
		}
		if req.CollegeID != nil {
			user.CollegeID = req.CollegeID
		}
		if req.StateID != nil {
			user.StateID = req.StateID
		}
		if req.IsActive != nil {
			user.IsActive = *req.IsActive
		}

		if err := store.UpdateUser(db, user); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to update user")
			return
		}

		writeJSON(w, http.StatusOK, user)
	}
}

// Admin: Delete user
func adminDeleteUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := chi.URLParam(r, "id")
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid user ID")
			return
		}

		if err := store.DeleteUser(db, uint(userID)); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to delete user")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"message": "user deleted"})
	}
}

// Admin: Reset user password
func adminResetPasswordHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := chi.URLParam(r, "id")
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid user ID")
			return
		}

		var req struct {
			NewPassword string `json:"new_password"`
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

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to hash password")
			return
		}

		user.PasswordHash = string(hashedPassword)
		if err := store.UpdateUser(db, user); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to update password")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"message": "password reset successfully"})
	}
}

// Admin: Get pending submissions
func adminGetPendingSubmissionsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var submissions []store.Submission
		if err := db.Where("status = ?", "pending").
			Order("created_at ASC").
			Find(&submissions).Error; err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch submissions")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"submissions": submissions,
		})
	}
}

// Admin: Review submission
func adminReviewSubmissionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		submissionIDStr := chi.URLParam(r, "id")
		submissionID, err := strconv.ParseUint(submissionIDStr, 10, 32)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid submission ID")
			return
		}

		var req struct {
			Status  string `json:"status"` // approved, rejected
			Comment string `json:"comment,omitempty"`
		}
		if err := readJSON(w, r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Status != "approved" && req.Status != "rejected" {
			writeJSONError(w, http.StatusBadRequest, "status must be 'approved' or 'rejected'")
			return
		}

		var submission store.Submission
		if err := db.First(&submission, submissionID).Error; err != nil {
			writeJSONError(w, http.StatusNotFound, "submission not found")
			return
		}

		submission.Status = req.Status
		if err := db.Save(&submission).Error; err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to update submission")
			return
		}

		// Update user stats if approved
		if req.Status == "approved" {
			db.Model(&store.User{}).Where("id = ?", submission.UserID).
				UpdateColumn("approved_submissions", gorm.Expr("approved_submissions + 1"))
		}

		writeJSON(w, http.StatusOK, submission)
	}
}

// Admin: Get submission statistics
func adminGetSubmissionStatsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pending, approved, rejected int64
		db.Model(&store.Submission{}).Where("status = ?", "pending").Count(&pending)
		db.Model(&store.Submission{}).Where("status = ?", "approved").Count(&approved)
		db.Model(&store.Submission{}).Where("status = ?", "rejected").Count(&rejected)

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"pending":  pending,
			"approved": approved,
			"rejected": rejected,
			"total":    pending + approved + rejected,
		})
	}
}

// Admin: Award XP
func adminAwardXPHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID      uint   `json:"user_id"`
			Amount      int    `json:"amount"`
			Description string `json:"description"`
		}

		if err := readJSON(w, r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		user, err := store.GetUserByID(db, req.UserID)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, "user not found")
			return
		}

		// Create XP transaction
		userIDInt := int(req.UserID)
		desc := req.Description
		tx := &store.XPTransaction{
			UserID:          &userIDInt,
			Amount:          req.Amount,
			TransactionType: "correction",
			BalanceAfter:    user.XP + req.Amount,
			Description:     &desc,
		}
		if err := store.CreateXPTransaction(db, tx); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to create transaction")
			return
		}

		// Update user XP
		user.XP += req.Amount
		if err := store.UpdateUser(db, user); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to update user XP")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"message": "XP awarded successfully",
			"user":    user,
		})
	}
}

// Admin: Penalize XP
func adminPenalizeXPHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID      uint   `json:"user_id"`
			Amount      int    `json:"amount"`
			Description string `json:"description"`
		}

		if err := readJSON(w, r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		user, err := store.GetUserByID(db, req.UserID)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, "user not found")
			return
		}

		// Create XP transaction
		userIDInt := int(req.UserID)
		desc := req.Description
		newXP := user.XP - req.Amount
		if newXP < 0 {
			newXP = 0
		}
		tx := &store.XPTransaction{
			UserID:          &userIDInt,
			Amount:          -req.Amount,
			TransactionType: "correction",
			BalanceAfter:    newXP,
			Description:     &desc,
		}
		if err := store.CreateXPTransaction(db, tx); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to create transaction")
			return
		}

		// Update user XP (ensure it doesn't go below 0)
		newXPValue := user.XP - req.Amount
		if newXPValue < 0 {
			newXPValue = 0
		}
		user.XP = newXPValue
		if err := store.UpdateUser(db, user); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to update user XP")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"message": "XP penalized successfully",
			"user":    user,
		})
	}
}

// Admin: Award badge
func adminAwardBadgeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID uint `json:"user_id"`
			BadgeID uint `json:"badge_id"`
		}

		if err := readJSON(w, r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		userIDInt := int(req.UserID)
		badgeIDInt := int(req.BadgeID)
		userBadge := &store.UserBadge{
			UserID:  userIDInt,
			BadgeID: badgeIDInt,
		}

		if err := store.CreateUserBadge(db, userBadge); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to award badge")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"message": "badge awarded successfully",
			"badge":   userBadge,
		})
	}
}

// Admin: Get dashboard overview
func adminDashboardHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var totalUsers, activeUsers, totalTasks, totalSubmissions int64
		db.Model(&store.User{}).Count(&totalUsers)
		db.Model(&store.User{}).Where("is_active = ?", true).Count(&activeUsers)
		db.Model(&store.Task{}).Count(&totalTasks)
		db.Model(&store.Submission{}).Count(&totalSubmissions)

		var pendingSubmissions int64
		db.Model(&store.Submission{}).Where("status = ?", "pending").Count(&pendingSubmissions)

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"total_users":        totalUsers,
			"active_users":       activeUsers,
			"total_tasks":        totalTasks,
			"total_submissions":  totalSubmissions,
			"pending_submissions": pendingSubmissions,
		})
	}
}

