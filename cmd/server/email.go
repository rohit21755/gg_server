package main

import (
	"net/http"

	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// Get email preferences
func getEmailPreferencesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		prefs, err := store.GetUserEmailPreferences(db, uint(user.ID))
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch preferences")
			return
		}

		writeJSON(w, http.StatusOK, prefs)
	}
}

// Update email preferences
func updateEmailPreferencesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		var req struct {
			MarketingEmails   *bool `json:"marketing_emails"`
			TaskNotifications *bool `json:"task_notifications"`
			AchievementEmails  *bool `json:"achievement_emails"`
			WeeklyDigest      *bool `json:"weekly_digest"`
		}

		if err := readJSON(w, r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		prefs, err := store.GetUserEmailPreferences(db, uint(user.ID))
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch preferences")
			return
		}

		if req.MarketingEmails != nil {
			prefs.MarketingEmails = *req.MarketingEmails
		}
		if req.TaskNotifications != nil {
			prefs.TaskNotifications = *req.TaskNotifications
		}
		if req.AchievementEmails != nil {
			prefs.AchievementEmails = *req.AchievementEmails
		}
		if req.WeeklyDigest != nil {
			prefs.WeeklyDigest = *req.WeeklyDigest
		}

		if err := store.UpdateUserEmailPreferences(db, prefs); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to update preferences")
			return
		}

		writeJSON(w, http.StatusOK, prefs)
	}
}

// Resend verification email
func resendVerificationEmailHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := GetUserFromContext(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		// TODO: Implement email sending logic
		// For now, just return success
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "verification email sent",
		})
	}
}

