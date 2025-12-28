package main

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// Get Active Flash Challenges
func getActiveFlashChallengesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var challenges []store.FlashChallenge
		result := db.Where("status = 'active' AND start_time <= ? AND end_time >= ?",
			time.Now(), time.Now()).
			Order("end_time ASC").
			Find(&challenges)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Calculate time remaining for each challenge
		var response []map[string]interface{}
		for _, challenge := range challenges {
			timeRemaining := challenge.EndTime.Sub(time.Now())
			response = append(response, map[string]interface{}{
				"challenge":       challenge,
				"time_remaining":  int(timeRemaining.Seconds()),
				"hours_remaining": int(timeRemaining.Hours()),
			})
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Participate in Flash Challenge
func participateFlashChallengeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		challengeIDStr := chi.URLParam(r, "id")
		challengeID, err := strconv.ParseUint(challengeIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid challenge ID"))
			return
		}

		var challenge store.FlashChallenge
		result := db.First(&challenge, challengeID)
		if result.Error != nil {
			notFoundResponse(w, r, errors.New("challenge not found"))
			return
		}

		// Check if challenge is active
		if challenge.Status != "active" || time.Now().Before(challenge.StartTime) || time.Now().After(challenge.EndTime) {
			badRequestResponse(w, r, errors.New("challenge is not active"))
			return
		}

		// Check max participants
		var participantCount int64
		db.Model(&store.Submission{}).
			Where("task_id IN (SELECT id FROM tasks WHERE campaign_id = ?)", challenge.ID).
			Group("user_id").
			Count(&participantCount)

		if challenge.MaxParticipants != nil && *challenge.MaxParticipants > 0 && participantCount >= int64(*challenge.MaxParticipants) {
			conflictResponse(w, r, errors.New("challenge has reached maximum participants"))
			return
		}

		// Create a task for this challenge if it doesn't exist
		var task store.Task
		db.Where("campaign_id = ?", challenge.ID).First(&task)
		if task.ID == 0 {
			// Create task for flash challenge
			challengeIDInt := int(challenge.ID)
			var description string
			if challenge.Description != nil {
				description = *challenge.Description
			}
			task = store.Task{
				CampaignID:  &challengeIDInt,
				Title:       challenge.Title,
				Description: description,
				TaskType:    "solo",
				ProofType:   "url",
				XPReward:    challenge.XPReward,
				Priority:    "flash",
				IsActive:    true,
			}
			db.Create(&task)
		}

		response := map[string]interface{}{
			"message":   "Successfully joined flash challenge",
			"challenge": challenge,
			"task_id":   task.ID,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Active Trivia
func getActiveTriviaHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var trivia []store.TriviaTournament
		result := db.Where("status = 'active' AND start_date <= ? AND end_date >= ?",
			time.Now(), time.Now()).
			Order("end_date ASC").
			Find(&trivia)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Don't send questions, just metadata
		var response []map[string]interface{}
		for _, t := range trivia {
			// Parse questions to get count
			var questions []interface{}
			json.Unmarshal([]byte(t.Questions), &questions)

			response = append(response, map[string]interface{}{
				"id":             t.ID,
				"title":          t.Title,
				"description":    t.Description,
				"start_date":     t.StartDate,
				"end_date":       t.EndDate,
				"duration":       t.DurationMinutes,
				"question_count": len(questions),
				"entry_fee":      t.EntryFeeXP,
				"time_remaining": int(t.EndDate.Sub(time.Now()).Seconds()),
			})
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Start Trivia
func startTriviaHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		triviaIDStr := chi.URLParam(r, "id")
		triviaID, err := strconv.ParseUint(triviaIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid trivia ID"))
			return
		}

		var trivia store.TriviaTournament
		result := db.First(&trivia, triviaID)
		if result.Error != nil {
			notFoundResponse(w, r, errors.New("trivia not found"))
			return
		}

		// Check if trivia is active
		if trivia.Status != "active" || time.Now().Before(trivia.StartDate) || time.Now().After(trivia.EndDate) {
			badRequestResponse(w, r, errors.New("trivia is not active"))
			return
		}

		// Check if user has already participated
		var participant store.TriviaParticipant
		db.Where("trivia_id = ? AND user_id = ?", trivia.ID, user.ID).First(&participant)
		if participant.ID > 0 {
			conflictResponse(w, r, errors.New("you have already participated in this trivia"))
			return
		}

		// Check entry fee
		if trivia.EntryFeeXP > 0 && user.XP < trivia.EntryFeeXP {
			badRequestResponse(w, r, errors.New("insufficient XP for entry fee"))
			return
		}

		// Deduct entry fee
		if trivia.EntryFeeXP > 0 {
			userIDInt := int(user.ID)
			triviaIDInt := int(trivia.ID)
			xpTransaction := &store.XPTransaction{
				UserID:          &userIDInt,
				TransactionType: "quiz",
				Amount:          -trivia.EntryFeeXP,
				SourceType:      stringPtr("trivia"),
				SourceID:        &triviaIDInt,
				Description:     stringPtr("Trivia entry fee"),
			}
			store.CreateXPTransaction(db, xpTransaction)

			user.XP -= trivia.EntryFeeXP
			store.UpdateUser(db, user)
		}

		// Create participant record
		triviaIDInt := int(trivia.ID)
		userIDInt := int(user.ID)
		participant = store.TriviaParticipant{
			TriviaID:       &triviaIDInt,
			UserID:         &userIDInt,
			ParticipatedAt: time.Now(),
		}
		store.CreateTriviaParticipant(db, &participant)

		// Return questions (without answers)
		var questions []map[string]interface{}
		json.Unmarshal([]byte(trivia.Questions), &questions)

		// Remove correct answers
		for i := range questions {
			delete(questions[i], "correct_answer")
		}

		response := map[string]interface{}{
			"trivia_id":      trivia.ID,
			"title":          trivia.Title,
			"duration":       trivia.DurationMinutes,
			"questions":      questions,
			"start_time":     time.Now(),
			"end_time":       time.Now().Add(time.Duration(trivia.DurationMinutes) * time.Minute),
			"participant_id": participant.ID,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Submit Trivia Answers
func submitTriviaAnswersHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		triviaIDStr := chi.URLParam(r, "id")
		triviaID, err := strconv.ParseUint(triviaIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid trivia ID"))
			return
		}

		var req struct {
			Answers []struct {
				QuestionID int    `json:"question_id"`
				Answer     string `json:"answer"`
			} `json:"answers" validate:"required"`
			TimeTaken int `json:"time_taken" validate:"required"`
		}

		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Get trivia
		var trivia store.TriviaTournament
		result := db.First(&trivia, triviaID)
		if result.Error != nil {
			notFoundResponse(w, r, errors.New("trivia not found"))
			return
		}

		// Check if trivia is still active
		if time.Now().After(trivia.EndDate) {
			badRequestResponse(w, r, errors.New("trivia has ended"))
			return
		}

		// Get participant
		var participant store.TriviaParticipant
		db.Where("trivia_id = ? AND user_id = ?", trivia.ID, user.ID).First(&participant)
		if participant.ID == 0 {
			notFoundResponse(w, r, errors.New("participant not found"))
			return
		}

		// Check if already submitted
		if participant.Score > 0 {
			conflictResponse(w, r, errors.New("answers already submitted"))
			return
		}

		// Calculate score
		var questions []map[string]interface{}
		json.Unmarshal([]byte(trivia.Questions), &questions)

		score := 0
		correctAnswers := 0

		for _, answer := range req.Answers {
			if answer.QuestionID < len(questions) {
				question := questions[answer.QuestionID]
				if correctAnswer, ok := question["correct_answer"].(string); ok && correctAnswer == answer.Answer {
					score += 10 // 10 points per correct answer
					correctAnswers++
				}
			}
		}

		// Update participant
		participant.Score = score
		participant.CorrectAnswers = correctAnswers
		participant.TimeTakenSeconds = intPtr(req.TimeTaken)
		db.Save(&participant)

		// Award XP based on score
		if score > 0 {
			xpEarned := score * 5 // 5 XP per point
			userIDInt := int(user.ID)
			triviaIDInt := int(trivia.ID)
			xpTransaction := &store.XPTransaction{
				UserID:          &userIDInt,
				TransactionType: "quiz",
				Amount:          xpEarned,
				SourceType:      stringPtr("trivia"),
				SourceID:        &triviaIDInt,
				Description:     stringPtr("Trivia competition reward"),
			}
			store.CreateXPTransaction(db, xpTransaction)

			user.XP += xpEarned
			store.UpdateUser(db, user)
		}

		response := map[string]interface{}{
			"score":           score,
			"correct_answers": correctAnswers,
			"total_questions": len(questions),
			"xp_earned":       score * 5,
			"time_taken":      req.TimeTaken,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Mystery Boxes
func getMysteryBoxesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var boxes []store.MysteryBox
		result := db.Where("is_active = ?", true).
			Order("cost_xp ASC").
			Find(&boxes)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		if err := jsonResponse(w, http.StatusOK, boxes); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Open Mystery Box
func openMysteryBoxHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		boxIDStr := chi.URLParam(r, "id")
		boxID, err := strconv.ParseUint(boxIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid mystery box ID"))
			return
		}

		var box store.MysteryBox
		result := db.First(&box, boxID)
		if result.Error != nil {
			notFoundResponse(w, r, errors.New("mystery box not found"))
			return
		}

		// Check if box is active
		if !box.IsActive {
			badRequestResponse(w, r, errors.New("mystery box is not available"))
			return
		}

		// Check if user has enough XP
		if user.XP < box.CostXP {
			badRequestResponse(w, r, errors.New("insufficient XP"))
			return
		}

		// Parse contents
		var contents []map[string]interface{}
		if err := json.Unmarshal([]byte(box.Contents), &contents); err != nil {
			internalServerError(w, r, err)
			return
		}

		// Select random reward
		if len(contents) == 0 {
			internalServerError(w, r, errors.New("no rewards available"))
			return
		}

		rand.Seed(time.Now().UnixNano())
		selected := contents[rand.Intn(len(contents))]

		// Deduct XP
		userIDInt := int(user.ID)
		boxIDInt := int(box.ID)
		xpTransaction := &store.XPTransaction{
			UserID:          &userIDInt,
			TransactionType: "mystery_box",
			Amount:          -box.CostXP,
			SourceType:      stringPtr("mystery_box"),
			SourceID:        &boxIDInt,
			Description:     stringPtr("Mystery box purchase"),
		}
		store.CreateXPTransaction(db, xpTransaction)

		user.XP -= box.CostXP
		store.UpdateUser(db, user)

		// Record redemption
		redemption := &store.MysteryBoxRedemption{
			UserID:       &userIDInt,
			MysteryBoxID: &boxIDInt,
			RewardType:   selected["type"].(string),
			RewardValue:  int(selected["value"].(float64)),
			RedeemedAt:   time.Now(),
		}
		store.CreateMysteryBoxRedemption(db, redemption)

		// Award reward
		rewardType := selected["type"].(string)
		rewardValue := int(selected["value"].(float64))

		switch rewardType {
		case "xp":
			userIDInt := int(user.ID)
			redemptionIDInt := int(redemption.ID)
			xpTransaction := &store.XPTransaction{
				UserID:          &userIDInt,
				TransactionType: "mystery_box",
				Amount:          rewardValue,
				SourceType:      stringPtr("mystery_box"),
				SourceID:        &redemptionIDInt,
				Description:     stringPtr("Mystery box reward"),
			}
			store.CreateXPTransaction(db, xpTransaction)

			user.XP += rewardValue
			store.UpdateUser(db, user)

		case "badge":
			userIDInt := int(user.ID)
			userBadge := &store.UserBadge{
				UserID:  userIDInt,
				BadgeID: rewardValue,
			}
			store.CreateUserBadge(db, userBadge)

			// Add other reward types as needed
		}

		response := map[string]interface{}{
			"reward": map[string]interface{}{
				"type":  rewardType,
				"value": rewardValue,
				"label": selected["label"],
			},
			"remaining_xp": user.XP,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Redeem Secret Code
func redeemSecretCodeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		code := chi.URLParam(r, "code")
		if code == "" {
			badRequestResponse(w, r, errors.New("code is required"))
			return
		}

		var secretCode store.SecretCode
		result := db.Where("code = ? AND is_active = ? AND valid_from <= ? AND valid_until >= ?",
			code, true, time.Now(), time.Now()).
			First(&secretCode)

		if result.Error != nil {
			notFoundResponse(w, r, errors.New("invalid or expired code"))
			return
		}

		// Check max redemptions
		if secretCode.MaxRedemptions > 0 && secretCode.CurrentRedemptions >= secretCode.MaxRedemptions {
			conflictResponse(w, r, errors.New("code has reached maximum redemptions"))
			return
		}

		// Check if user has already redeemed
		var existingRedemption store.SecretCodeRedemption
		db.Where("secret_code_id = ? AND user_id = ?", secretCode.ID, user.ID).First(&existingRedemption)
		if existingRedemption.ID > 0 {
			conflictResponse(w, r, errors.New("code already redeemed"))
			return
		}

		// Record redemption
		secretCodeIDInt := int(secretCode.ID)
		userIDInt := int(user.ID)
		redemption := &store.SecretCodeRedemption{
			SecretCodeID: &secretCodeIDInt,
			UserID:       &userIDInt,
			RedeemedAt:   time.Now(),
		}
		store.CreateSecretCodeRedemption(db, redemption)

		// Update redemption count
		secretCode.CurrentRedemptions++
		db.Save(&secretCode)

		// Award XP
		if secretCode.XPReward > 0 {
			description := "Secret code redemption"
			if secretCode.Description != nil {
				description = description + ": " + *secretCode.Description
			}
			xpTransaction := &store.XPTransaction{
				UserID:          &userIDInt,
				TransactionType: "bonus",
				Amount:          secretCode.XPReward,
				SourceType:      stringPtr("secret_code"),
				SourceID:        &secretCodeIDInt,
				Description:     stringPtr(description),
			}
			store.CreateXPTransaction(db, xpTransaction)

			user.XP += secretCode.XPReward
			store.UpdateUser(db, user)
		}

		// Award badge if specified
		if secretCode.BadgeID != nil {
			userBadge := &store.UserBadge{
				UserID:  userIDInt,
				BadgeID: int(*secretCode.BadgeID),
			}
			store.CreateUserBadge(db, userBadge)
		}

		response := map[string]interface{}{
			"message": "Code redeemed successfully",
			"rewards": map[string]interface{}{
				"xp":       secretCode.XPReward,
				"coins":    secretCode.CoinReward,
				"badge_id": secretCode.BadgeID,
			},
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Weekly Vibe Challenge
func getWeeklyVibeChallengeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get current week's challenge
		startOfWeek := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -int(time.Now().Weekday()))
		endOfWeek := startOfWeek.AddDate(0, 0, 7)

		var challenge store.Campaign
		result := db.Where("campaign_type = 'weekly_vibe' AND start_date <= ? AND end_date >= ?",
			endOfWeek, startOfWeek).
			First(&challenge)

		if result.Error != nil {
			// Create a default response if no challenge found
			response := map[string]interface{}{
				"message": "No weekly challenge this week",
				"theme":   "Check back soon for next week's challenge!",
			}
			jsonResponse(w, http.StatusOK, response)
			return
		}

		// Get submission count
		var submissionCount int64
		db.Model(&store.Submission{}).Where("campaign_id = ?", challenge.ID).Count(&submissionCount)

		response := map[string]interface{}{
			"challenge": challenge,
			"stats": map[string]interface{}{
				"submission_count": submissionCount,
				"days_remaining":   int(endOfWeek.Sub(time.Now()).Hours() / 24),
			},
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Submit to Weekly Vibe Challenge
func submitWeeklyVibeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Find current weekly challenge
		startOfWeek := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -int(time.Now().Weekday()))
		endOfWeek := startOfWeek.AddDate(0, 0, 7)

		var challenge store.Campaign
		result := db.Where("campaign_type = 'weekly_vibe' AND start_date <= ? AND end_date >= ?",
			endOfWeek, startOfWeek).
			First(&challenge)

		if result.Error != nil {
			notFoundResponse(w, r, errors.New("no active weekly challenge"))
			return
		}

		// Check if user has already submitted
		var existingSubmission store.Submission
		db.Where("campaign_id = ? AND user_id = ?", challenge.ID, user.ID).First(&existingSubmission)
		if existingSubmission.ID > 0 {
			conflictResponse(w, r, errors.New("you have already submitted to this week's challenge"))
			return
		}

		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			ProofURL    string `json:"proof_url" validate:"required,url"`
			ProofType   string `json:"proof_type" validate:"required"`
		}

		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Create task for weekly challenge if it doesn't exist
		var task store.Task
		db.Where("campaign_id = ?", challenge.ID).First(&task)
		if task.ID == 0 {
			challengeIDInt := int(challenge.ID)
			var description string
			if challenge.Description != nil {
				description = *challenge.Description
			}
			task = store.Task{
				CampaignID:  &challengeIDInt,
				Title:       "Weekly Vibe Challenge",
				Description: description,
				TaskType:    "solo",
				ProofType:   "url",
				XPReward:    500,
				Priority:    "high",
				IsActive:    true,
			}
			store.CreateTask(db, &task)
		}

		// Create submission
		taskIDInt := int(task.ID)
		userIDInt := int(user.ID)
		challengeIDInt := int(challenge.ID)
		submission := &store.Submission{
			TaskID:      &taskIDInt,
			UserID:      &userIDInt,
			CampaignID:  &challengeIDInt,
			ProofType:   req.ProofType,
			ProofURL:    req.ProofURL,
			ProofText:   stringPtr(req.Description),
			Status:      "pending",
			SubmittedAt: time.Now(),
		}
		store.CreateSubmission(db, submission)

		response := map[string]interface{}{
			"message":       "Submission received! Voting starts after the submission deadline.",
			"submission_id": submission.ID,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Weekly Challenge Submissions
func getWeeklyChallengeSubmissionsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Find current weekly challenge
		startOfWeek := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -int(time.Now().Weekday()))
		endOfWeek := startOfWeek.AddDate(0, 0, 7)

		var challenge store.Campaign
		result := db.Where("campaign_type = 'weekly_vibe' AND start_date <= ? AND end_date >= ?",
			endOfWeek, startOfWeek).
			First(&challenge)

		if result.Error != nil {
			notFoundResponse(w, r, errors.New("no active weekly challenge"))
			return
		}

		// Get submissions for this challenge
		var submissions []store.Submission
		db.Where("campaign_id = ?", challenge.ID).
			Preload("User").
			Order("submitted_at DESC").
			Find(&submissions)

		// Format response
		var response []map[string]interface{}
		for _, sub := range submissions {
			response = append(response, map[string]interface{}{
				"id":           sub.ID,
				"proof_url":    sub.ProofURL,
				"description":  sub.ProofText,
				"submitted_at": sub.SubmittedAt,
				"user": map[string]interface{}{
					"id":         sub.User.ID,
					"first_name": sub.User.FirstName,
					"last_name":  sub.User.LastName,
					"college_id": sub.User.CollegeID,
				},
			})
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Vote for Weekly Challenge Submission
func voteWeeklyChallengeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		submissionIDStr := chi.URLParam(r, "submissionId")
		submissionID, err := strconv.ParseUint(submissionIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid submission ID"))
			return
		}

		// Get submission
		var submission store.Submission
		result := db.First(&submission, submissionID)
		if result.Error != nil {
			notFoundResponse(w, r, errors.New("submission not found"))
			return
		}

		// Find the campaign
		var campaign store.Campaign
		db.First(&campaign, submission.CampaignID)
		if campaign.CampaignType != "weekly_vibe" {
			badRequestResponse(w, r, errors.New("not a weekly challenge submission"))
			return
		}

		// Check if voting period
		now := time.Now()
		if now.Before(campaign.StartDate) || now.After(campaign.EndDate) {
			badRequestResponse(w, r, errors.New("voting is not currently active"))
			return
		}

		// Check if user has already voted in this campaign
		var existingVote store.BattleVote
		db.Where("battle_id = ? AND voter_id = ?", campaign.ID, user.ID).First(&existingVote)
		if existingVote.ID > 0 {
			conflictResponse(w, r, errors.New("you have already voted in this challenge"))
			return
		}

		// Create vote
		campaignIDInt := int(campaign.ID)
		submissionIDInt := int(submission.ID)
		userIDInt := int(user.ID)
		vote := store.BattleVote{
			BattleID:     &campaignIDInt,
			SubmissionID: &submissionIDInt,
			VoterID:      &userIDInt,
			VotedAt:      now,
		}
		store.CreateBattleVote(db, &vote)

		// Update submission vote count
		submission.Score = func(score *float64) *float64 {
			if score == nil {
				newScore := 1.0
				return &newScore
			}
			*score += 1
			return score
		}(submission.Score)
		db.Save(&submission)

		response := map[string]string{
			"message": "Vote recorded successfully",
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Get Active Content Battles
func getActiveBattlesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		var battles []store.ContentBattle
		result := db.Where("status IN ('submissions', 'voting') OR (status = 'upcoming' AND submission_deadline > ?)", now).
			Order("submission_deadline ASC").
			Find(&battles)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Add status info
		var response []map[string]interface{}
		for _, battle := range battles {
			status := "upcoming"
			if now.After(battle.VotingEnd) {
				status = "completed"
			} else if now.After(battle.VotingStart) {
				status = "voting"
			} else if now.After(battle.SubmissionDeadline) {
				status = "submissions_closed"
			} else {
				status = "submissions_open"
			}

			// Get submission count
			var submissionCount int64
			db.Model(&store.BattleSubmission{}).Where("battle_id = ?", battle.ID).Count(&submissionCount)

			response = append(response, map[string]interface{}{
				"battle":           battle,
				"status":           status,
				"submission_count": submissionCount,
				"time_remaining": func() int {
					switch status {
					case "submissions_open":
						return int(battle.SubmissionDeadline.Sub(now).Seconds())
					case "voting":
						return int(battle.VotingEnd.Sub(now).Seconds())
					default:
						return 0
					}
				}(),
			})
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Submit to Content Battle
func submitBattleHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		battleIDStr := chi.URLParam(r, "id")
		battleID, err := strconv.ParseUint(battleIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid battle ID"))
			return
		}

		var battle store.ContentBattle
		result := db.First(&battle, battleID)
		if result.Error != nil {
			notFoundResponse(w, r, errors.New("battle not found"))
			return
		}

		// Check if submissions are open
		now := time.Now()
		if now.After(battle.SubmissionDeadline) {
			badRequestResponse(w, r, errors.New("submission deadline has passed"))
			return
		}

		// Check if user has already submitted
		var existingSubmission store.BattleSubmission
		db.Where("battle_id = ? AND user_id = ?", battle.ID, user.ID).First(&existingSubmission)
		if existingSubmission.ID > 0 {
			conflictResponse(w, r, errors.New("you have already submitted to this battle"))
			return
		}

		// Check max participants
		var submissionCount int64
		db.Model(&store.BattleSubmission{}).Where("battle_id = ?", battle.ID).Count(&submissionCount)
		if battle.MaxParticipants != nil && *battle.MaxParticipants > 0 && submissionCount >= int64(*battle.MaxParticipants) {
			conflictResponse(w, r, errors.New("battle has reached maximum participants"))
			return
		}

		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			MediaURL    string `json:"media_url" validate:"required,url"`
		}

		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Create submission
		battleIDInt := int(battle.ID)
		userIDInt := int(user.ID)
		submission := &store.BattleSubmission{
			BattleID:    &battleIDInt,
			UserID:      &userIDInt,
			Title:       stringPtr(req.Title),
			Description: stringPtr(req.Description),
			MediaURL:    req.MediaURL,
			SubmittedAt: now,
		}
		store.CreateBattleSubmission(db, submission)

		// Award participation XP
		xpTransaction := &store.XPTransaction{
			UserID:          &userIDInt,
			TransactionType: "battle_participation",
			Amount:          100,
			SourceType:      stringPtr("content_battle"),
			SourceID:        &battleIDInt,
			Description:     stringPtr("Content battle participation"),
		}
		store.CreateXPTransaction(db, xpTransaction)

		user.XP += 100
		store.UpdateUser(db, user)

		response := map[string]interface{}{
			"message":       "Submission received! Good luck!",
			"submission_id": submission.ID,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Vote in Content Battle
func voteBattleHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		battleIDStr := chi.URLParam(r, "id")
		battleID, err := strconv.ParseUint(battleIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid battle ID"))
			return
		}

		submissionIDStr := chi.URLParam(r, "submissionId")
		submissionID, err := strconv.ParseUint(submissionIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid submission ID"))
			return
		}

		var battle store.ContentBattle
		result := db.First(&battle, battleID)
		if result.Error != nil {
			notFoundResponse(w, r, errors.New("battle not found"))
			return
		}

		// Check if voting is open
		now := time.Now()
		if now.Before(battle.VotingStart) || now.After(battle.VotingEnd) {
			badRequestResponse(w, r, errors.New("voting is not currently active"))
			return
		}

		// Check if submission exists
		var submission store.BattleSubmission
		db.First(&submission, submissionID)
		battleIDInt := int(battle.ID)
		if submission.ID == 0 || submission.BattleID == nil || *submission.BattleID != battleIDInt {
			notFoundResponse(w, r, errors.New("submission not found"))
			return
		}

		// Check if user has already voted in this battle
		var existingVote store.BattleVote
		userIDInt := int(user.ID)
		db.Where("battle_id = ? AND voter_id = ?", battleIDInt, userIDInt).First(&existingVote)
		if existingVote.ID > 0 {
			conflictResponse(w, r, errors.New("you have already voted in this battle"))
			return
		}

		// Create vote
		submissionIDInt := int(submission.ID)
		vote := store.BattleVote{
			BattleID:     &battleIDInt,
			SubmissionID: &submissionIDInt,
			VoterID:      &userIDInt,
			VotedAt:      now,
		}
		store.CreateBattleVote(db, &vote)

		// Update submission vote count
		submission.VoteCount++
		db.Save(&submission)

		// Award XP for voting
		xpTransaction := &store.XPTransaction{
			UserID:          &userIDInt,
			TransactionType: "battle_vote",
			Amount:          10,
			SourceType:      stringPtr("content_battle"),
			SourceID:        &battleIDInt,
			Description:     stringPtr("Content battle voting"),
		}
		store.CreateXPTransaction(db, xpTransaction)

		user.XP += 10
		store.UpdateUser(db, user)

		response := map[string]interface{}{
			"message":     "Vote recorded successfully",
			"total_votes": submission.VoteCount,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}
