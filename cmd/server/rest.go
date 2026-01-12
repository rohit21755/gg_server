package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func setupREST(r chi.Router, db *gorm.DB) {
	log.Println("Setting up REST API")

	if db == nil {
		log.Fatal("Database connection is nil, cannot setup REST API")
	}

	// CORS middleware

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Get("/health", healthHandler)
		log.Println("Health endpoint registered at /api/v1/health")

		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", loginHandler(db))
			r.Post("/register", registerHandler(db))
			r.Post("/refresh", refreshTokenHandler(db))
			r.Post("/logout", logoutHandler(db))
			r.Post("/forgot-password", forgotPasswordHandler(db))
			r.Post("/reset-password", resetPasswordHandler(db))
			r.Get("/verify-email/{token}", verifyEmailHandler(db))
		})

		// Public routes
		r.Get("/colleges", getCollegesHandler(db))
		r.Get("/states", getStatesHandler(db))
		r.Get("/leaderboards/global", getGlobalLeaderboardHandler(db))

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(RequireAuth(db))

			// User routes
			r.Route("/users", func(r chi.Router) {
				r.Get("/me", getCurrentUserProfileHandler(db))
				r.Put("/me", updateUserProfileHandler(db))
				r.Patch("/me/avatar", updateAvatarHandler(db))
				r.Post("/me/resume", uploadResumeHandler(db))
				r.Get("/me/certificates", getUserCertificatesHandler(db))
				r.Get("/me/certificates/{id}/download", downloadCertificateHandler(db))
				r.Get("/me/stats", getUserDashboardStatsHandler(db))
				r.Get("/me/activity", getUserActivityHandler(db))
				r.Get("/search", searchUsersHandler(db))
				r.Get("/{id}/stats", getUserStatsHandler(db))
			})

			// Task routes
			r.Route("/tasks", func(r chi.Router) {
				r.Get("/", getTasksHandler(db))
				r.Get("/{id}", getTaskHandler(db))
				r.Get("/assigned", getAssignedTasksHandler(db))
				r.Get("/available", getAvailableTasksHandler(db))
			})

			// Submission routes
			r.Route("/submissions", func(r chi.Router) {
				r.Get("/", getSubmissionsHandler(db))
				r.Post("/", createSubmissionHandler(db))
				r.Get("/{id}", getSubmissionHandler(db))
				r.Put("/{id}", updateSubmissionHandler(db))
				r.Delete("/{id}", deleteSubmissionHandler(db))
				r.Get("/{id}/proof", getSubmissionProofHandler(db))
				r.Post("/{id}/appeal", func(w http.ResponseWriter, r *http.Request) {
					// TODO: Implement appeal handler
					writeJSONError(w, http.StatusNotImplemented, "appeal not yet implemented")
				})
			})

			// Campaign routes
			r.Route("/campaigns", func(r chi.Router) {
				r.Get("/", getCampaignsHandler(db))
				r.Get("/{id}", getCampaignHandler(db))
				r.Post("/{id}/join", joinCampaignHandler(db))
				r.Get("/{id}/tasks", getCampaignTasksHandler(db))
				r.Get("/{id}/leaderboard", getCampaignLeaderboardHandler(db))
			})

			// Gamification routes
			r.Route("/xp", func(r chi.Router) {
				r.Get("/transactions", getXPTransactionsHandler(db))
				r.Post("/award", awardXPHandler(db))
			})

			r.Route("/levels", func(r chi.Router) {
				r.Get("/", getLevelsHandler(db))
				r.Get("/current", getCurrentLevelHandler(db))
				r.Get("/{id}/next", func(w http.ResponseWriter, r *http.Request) {
					// TODO: Implement next level handler
					writeJSONError(w, http.StatusNotImplemented, "next level not yet implemented")
				})
			})

			r.Route("/badges", func(r chi.Router) {
				r.Get("/", getBadgesHandler(db))
				r.Get("/{id}", getBadgeHandler(db))
				r.Get("/me", getUserBadgesHandler(db))
			})

			r.Route("/streaks", func(r chi.Router) {
				r.Get("/", getStreakHandler(db))
				r.Post("/log", logStreakHandler(db))
			})

			r.Route("/spin-wheel", func(r chi.Router) {
				r.Get("/", getSpinWheelHandler(db))
				r.Post("/spin", spinWheelHandler(db))
				r.Get("/history", getSpinHistoryHandler(db))
			})

			// Engagement routes
			r.Route("/flash-challenges", func(r chi.Router) {
				r.Get("/active", getActiveFlashChallengesHandler(db))
				r.Post("/{id}/participate", participateFlashChallengeHandler(db))
			})

			r.Route("/trivia", func(r chi.Router) {
				r.Get("/active", getActiveTriviaHandler(db))
				r.Post("/{id}/start", startTriviaHandler(db))
				r.Post("/{id}/submit-answers", submitTriviaAnswersHandler(db))
			})

			r.Route("/mystery-boxes", func(r chi.Router) {
				r.Get("/", getMysteryBoxesHandler(db))
				r.Post("/{id}/open", openMysteryBoxHandler(db))
			})

			r.Route("/secret-codes", func(r chi.Router) {
				r.Post("/redeem/{code}", redeemSecretCodeHandler(db))
			})

			r.Route("/weekly-challenge", func(r chi.Router) {
				r.Get("/current", getWeeklyVibeChallengeHandler(db))
				r.Post("/submit", submitWeeklyVibeHandler(db))
				r.Get("/submissions", getWeeklyChallengeSubmissionsHandler(db))
				r.Post("/vote/{submissionId}", voteWeeklyChallengeHandler(db))
			})

			r.Route("/battles", func(r chi.Router) {
				r.Get("/active", getActiveBattlesHandler(db))
				r.Post("/{id}/submit", submitBattleHandler(db))
				r.Post("/{id}/vote/{submissionId}", voteBattleHandler(db))
			})

			// Rewards routes
			r.Route("/rewards", func(r chi.Router) {
				r.Get("/", getRewardsHandler(db))
				r.Get("/{id}", getRewardHandler(db))
				r.Post("/{id}/redeem", redeemRewardHandler(db))
				r.Get("/redemptions", getRewardRedemptionsHandler(db))
			})

			// Referral routes
			r.Route("/referrals", func(r chi.Router) {
				r.Get("/", getReferralsHandler(db))
				r.Get("/code", getReferralCodeHandler(db))
				r.Get("/invites", getReferralInvitesHandler(db))
				r.Post("/invite", sendReferralInviteHandler(db))
			})

			// College & State routes
			r.Route("/colleges", func(r chi.Router) {
				r.Get("/{id}", getCollegeHandler(db))
				r.Get("/{id}/leaderboard", getCollegeLeaderboardHandler(db))
			})

			r.Route("/states", func(r chi.Router) {
				r.Get("/{id}", getStateHandler(db))
				r.Get("/{id}/leaderboard", getStateLeaderboardHandler(db))
			})

			// Campus Wars routes
			r.Route("/wars", func(r chi.Router) {
				r.Get("/active", getActiveWarsHandler(db))
				r.Get("/{id}", getWarHandler(db))
				r.Get("/{id}/participants", getWarParticipantsHandler(db))
				r.Get("/{id}/leaderboard", getWarLeaderboardHandler(db))
			})

			// Survey routes
			r.Route("/surveys", func(r chi.Router) {
				r.Get("/available", getAvailableSurveysHandler(db))
				r.Get("/{id}", getSurveyHandler(db))
				r.Post("/{id}/submit", submitSurveyHandler(db))
				r.Get("/responses", getSurveyResponsesHandler(db))
			})

			// Notification routes
			r.Route("/notifications", func(r chi.Router) {
				r.Get("/", getNotificationsHandler(db))
				r.Get("/unread-count", getUnreadNotificationsCountHandler(db))
				r.Put("/{id}/read", markNotificationReadHandler(db))
				r.Put("/read-all", markAllNotificationsReadHandler(db))
				r.Delete("/{id}", deleteNotificationHandler(db))
			})

			// Wallet routes
			r.Route("/wallet", func(r chi.Router) {
				r.Get("/", getWalletHandler(db))
				r.Get("/transactions", getWalletTransactionsHandler(db))
				r.Post("/transfer", transferWalletHandler(db))
			})

			// Social & Feed routes
			r.Route("/feed", func(r chi.Router) {
				r.Get("/", getActivityFeedHandler(db))
			})

			r.Route("/posts", func(r chi.Router) {
				r.Post("/", createPostHandler(db))
				r.Post("/{id}/like", likePostHandler(db))
				r.Post("/{id}/unlike", unlikePostHandler(db))
				r.Post("/{id}/comment", commentPostHandler(db))
				r.Get("/{id}/comments", getPostCommentsHandler(db))
			})

			// Activity routes
			r.Route("/activities", func(r chi.Router) {
				r.Get("/", getUserActivityHandler(db))
				r.Get("/global", getGlobalActivityFeedHandler(db))
			})

			// Dashboard routes
			r.Route("/dashboard", func(r chi.Router) {
				r.Get("/", getUserDashboardStatsHandler(db))
				r.Get("/quick-stats", getUserDashboardStatsHandler(db))
			})

			// Email preferences
			r.Route("/email", func(r chi.Router) {
				r.Get("/preferences", getEmailPreferencesHandler(db))
				r.Put("/preferences", updateEmailPreferencesHandler(db))
				r.Post("/verify/resend", resendVerificationEmailHandler(db))
			})
		})

		// Admin routes
		r.Route("/admin", func(r chi.Router) {
			r.Use(RequireAuth(db))
			r.Use(RequireAdmin(db))

			// User management
			r.Route("/users", func(r chi.Router) {
				r.Get("/", adminGetUsersHandler(db))
				r.Post("/", adminCreateUserHandler(db))
				r.Get("/{id}", adminGetUserHandler(db))
				r.Put("/{id}", adminUpdateUserHandler(db))
				r.Delete("/{id}", adminDeleteUserHandler(db))
				r.Post("/{id}/block", blockUserHandler(db))
				r.Post("/{id}/reset-password", adminResetPasswordHandler(db))
			})

			// Task management
			r.Route("/tasks", func(r chi.Router) {
				r.Post("/", func(w http.ResponseWriter, r *http.Request) {
					writeJSONError(w, http.StatusNotImplemented, "create task not yet implemented")
				})
				r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
					writeJSONError(w, http.StatusNotImplemented, "update task not yet implemented")
				})
				r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
					writeJSONError(w, http.StatusNotImplemented, "delete task not yet implemented")
				})
			})

			// Campaign management
			r.Route("/campaigns", func(r chi.Router) {
				r.Post("/", func(w http.ResponseWriter, r *http.Request) {
					writeJSONError(w, http.StatusNotImplemented, "create campaign not yet implemented")
				})
				r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
					writeJSONError(w, http.StatusNotImplemented, "update campaign not yet implemented")
				})
				r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
					writeJSONError(w, http.StatusNotImplemented, "delete campaign not yet implemented")
				})
			})

			// Submission review
			r.Route("/submissions", func(r chi.Router) {
				r.Get("/pending", adminGetPendingSubmissionsHandler(db))
				r.Get("/stats", adminGetSubmissionStatsHandler(db))
				r.Post("/{id}/review", adminReviewSubmissionHandler(db))
			})

			// Gamification management
			r.Route("/xp", func(r chi.Router) {
				r.Post("/award", adminAwardXPHandler(db))
				r.Post("/penalize", adminPenalizeXPHandler(db))
			})

			r.Route("/badges", func(r chi.Router) {
				r.Post("/award", adminAwardBadgeHandler(db))
			})

			// Dashboard & Analytics
			r.Route("/dashboard", func(r chi.Router) {
				r.Get("/", adminDashboardHandler(db))
			})

			r.Route("/analytics", func(r chi.Router) {
				r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
					writeJSONError(w, http.StatusNotImplemented, "user analytics not yet implemented")
				})
				r.Get("/engagement", func(w http.ResponseWriter, r *http.Request) {
					writeJSONError(w, http.StatusNotImplemented, "engagement analytics not yet implemented")
				})
			})
		})
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
