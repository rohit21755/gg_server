package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// Get Rewards
func getRewardsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		category := r.URL.Query().Get("category")
		isFeatured := r.URL.Query().Get("is_featured")
		isActive := r.URL.Query().Get("is_active")
		rewardType := r.URL.Query().Get("type")

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
		query := db.Model(&store.RewardStore{})

		if category != "" {
			query = query.Where("category = ?", category)
		}
		if rewardType != "" {
			query = query.Where("reward_type = ?", rewardType)
		}
		if isFeatured != "" {
			query = query.Where("is_featured = ?", isFeatured == "true")
		}
		if isActive != "" {
			query = query.Where("is_active = ?", isActive == "true")
		} else {
			// Default to active rewards
			query = query.Where("is_active = ?", true)
		}

		// Get total count
		var totalCount int64
		query.Count(&totalCount)

		// Get rewards
		var rewards []store.RewardStore
		result := query.
			Order("is_featured DESC, xp_cost ASC, created_at DESC").
			Offset(offset).
			Limit(limit).
			Find(&rewards)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		response := map[string]interface{}{
			"rewards": rewards,
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

// Get Reward Details
func getRewardHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rewardIDStr := chi.URLParam(r, "id")
		rewardID, err := strconv.ParseUint(rewardIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid reward ID"))
			return
		}

		var reward store.RewardStore
		result := db.First(&reward, rewardID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("reward not found"))
			} else {
				internalServerError(w, r, result.Error)
			}
			return
		}

		// Check if reward is active
		if !reward.IsActive {
			badRequestResponse(w, r, errors.New("reward is not available"))
			return
		}

		// Add availability info
		available := reward.QuantityAvailable == nil || *reward.QuantityAvailable == 0 || *reward.QuantityAvailable > reward.QuantitySold
		rewardResponse := map[string]interface{}{
			"reward": reward,
			"availability": map[string]interface{}{
				"is_available": available,
				"remaining_quantity": func() int {
					if reward.QuantityAvailable == nil || *reward.QuantityAvailable == 0 {
						return -1 // Unlimited
					}
					return *reward.QuantityAvailable - reward.QuantitySold
				}(),
				"total_quantity": reward.QuantityAvailable,
				"sold_quantity":  reward.QuantitySold,
			},
		}

		if err := jsonResponse(w, http.StatusOK, rewardResponse); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Redeem Reward
func redeemRewardHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		rewardIDStr := chi.URLParam(r, "id")
		rewardID, err := strconv.ParseUint(rewardIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid reward ID"))
			return
		}

		// Get reward
		var reward store.RewardStore
		result := db.First(&reward, rewardID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("reward not found"))
			} else {
				internalServerError(w, r, result.Error)
			}
			return
		}

		// Validate reward
		if !reward.IsActive {
			badRequestResponse(w, r, errors.New("reward is not available"))
			return
		}

		// Check quantity
		if reward.QuantityAvailable != nil && *reward.QuantityAvailable > 0 && reward.QuantitySold >= *reward.QuantityAvailable {
			conflictResponse(w, r, errors.New("reward is out of stock"))
			return
		}

		// Check if user has enough XP
		dbUser, err := store.GetUserByID(db, user.ID)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		if dbUser.XP < reward.XPCost {
			badRequestResponse(w, r, errors.New("insufficient XP"))
			return
		}

		// Check if user has already redeemed this reward (if it's limited per user)
		// For now, we'll allow multiple redemptions unless specified in metadata
		var existingRedemptions int64
		db.Model(&store.UserReward{}).
			Where("user_id = ? AND reward_id = ?", user.ID, reward.ID).
			Count(&existingRedemptions)

		if existingRedemptions > 0 {
			// Check if reward allows multiple redemptions
			if reward.QuantityAvailable != nil && *reward.QuantityAvailable == 1 {
				conflictResponse(w, r, errors.New("you have already redeemed this reward"))
				return
			}
		}

		// Start transaction
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				internalServerError(w, r.(*http.Request), errors.New("transaction failed"))
			}
		}()

		// Create redemption record
		redemption := &store.UserReward{
			UserID:    intPtr(int(dbUser.ID)),
			RewardID:  intPtr(int(reward.ID)),
			Status:    "pending",
			XPPaid:    reward.XPCost,
			CoinsPaid: reward.CoinCost,
			CashPaid:  *reward.CashCost,
			ClaimedAt: time.Now(),
		}

		// Generate redemption code for digital rewards
		if reward.RewardType == "digital" || reward.RewardType == "gift_card" {
			redemption.RedemptionCode = stringPtr(generateRedemptionCode())
		}

		// Parse shipping address if required
		if reward.RewardType == "physical" {
			var req struct {
				ShippingAddress map[string]interface{} `json:"shipping_address"`
			}

			if err := readJSON(w, r, &req); err != nil {
				badRequestResponse(w, r, errors.New("shipping address is required for physical rewards"))
				return
			}

			// Validate shipping address
			if req.ShippingAddress == nil {
				badRequestResponse(w, r, errors.New("shipping address is required"))
				return
			}

			// Check required fields
			requiredFields := []string{"full_name", "address_line1", "city", "state", "country", "postal_code", "phone"}
			for _, field := range requiredFields {
				if _, ok := req.ShippingAddress[field]; !ok {
					badRequestResponse(w, r, errors.New("missing required field: "+field))
					return
				}
			}

			shippingAddrJSON, err := json.Marshal(req.ShippingAddress)
			if err != nil {
				badRequestResponse(w, r, errors.New("invalid shipping address format"))
				return
			}
			redemption.ShippingAddress = stringPtr(string(shippingAddrJSON))
		}

		if err := tx.Create(redemption).Error; err != nil {
			tx.Rollback()
			internalServerError(w, r, err)
			return
		}

		// Deduct XP from user
		dbUser.XP = dbUser.XP - reward.XPCost
		if err := tx.Save(dbUser).Error; err != nil {
			tx.Rollback()
			internalServerError(w, r, err)
			return
		}

		// Update reward quantity
		reward.QuantitySold++
		if err := tx.Save(&reward).Error; err != nil {
			tx.Rollback()
			internalServerError(w, r, err)
			return
		}

		// Create XP transaction
		metadataJSON, err := json.Marshal(map[string]interface{}{
			"reward_id":     reward.ID,
			"reward_name":   reward.Name,
			"redemption_id": redemption.ID,
		})
		if err != nil {
			tx.Rollback()
			internalServerError(w, r, err)
			return
		}
		xpTransaction := &store.XPTransaction{
			UserID:          intPtr(int(dbUser.ID)),
			TransactionType: "redemption",
			Amount:          -reward.XPCost,
			SourceType:      stringPtr("reward"),
			SourceID:        intPtr(int(redemption.ID)),
			Description:     stringPtr("Reward redemption: " + reward.Name),
			Metadata:        stringPtr(string(metadataJSON)),
		}

		if err := tx.Create(xpTransaction).Error; err != nil {
			tx.Rollback()
			internalServerError(w, r, err)
			return
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			internalServerError(w, r, err)
			return
		}

		// Create notification
		notificationDataJSON, err := json.Marshal(map[string]interface{}{
			"reward_id":       reward.ID,
			"reward_name":     reward.Name,
			"redemption_id":   redemption.ID,
			"redemption_code": redemption.RedemptionCode,
			"status":          redemption.Status,
		})
		if err != nil {
			internalServerError(w, r, err)
			return
		}
		notification := &store.Notification{
			UserID:           intPtr(int(dbUser.ID)),
			NotificationType: "reward_unlocked",
			Title:            "Reward Redeemed!",
			Message:          fmt.Sprintf("You have successfully redeemed: %s. Status: %s", reward.Name, redemption.Status),
			Data:             stringPtr(string(notificationDataJSON)),
			IsActionable:     true,
			ActionURL:        stringPtr("/rewards/redemptions/" + fmt.Sprintf("%d", redemption.ID)),
		}
		store.CreateNotification(db, notification)

		// Prepare response based on reward type
		responseData := map[string]interface{}{
			"redemption": map[string]interface{}{
				"id":          intPtr(int(redemption.ID)),
				"status":      redemption.Status,
				"reward_id":   reward.ID,
				"reward_name": reward.Name,
				"xp_paid":     redemption.XPPaid,
				"coins_paid":  redemption.CoinsPaid,
				"cash_paid":   redemption.CashPaid,
				"claimed_at":  redemption.ClaimedAt,
			},
			"user": map[string]interface{}{
				"remaining_xp": dbUser.XP,
			},
			"message": "Reward redeemed successfully",
		}

		// Add redemption code for digital rewards
		if redemption.RedemptionCode != nil && *redemption.RedemptionCode != "" {
			responseData["redemption"].(map[string]interface{})["redemption_code"] = redemption.RedemptionCode
			responseData["message"] = "Reward redeemed successfully. Your redemption code: " + *redemption.RedemptionCode
		}

		// Add shipping info for physical rewards
		if reward.RewardType == "physical" {
			responseData["redemption"].(map[string]interface{})["shipping_info"] = map[string]interface{}{
				"estimated_delivery": "7-14 business days",
				"tracking_available": false,
			}
		}

		if err := jsonResponse(w, http.StatusOK, responseData); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Reward Redemptions
func getRewardRedemptionsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Parse query parameters
		status := r.URL.Query().Get("status")
		rewardType := r.URL.Query().Get("type")

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
		query := db.Model(&store.UserReward{}).Where("user_id = ?", user.ID)

		if status != "" {
			query = query.Where("status = ?", status)
		}
		if rewardType != "" {
			// Join with rewards table to filter by reward type
			query = query.Joins("JOIN rewards_store ON rewards_store.id = user_rewards.reward_id").
				Where("rewards_store.reward_type = ?", rewardType)
		}

		// Get total count
		var totalCount int64
		query.Count(&totalCount)

		// Get redemptions with reward details
		var redemptions []store.UserReward
		result := query.
			Preload("Reward").
			Order("claimed_at DESC").
			Offset(offset).
			Limit(limit).
			Find(&redemptions)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Format response
		var responseRedemptions []map[string]interface{}
		for _, redemption := range redemptions {
			redemptionData := map[string]interface{}{
				"id":              redemption.ID,
				"status":          redemption.Status,
				"claimed_at":      redemption.ClaimedAt,
				"delivered_at":    redemption.DeliveredAt,
				"xp_paid":         redemption.XPPaid,
				"coins_paid":      redemption.CoinsPaid,
				"cash_paid":       redemption.CashPaid,
				"tracking_number": redemption.TrackingNumber,
				"redemption_code": redemption.RedemptionCode,
			}

			// Add reward details
			if redemption.Reward != nil {
				redemptionData["reward"] = map[string]interface{}{
					"id":          redemption.Reward.ID,
					"name":        redemption.Reward.Name,
					"description": redemption.Reward.Description,
					"reward_type": redemption.Reward.RewardType,
					"image_url":   redemption.Reward.ImageURL,
				}
			}

			// Add shipping address for physical rewards
			if redemption.ShippingAddress != nil {
				redemptionData["shipping_address"] = redemption.ShippingAddress
			}

			responseRedemptions = append(responseRedemptions, redemptionData)
		}

		// Get statistics
		var stats struct {
			TotalRedemptions  int64 `gorm:"column:total_redemptions"`
			TotalXPSpent      int64 `gorm:"column:total_xp_spent"`
			PendingDeliveries int64 `gorm:"column:pending_deliveries"`
		}

		db.Raw(`
			SELECT 
				COUNT(*) as total_redemptions,
				SUM(xp_paid) as total_xp_spent,
				COUNT(CASE WHEN status IN ('pending', 'processing', 'shipped') THEN 1 END) as pending_deliveries
			FROM user_rewards 
			WHERE user_id = ?
		`, user.ID).Scan(&stats)

		response := map[string]interface{}{
			"redemptions": responseRedemptions,
			"stats": map[string]interface{}{
				"total_redemptions":    stats.TotalRedemptions,
				"total_xp_spent":       stats.TotalXPSpent,
				"pending_deliveries":   stats.PendingDeliveries,
				"completed_deliveries": stats.TotalRedemptions - stats.PendingDeliveries,
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

// Get User's Specific Redemption Details
func getRedemptionDetailsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		redemptionIDStr := chi.URLParam(r, "id")
		redemptionID, err := strconv.ParseUint(redemptionIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid redemption ID"))
			return
		}

		var redemption store.UserReward
		result := db.Preload("Reward").First(&redemption, redemptionID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("redemption not found"))
			} else {
				internalServerError(w, r, result.Error)
			}
			return
		}

		// Check ownership
		if redemption.UserID == nil || *redemption.UserID != int(user.ID) {
			unauthorizedResponse(w, r, errors.New("not authorized to view this redemption"))
			return
		}

		// Format response
		response := map[string]interface{}{
			"id":               redemption.ID,
			"status":           redemption.Status,
			"claimed_at":       redemption.ClaimedAt,
			"delivered_at":     redemption.DeliveredAt,
			"xp_paid":          redemption.XPPaid,
			"coins_paid":       redemption.CoinsPaid,
			"cash_paid":        redemption.CashPaid,
			"tracking_number":  redemption.TrackingNumber,
			"redemption_code":  redemption.RedemptionCode,
			"shipping_address": redemption.ShippingAddress,
		}

		if redemption.Reward != nil {
			response["reward"] = redemption.Reward
		}

		// Add estimated delivery for physical rewards
		if redemption.Reward != nil && redemption.Reward.RewardType == "physical" {
			response["delivery_info"] = map[string]interface{}{
				"estimated_delivery": func() string {
					switch redemption.Status {
					case "pending":
						return "Processing - 3-5 business days"
					case "processing":
						return "Preparing for shipment - 2-3 business days"
					case "shipped":
						return "In transit - 5-7 business days"
					default:
						return "Contact support"
					}
				}(),
				"tracking_url": func() string {
					if redemption.TrackingNumber != nil && *redemption.TrackingNumber != "" {
						return "https://tracking.example.com/" + *redemption.TrackingNumber
					}
					return ""
				}(),
			}
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Cancel Redemption (if allowed)
func cancelRedemptionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		redemptionIDStr := chi.URLParam(r, "id")
		redemptionID, err := strconv.ParseUint(redemptionIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid redemption ID"))
			return
		}

		var redemption store.UserReward
		result := db.First(&redemption, redemptionID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("redemption not found"))
			} else {
				internalServerError(w, r, result.Error)
			}
			return
		}

		// Check ownership
		if redemption.UserID == nil || *redemption.UserID != int(user.ID) {
			unauthorizedResponse(w, r, errors.New("not authorized to cancel this redemption"))
			return
		}

		// Check if redemption can be cancelled
		if redemption.Status != "pending" && redemption.Status != "processing" {
			badRequestResponse(w, r, errors.New("redemption cannot be cancelled at this stage"))
			return
		}

		// Start transaction
		tx := db.Begin()

		// Update redemption status
		redemption.Status = "cancelled"
		if err := tx.Save(&redemption).Error; err != nil {
			tx.Rollback()
			internalServerError(w, r, err)
			return
		}

		// Refund XP to user
		dbUser, _ := store.GetUserByID(tx, user.ID)
		dbUser.XP += redemption.XPPaid
		if err := tx.Save(dbUser).Error; err != nil {
			tx.Rollback()
			internalServerError(w, r, err)
			return
		}

		// Create XP transaction for refund
		xpTransaction := &store.XPTransaction{
			UserID:          intPtr(int(dbUser.ID)),
			TransactionType: "redemption_refund",
			Amount:          redemption.XPPaid,
			SourceType:      stringPtr("reward"),
			SourceID:        intPtr(int(redemption.ID)),
			Description:     stringPtr("Redemption cancelled - XP refunded"),
		}

		if err := tx.Create(xpTransaction).Error; err != nil {
			tx.Rollback()
			internalServerError(w, r, err)
			return
		}

		// Update reward quantity
		var reward store.RewardStore
		tx.First(&reward, redemption.RewardID)
		if reward.ID > 0 {
			reward.QuantitySold--
			if err := tx.Save(&reward).Error; err != nil {
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

		// Create notification
		notificationDataJSON, err := json.Marshal(map[string]interface{}{
			"redemption_id":  redemption.ID,
			"xp_refunded":    redemption.XPPaid,
			"new_xp_balance": dbUser.XP,
		})
		if err != nil {
			internalServerError(w, r, err)
			return
		}
		notification := &store.Notification{
			UserID:           intPtr(int(dbUser.ID)),
			NotificationType: "system",
			Title:            "Redemption Cancelled",
			Message:          fmt.Sprintf("Your redemption has been cancelled. %d XP has been refunded to your account.", redemption.XPPaid),
			Data:             stringPtr(string(notificationDataJSON)),
		}
		store.CreateNotification(db, notification)

		response := map[string]interface{}{
			"message":        "Redemption cancelled successfully",
			"xp_refunded":    redemption.XPPaid,
			"new_xp_balance": dbUser.XP,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Admin: Update Redemption Status
func adminUpdateRedemptionStatusHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminUser, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Check if user is admin
		if adminUser.Role != "admin" && adminUser.Role != "state_lead" {
			unauthorizedResponse(w, r, errors.New("only admins can update redemption status"))
			return
		}

		redemptionIDStr := chi.URLParam(r, "id")
		redemptionID, err := strconv.ParseUint(redemptionIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid redemption ID"))
			return
		}

		var req struct {
			Status         string `json:"status" validate:"required,oneof=processing shipped delivered cancelled"`
			TrackingNumber string `json:"tracking_number"`
			Notes          string `json:"notes"`
		}

		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		var redemption store.UserReward
		result := db.Preload("User").First(&redemption, redemptionID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("redemption not found"))
			} else {
				internalServerError(w, r, result.Error)
			}
			return
		}

		oldStatus := redemption.Status
		redemption.Status = req.Status

		if req.TrackingNumber != "" {
			redemption.TrackingNumber = stringPtr(req.TrackingNumber)
		}

		if req.Status == "delivered" {
			now := time.Now()
			redemption.DeliveredAt = &now
		}

		if err := db.Save(&redemption).Error; err != nil {
			internalServerError(w, r, err)
			return
		}

		// Create notification for user
		var message string
		switch req.Status {
		case "processing":
			message = "Your reward is being processed and will be shipped soon."
		case "shipped":
			message = fmt.Sprintf("Your reward has been shipped! Tracking number: %s", req.TrackingNumber)
		case "delivered":
			message = "Your reward has been delivered!"
		case "cancelled":
			message = "Your redemption has been cancelled by admin."
		}

		if req.Notes != "" {
			message += " Note: " + req.Notes
		}

		notificationDataJSON, err := json.Marshal(map[string]interface{}{
			"redemption_id":   redemption.ID,
			"old_status":      oldStatus,
			"new_status":      req.Status,
			"tracking_number": req.TrackingNumber,
			"updated_by":      adminUser.ID,
			"updated_by_name": adminUser.FirstName + " " + adminUser.LastName,
		})
		if err != nil {
			internalServerError(w, r, err)
			return
		}
		notification := &store.Notification{
			UserID:           redemption.UserID,
			NotificationType: "reward_unlocked",
			Title:            "Redemption Status Updated",
			Message:          message,
			Data:             stringPtr(string(notificationDataJSON)),
			IsActionable:     true,
			ActionURL:        stringPtr("/rewards/redemptions/" + fmt.Sprintf("%d", redemption.ID)),
		}
		store.CreateNotification(db, notification)

		// Log admin action
		changesJSON, err := json.Marshal(map[string]interface{}{
			"old_status":      oldStatus,
			"new_status":      req.Status,
			"tracking_number": req.TrackingNumber,
		})
		if err != nil {
			internalServerError(w, r, err)
			return
		}
		adminAction := &store.AdminAction{
			AdminID:      intPtr(int(adminUser.ID)),
			ActionType:   "update_redemption_status",
			ResourceType: "user_reward",
			ResourceID:   intPtr(int(redemption.ID)),
			Changes:      stringPtr(string(changesJSON)),
		}
		store.CreateAdminAction(db, adminAction)

		response := map[string]interface{}{
			"message": "Redemption status updated successfully",
			"redemption": map[string]interface{}{
				"id":              redemption.ID,
				"status":          redemption.Status,
				"tracking_number": redemption.TrackingNumber,
				"created_at":      redemption.CreatedAt,
			},
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Helper function to generate redemption code
func generateRedemptionCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 12)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	// Format as XXXX-XXXX-XXXX
	return fmt.Sprintf("%s-%s-%s", string(b[0:4]), string(b[4:8]), string(b[8:12]))
}
