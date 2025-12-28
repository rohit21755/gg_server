package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// Get XP Transactions
func getXPTransactionsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit

		var transactions []store.XPTransaction
		var totalCount int64

		// Get total count
		db.Model(&store.XPTransaction{}).Where("user_id = ?", user.ID).Count(&totalCount)

		// Get transactions
		result := db.Where("user_id = ?", user.ID).
			Order("created_at DESC").
			Offset(offset).
			Limit(limit).
			Find(&transactions)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		response := map[string]interface{}{
			"transactions": transactions,
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

// Award XP (Admin only)
func awardXPHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Check if user is admin
		if user.Role != "admin" && user.Role != "state_lead" {
			unauthorizedResponse(w, r, errors.New("only admins can award XP"))
			return
		}

		var req struct {
			UserID     uint   `json:"user_id" validate:"required"`
			Amount     int    `json:"amount" validate:"required,min=1"`
			Reason     string `json:"reason" validate:"required"`
			SourceType string `json:"source_type"`
			SourceID   uint   `json:"source_id"`
		}

		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		if err := Validate.Struct(req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Get target user
		targetUser, err := store.GetUserByID(db, req.UserID)
		if err != nil {
			notFoundResponse(w, r, errors.New("user not found"))
			return
		}

		// Create XP transaction
		userIDInt := int(req.UserID)
		sourceIDInt := int(req.SourceID)
		metadataJSON, _ := json.Marshal(map[string]interface{}{
			"awarded_by":      user.ID,
			"awarded_by_name": user.FirstName + " " + user.LastName,
		})
		xpTransaction := &store.XPTransaction{
			UserID:          &userIDInt,
			TransactionType: "bonus",
			Amount:          req.Amount,
			SourceType:      stringPtr(req.SourceType),
			SourceID:        &sourceIDInt,
			Description:     stringPtr(req.Reason),
			Metadata:        stringPtr(string(metadataJSON)),
		}

		if err := store.CreateXPTransaction(db, xpTransaction); err != nil {
			internalServerError(w, r, err)
			return
		}

		// Update user XP
		targetUser.XP += req.Amount
		if err := store.UpdateUser(db, targetUser); err != nil {
			internalServerError(w, r, err)
			return
		}

		// Create notification
		dataJSON, _ := json.Marshal(map[string]interface{}{
			"xp_amount": req.Amount,
			"reason":    req.Reason,
		})
		notification := &store.Notification{
			UserID:           &userIDInt,
			NotificationType: "reward_unlocked",
			Title:            "XP Awarded!",
			Message:          fmt.Sprintf("You received %d XP: %s", req.Amount, req.Reason),
			Data:             stringPtr(string(dataJSON)),
			IsActionable:     true,
		}
		store.CreateNotification(db, notification)

		response := map[string]interface{}{
			"message":     "XP awarded successfully",
			"new_balance": targetUser.XP,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Levels
func getLevelsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var levels []store.Level
		result := db.Order("rank_order ASC").Find(&levels)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		if err := jsonResponse(w, http.StatusOK, levels); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Current Level Info
func getCurrentLevelHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Get current level
		var levelID uint = 1
		if user.LevelID != nil {
			levelID = uint(*user.LevelID)
		}
		currentLevel, err := store.GetLevelByID(db, levelID)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		// Get next level
		var nextLevel store.Level
		db.Where("rank_order > ?", currentLevel.RankOrder).
			Order("rank_order ASC").
			First(&nextLevel)

		// Get progress to next level
		progress := 0
		if nextLevel.ID > 0 {
			xpRange := nextLevel.MinXP - currentLevel.MinXP
			currentProgress := user.XP - currentLevel.MinXP
			if xpRange > 0 {
				progress = (currentProgress * 100) / xpRange
			}
		}

		response := map[string]interface{}{
			"current_level": currentLevel,
			"next_level":    nextLevel,
			"progress": map[string]interface{}{
				"percentage": progress,
				"current_xp": user.XP,
				"xp_to_next": func() int {
					if nextLevel.ID > 0 {
						return nextLevel.MinXP - user.XP
					}
					return 0
				}(),
			},
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Badges
func getBadgesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("category")

		query := db.Model(&store.Badge{})
		if category != "" {
			query = query.Where("category = ?", category)
		}

		var badges []store.Badge
		result := query.Order("created_at DESC").Find(&badges)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		if err := jsonResponse(w, http.StatusOK, badges); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Badge Details
func getBadgeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		badgeIDStr := chi.URLParam(r, "id")
		badgeID, err := strconv.ParseUint(badgeIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid badge ID"))
			return
		}

		badge, err := store.GetBadgeByID(db, uint(badgeID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("badge not found"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		if err := jsonResponse(w, http.StatusOK, badge); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get User Badges
func getUserBadgesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		badges, err := store.GetUserBadges(db, user.ID)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		// Get all badges to show locked ones
		var allBadges []store.Badge
		db.Where("is_secret = ?", false).Order("created_at DESC").Find(&allBadges)

		// Create map of earned badges
		earnedMap := make(map[int]store.UserBadge)
		for _, badge := range badges {
			earnedMap[badge.BadgeID] = badge
		}

		// Build response with earned status
		var response []map[string]interface{}
		for _, badge := range allBadges {
			badgeData := map[string]interface{}{
				"badge":  badge,
				"earned": false,
			}
			badgeIDInt := int(badge.ID)
			if _, earned := earnedMap[badgeIDInt]; earned {
				badgeData["earned"] = true
				badgeData["earned_at"] = earnedMap[badgeIDInt].EarnedAt
			}
			response = append(response, badgeData)
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Streak Info
func getStreakHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		streakType := r.URL.Query().Get("type")
		if streakType == "" {
			streakType = "daily_login"
		}

		streak, err := store.GetUserStreak(db, user.ID, streakType)
		if err != nil {
			// Return empty streak if not found
			response := map[string]interface{}{
				"streak_type":    streakType,
				"current_streak": 0,
				"longest_streak": 0,
				"total_days":     0,
				"last_activity":  nil,
			}
			jsonResponse(w, http.StatusOK, response)
			return
		}

		// Get recent streak logs
		var logs []store.StreakLog
		db.Where("user_id = ? AND streak_type = ?", user.ID, streakType).
			Order("activity_date DESC").
			Limit(30).
			Find(&logs)

		response := map[string]interface{}{
			"streak_type":    streak.StreakType,
			"current_streak": streak.CurrentStreak,
			"longest_streak": streak.LongestStreak,
			"total_days":     streak.TotalDays,
			"last_activity":  streak.LastActivityDate,
			"recent_logs":    logs,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Log Daily Activity
func logStreakHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		var req struct {
			ActivityType string `json:"activity_type" validate:"required"`
		}

		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		if err := Validate.Struct(req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Update streak
		if err := updateUserStreak(db, user.ID, req.ActivityType); err != nil {
			internalServerError(w, r, err)
			return
		}

		// Create streak log
		userIDInt := int(user.ID)
		log := &store.StreakLog{
			UserID:       &userIDInt,
			StreakType:   req.ActivityType,
			ActivityDate: time.Now().Truncate(24 * time.Hour),
			EarnedXP:     10, // XP for maintaining streak
		}
		store.CreateStreakLog(db, log)

		response := map[string]string{
			"message": "Activity logged successfully",
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Spin Wheel Config
func getSpinWheelHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Get active spin wheel
		var spinWheel store.SpinWheel
		result := db.Where("is_active = ? AND start_date <= ? AND (end_date IS NULL OR end_date >= ?)",
			true, time.Now(), time.Now()).
			First(&spinWheel)

		if result.Error != nil {
			notFoundResponse(w, r, errors.New("no active spin wheel found"))
			return
		}

		// Get wheel items
		var items []store.SpinWheelItem
		db.Where("spin_wheel_id = ? AND is_active = ?", spinWheel.ID, true).
			Order("sort_order ASC").
			Find(&items)

		// Check if user has spins remaining
		spinsToday, _ := store.GetUserSpinsToday(db, user.ID, spinWheel.ID)
		spinsRemaining := spinWheel.SpinsPerUser - spinsToday

		response := map[string]interface{}{
			"spin_wheel": spinWheel,
			"items":      items,
			"user_stats": map[string]interface{}{
				"spins_remaining": spinsRemaining,
				"spins_today":     spinsToday,
				"can_spin":        spinsRemaining > 0,
			},
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Spin the Wheel
func spinWheelHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Get active spin wheel
		var spinWheel store.SpinWheel
		result := db.Where("is_active = ? AND start_date <= ? AND (end_date IS NULL OR end_date >= ?)",
			true, time.Now(), time.Now()).
			First(&spinWheel)

		if result.Error != nil {
			notFoundResponse(w, r, errors.New("no active spin wheel found"))
			return
		}

		// Check spins remaining
		spinsToday, _ := store.GetUserSpinsToday(db, user.ID, spinWheel.ID)
		if spinsToday >= spinWheel.SpinsPerUser {
			badRequestResponse(w, r, errors.New("no spins remaining for today"))
			return
		}

		// Check min activity level
		var userLevelID int = 1
		if user.LevelID != nil {
			userLevelID = *user.LevelID
		}
		if userLevelID < spinWheel.MinActivityLevel {
			badRequestResponse(w, r, errors.New("minimum level not reached"))
			return
		}

		// Get wheel items with probabilities
		var items []store.SpinWheelItem
		db.Where("spin_wheel_id = ? AND is_active = ?", spinWheel.ID, true).
			Order("sort_order ASC").
			Find(&items)

		if len(items) == 0 {
			internalServerError(w, r, errors.New("no items available on spin wheel"))
			return
		}

		// Select random item based on probability
		selectedItem := selectRandomItem(items)

		// Check quantity
		if selectedItem.MaxQuantity != nil && *selectedItem.MaxQuantity > 0 {
			if selectedItem.CurrentQuantity != nil && *selectedItem.CurrentQuantity >= *selectedItem.MaxQuantity {
				badRequestResponse(w, r, errors.New("item out of stock"))
				return
			}
		}

		// Record spin
		userIDInt := int(user.ID)
		wheelIDInt := int(spinWheel.ID)
		itemIDInt := int(selectedItem.ID)
		userSpin := &store.UserSpin{
			UserID:          &userIDInt,
			SpinWheelID:     &wheelIDInt,
			SpinWheelItemID: &itemIDInt,
			EarnedValue:     selectedItem.ItemValue,
			SpunAt:          time.Now(),
		}

		if err := store.CreateUserSpin(db, userSpin); err != nil {
			internalServerError(w, r, err)
			return
		}

		// Update item quantity
		if selectedItem.MaxQuantity != nil && *selectedItem.MaxQuantity > 0 {
			if selectedItem.CurrentQuantity == nil {
				currentQty := 1
				selectedItem.CurrentQuantity = &currentQty
			} else {
				newQty := *selectedItem.CurrentQuantity + 1
				selectedItem.CurrentQuantity = &newQty
			}
			db.Save(&selectedItem)
		}

		// Award prize
		var rewardDetails map[string]interface{}
		switch selectedItem.ItemType {
		case "xp":
			// Award XP
			spinIDInt := int(userSpin.ID)
			xpTransaction := &store.XPTransaction{
				UserID:          &userIDInt,
				TransactionType: "spin_wheel",
				Amount:          selectedItem.ItemValue,
				SourceType:      stringPtr("spin_wheel"),
				SourceID:        &spinIDInt,
				Description:     stringPtr("Spin wheel reward: " + selectedItem.ItemLabel),
			}
			store.CreateXPTransaction(db, xpTransaction)

			user.XP += selectedItem.ItemValue
			store.UpdateUser(db, user)

			rewardDetails = map[string]interface{}{
				"type":  "xp",
				"value": selectedItem.ItemValue,
			}

		case "badge":
			// Award badge
			userBadge := &store.UserBadge{
				UserID:  userIDInt,
				BadgeID: selectedItem.ItemValue,
			}
			store.CreateUserBadge(db, userBadge)

			rewardDetails = map[string]interface{}{
				"type":  "badge",
				"value": selectedItem.ItemValue,
			}

		case "physical":
			// Create reward redemption
			userReward := &store.UserReward{
				UserID:    &userIDInt,
				Status:    "pending",
				ClaimedAt: time.Now(),
			}
			store.CreateUserReward(db, userReward)

			rewardDetails = map[string]interface{}{
				"type":  "physical",
				"value": selectedItem.ItemLabel,
			}

		default:
			rewardDetails = map[string]interface{}{
				"type":  selectedItem.ItemType,
				"value": selectedItem.ItemValue,
			}
		}

		response := map[string]interface{}{
			"spin_result": map[string]interface{}{
				"item":      selectedItem,
				"reward":    rewardDetails,
				"spin_id":   userSpin.ID,
				"timestamp": userSpin.SpunAt,
			},
			"remaining_spins": spinWheel.SpinsPerUser - spinsToday - 1,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Spin History
func getSpinHistoryHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit

		var spins []store.UserSpin
		var totalCount int64

		// Get total count
		db.Model(&store.UserSpin{}).Where("user_id = ?", user.ID).Count(&totalCount)

		// Get spins with item details
		result := db.Where("user_id = ?", user.ID).
			Preload("SpinWheelItem").
			Order("spun_at DESC").
			Offset(offset).
			Limit(limit).
			Find(&spins)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		response := map[string]interface{}{
			"spins": spins,
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

// Helper function to select random item based on probability
func selectRandomItem(items []store.SpinWheelItem) store.SpinWheelItem {
	// Calculate cumulative probabilities
	var cumulative float64
	var ranges []struct {
		item  store.SpinWheelItem
		start float64
		end   float64
	}

	for _, item := range items {
		ranges = append(ranges, struct {
			item  store.SpinWheelItem
			start float64
			end   float64
		}{
			item:  item,
			start: cumulative,
			end:   cumulative + item.Probability,
		})
		cumulative += item.Probability
	}

	// Generate random number
	rand.Seed(time.Now().UnixNano())
	r := rand.Float64()

	// Find selected item
	for _, rng := range ranges {
		if r >= rng.start && r < rng.end {
			return rng.item
		}
	}

	// Fallback to first item
	return items[0]
}
