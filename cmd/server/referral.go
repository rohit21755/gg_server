package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// Get Referrals
func getReferralsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Parse query parameters
		status := r.URL.Query().Get("status")
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
		query := db.Model(&store.Referral{}).Where("referrer_id = ?", user.ID)

		if status != "" {
			query = query.Where("status = ?", status)
		}

		// Get total count
		var totalCount int64
		query.Count(&totalCount)

		// Get referrals
		var referrals []store.Referral
		result := query.
			Preload("ReferredUser").
			Order("created_at DESC").
			Offset(offset).
			Limit(limit).
			Find(&referrals)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Format response
		var responseReferrals []map[string]interface{}
		for _, referral := range referrals {
			referralData := map[string]interface{}{
				"id":                     referral.ID,
				"referred_email":         referral.ReferredEmail,
				"status":                 referral.Status,
				"xp_awarded":             referral.XPAwarded,
				"xp_awarded_to_referred": referral.XPAwardedToReferred,
				"conversion_stage":       referral.ConversionStage,
				"created_at":             referral.CreatedAt,
				"updated_at":             referral.UpdatedAt,
			}

			// Add referred user info if exists
			if referral.ReferredUser != nil {
				referralData["referred_user"] = map[string]interface{}{
					"id":         referral.ReferredUser.ID,
					"first_name": referral.ReferredUser.FirstName,
					"last_name":  referral.ReferredUser.LastName,
					"email":      referral.ReferredUser.Email,
				}
			}

			responseReferrals = append(responseReferrals, referralData)
		}

		// Get statistics
		var stats struct {
			TotalReferrals     int64 `gorm:"column:total_referrals"`
			PendingReferrals   int64 `gorm:"column:pending_referrals"`
			JoinedReferrals    int64 `gorm:"column:joined_referrals"`
			ConvertedReferrals int64 `gorm:"column:converted_referrals"`
			TotalXPAwarded     int64 `gorm:"column:total_xp_awarded"`
		}

		db.Raw(`
			SELECT 
				COUNT(*) as total_referrals,
				COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_referrals,
				COUNT(CASE WHEN status = 'joined' THEN 1 END) as joined_referrals,
				COUNT(CASE WHEN status = 'converted' THEN 1 END) as converted_referrals,
				COALESCE(SUM(xp_awarded), 0) as total_xp_awarded
			FROM referrals 
			WHERE referrer_id = ?
		`, user.ID).Scan(&stats)

		response := map[string]interface{}{
			"referrals": responseReferrals,
			"stats": map[string]interface{}{
				"total_referrals":     stats.TotalReferrals,
				"pending_referrals":   stats.PendingReferrals,
				"joined_referrals":    stats.JoinedReferrals,
				"converted_referrals": stats.ConvertedReferrals,
				"total_xp_awarded":    stats.TotalXPAwarded,
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

// Get Referral Code
func getReferralCodeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Get user details to get referral code
		dbUser, err := store.GetUserByID(db, user.ID)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		response := map[string]interface{}{
			"referral_code": dbUser.ReferralCode,
			"referral_url":  fmt.Sprintf("https://app.example.com/register?ref=%s", dbUser.ReferralCode),
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Referral Invites
func getReferralInvitesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Parse query parameters
		status := r.URL.Query().Get("status")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit

		// Build query for invites (pending or joined status)
		query := db.Model(&store.Referral{}).Where("referrer_id = ?", user.ID)

		if status != "" {
			query = query.Where("status = ?", status)
		} else {
			// Default to pending and joined
			query = query.Where("status IN (?, ?)", "pending", "joined")
		}

		// Get total count
		var totalCount int64
		query.Count(&totalCount)

		// Get invites
		var referrals []store.Referral
		result := query.
			Preload("ReferredUser").
			Order("created_at DESC").
			Offset(offset).
			Limit(limit).
			Find(&referrals)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Format response
		var responseInvites []map[string]interface{}
		for _, referral := range referrals {
			inviteData := map[string]interface{}{
				"id":             referral.ID,
				"referred_email": referral.ReferredEmail,
				"status":         referral.Status,
				"created_at":     referral.CreatedAt,
				"updated_at":     referral.UpdatedAt,
			}

			// Add referred user info if exists (status is joined or beyond)
			if referral.ReferredUser != nil {
				inviteData["referred_user"] = map[string]interface{}{
					"id":         referral.ReferredUser.ID,
					"first_name": referral.ReferredUser.FirstName,
					"last_name":  referral.ReferredUser.LastName,
					"email":      referral.ReferredUser.Email,
				}
			}

			responseInvites = append(responseInvites, inviteData)
		}

		response := map[string]interface{}{
			"invites": responseInvites,
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

// Send Referral Invite
func sendReferralInviteHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		var req struct {
			Email string `json:"email" validate:"required,email"`
		}

		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Validate email
		if req.Email == "" {
			badRequestResponse(w, r, errors.New("email is required"))
			return
		}

		// Check if user is trying to refer themselves
		if req.Email == user.Email {
			badRequestResponse(w, r, errors.New("cannot refer yourself"))
			return
		}

		// Check if referral already exists
		var existingReferral store.Referral
		result := db.Where("referrer_id = ? AND referred_email = ?", user.ID, req.Email).First(&existingReferral)
		if result.Error == nil {
			// Referral already exists
			conflictResponse(w, r, errors.New("referral invite already sent to this email"))
			return
		}
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			internalServerError(w, r, result.Error)
			return
		}

		// Check if email already belongs to a user
		var existingUser store.User
		userCheck := db.Where("email = ?", req.Email).First(&existingUser)
		if userCheck.Error == nil {
			// User already exists, check if they were already referred by this user
			existingReferralCheck := db.Where("referrer_id = ? AND referred_user_id = ?", user.ID, existingUser.ID).First(&existingReferral)
			if existingReferralCheck.Error == nil {
				conflictResponse(w, r, errors.New("this user has already been referred by you"))
				return
			}

			// User exists but not referred by this user - could create referral with joined status
			referrerID := uint(user.ID)
			referredUserID := existingUser.ID
			referral := &store.Referral{
				ReferrerID:     &referrerID,
				ReferredEmail:  req.Email,
				ReferredUserID: &referredUserID,
				Status:         "joined",
			}

			if err := store.CreateReferral(db, referral); err != nil {
				internalServerError(w, r, err)
				return
			}

			response := map[string]interface{}{
				"message":        "Referral created successfully (user already exists)",
				"referral_id":    referral.ID,
				"status":         referral.Status,
				"referred_email": referral.ReferredEmail,
			}

			if err := jsonResponse(w, http.StatusOK, response); err != nil {
				internalServerError(w, r, err)
			}
			return
		}

		// Create referral invite with pending status
		referrerID := uint(user.ID)
		referral := &store.Referral{
			ReferrerID:    &referrerID,
			ReferredEmail: req.Email,
			Status:        "pending",
		}

		if err := store.CreateReferral(db, referral); err != nil {
			internalServerError(w, r, err)
			return
		}

		// Get user details for referral URL
		dbUser, err := store.GetUserByID(db, user.ID)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		// TODO: Send email invitation (integrate with email service)
		// For now, just return success response

		response := map[string]interface{}{
			"message":        "Referral invite sent successfully",
			"referral_id":    referral.ID,
			"status":         referral.Status,
			"referred_email": referral.ReferredEmail,
			"referral_url":   fmt.Sprintf("https://app.example.com/register?ref=%s", dbUser.ReferralCode),
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}
