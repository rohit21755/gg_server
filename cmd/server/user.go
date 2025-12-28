package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// Profile Update Request
type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	CollegeID uint   `json:"college_id"`
	StateID   uint   `json:"state_id"`
}

// Get Current User Profile
func getCurrentUserProfileHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Get user with related data
		fullUser, err := store.GetUserWithRelations(db, user.ID)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		// Get badges count
		badgeCount, err := store.GetUserBadgeCount(db, user.ID)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		// Get level info
		level, err := store.GetLevelByID(db, uint(*fullUser.LevelID))
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		// Get streak info
		streak, _ := store.GetUserStreak(db, user.ID, "daily_login")

		response := map[string]interface{}{
			"user": map[string]interface{}{
				"id":         fullUser.ID,
				"email":      fullUser.Email,
				"first_name": fullUser.FirstName,
				"last_name":  fullUser.LastName,
				"phone":      fullUser.Phone,
				"role":       fullUser.Role,
				"xp":         fullUser.XP,
				"college_id": fullUser.CollegeID,
				"state_id":   fullUser.StateID,
				"avatar_url": fullUser.AvatarURL,
				"resume_url": fullUser.ResumeURL,
				"created_at": fullUser.CreatedAt,
			},
			"stats": map[string]interface{}{
				"level":                level.Name,
				"badge_count":          badgeCount,
				"current_streak":       0,
				"longest_streak":       0,
				"total_submissions":    fullUser.TotalSubmissions,
				"approved_submissions": fullUser.ApprovedSubmissions,
				"win_rate":             fullUser.WinRate,
				"referral_code":        fullUser.ReferralCode,
			},
		}

		if streak != nil {
			response["stats"].(map[string]interface{})["current_streak"] = streak.CurrentStreak
			response["stats"].(map[string]interface{})["longest_streak"] = streak.LongestStreak
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Update User Profile
func updateUserProfileHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		var req UpdateProfileRequest
		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Update fields if provided
		if req.FirstName != "" {
			user.FirstName = req.FirstName
		}
		if req.LastName != "" {
			user.LastName = req.LastName
		}
		if req.Phone != "" {
			user.Phone = &req.Phone
		}
		if req.CollegeID > 0 {
			collegeID := int(req.CollegeID)
			user.CollegeID = &collegeID
		}
		if req.StateID > 0 {
			stateID := int(req.StateID)
			user.StateID = &stateID
		}

		if err := store.UpdateUser(db, user); err != nil {
			internalServerError(w, r, err)
			return
		}

		if err := jsonResponse(w, http.StatusOK, user); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Update Avatar
func updateAvatarHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			badRequestResponse(w, r, err)
			return
		}

		file, handler, err := r.FormFile("avatar")
		if err != nil {
			badRequestResponse(w, r, errors.New("avatar file is required"))
			return
		}
		defer file.Close()

		// Validate file type
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
		}
		if !allowedTypes[handler.Header.Get("Content-Type")] {
			badRequestResponse(w, r, errors.New("only JPEG, PNG, and GIF images are allowed"))
			return
		}

		// Upload to storage (S3, local, etc.)
		// For now, we'll just store the filename
		avatarURL := "/uploads/avatars/" + handler.Filename

		// In production, you would:
		// 1. Upload to S3/cloud storage
		// 2. Generate unique filename
		// 3. Create thumbnails

		user.AvatarURL = &avatarURL
		if err := store.UpdateUser(db, user); err != nil {
			internalServerError(w, r, err)
			return
		}

		response := map[string]string{
			"avatar_url": avatarURL,
			"message":    "Avatar updated successfully",
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Upload Resume
func uploadResumeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			badRequestResponse(w, r, err)
			return
		}

		file, handler, err := r.FormFile("resume")
		if err != nil {
			badRequestResponse(w, r, errors.New("resume file is required"))
			return
		}
		defer file.Close()

		// Validate file type
		allowedTypes := map[string]bool{
			"application/pdf":    true,
			"application/msword": true,
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		}
		if !allowedTypes[handler.Header.Get("Content-Type")] {
			badRequestResponse(w, r, errors.New("only PDF and DOC/DOCX files are allowed"))
			return
		}

		// Upload to storage
		resumeURL := "/uploads/resumes/" + handler.Filename

		user.ResumeURL = &resumeURL
		if err := store.UpdateUser(db, user); err != nil {
			internalServerError(w, r, err)
			return
		}

		response := map[string]string{
			"resume_url": resumeURL,
			"message":    "Resume uploaded successfully",
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get User Certificates
func getUserCertificatesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		certificates, err := store.GetUserCertificates(db, user.ID)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		if err := jsonResponse(w, http.StatusOK, certificates); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Download Certificate
func downloadCertificateHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		certificateIDStr := chi.URLParam(r, "id")
		certificateID, err := strconv.ParseUint(certificateIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid certificate ID"))
			return
		}

		certificate, err := store.GetCertificateByID(db, uint(certificateID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("certificate not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Check if certificate belongs to user
		if certificate.UserID != intPtr(int(user.ID)) {
			unauthorizedResponse(w, r, errors.New("not authorized to download this certificate"))
			return
		}

		// Serve file
		// In production, you would serve from S3 or storage
		http.ServeFile(w, r, certificate.CertificateURL)
	}
}
