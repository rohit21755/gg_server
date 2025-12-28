package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/gorm"
)

func setupREST(r chi.Router, db *gorm.DB) {
	log.Println("Setting up REST API")

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Get("/health", healthHandler)

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
			// r.Route("/rewards", func(r chi.Router) {
			// 	r.Get("/", getRewardsHandler(db))
			// 	r.Get("/{id}", getRewardHandler(db))
			// 	r.Post("/{id}/redeem", redeemRewardHandler(db))
			// 	r.Get("/redemptions", getRewardRedemptionsHandler(db))
			// })

			// Referral routes
			// r.Route("/referrals", func(r chi.Router) {
			// 	r.Get("/", getReferralsHandler(db))
			// 	r.Get("/code", getReferralCodeHandler(db))
			// 	r.Get("/invites", getReferralInvitesHandler(db))
			// 	r.Post("/invite", sendReferralInviteHandler(db))
			// })

			// College & State routes
			// r.Route("/colleges", func(r chi.Router) {
			// 	r.Get("/", getCollegesHandler(db))
			// 	r.Get("/{id}", getCollegeHandler(db))
			// 	r.Get("/{id}/stats", getCollegeStatsHandler(db))
			// })

			// r.Route("/states", func(r chi.Router) {
			// 	r.Get("/", getStatesHandler(db))
			// 	r.Get("/{id}", getStateHandler(db))
			// 	r.Get("/{id}/leaderboard", getStateLeaderboardHandler(db))
			// })

			// Campus Wars routes
			// r.Route("/wars", func(r chi.Router) {
			// 	r.Get("/active", getActiveWarsHandler(db))
			// 	r.Get("/{id}", getWarHandler(db))
			// 	r.Get("/{id}/participants", getWarParticipantsHandler(db))
			// 	r.Get("/{id}/leaderboard", getWarLeaderboardHandler(db))
			// })

			// Survey routes
			// r.Route("/surveys", func(r chi.Router) {
			// 	r.Get("/available", getAvailableSurveysHandler(db))
			// 	r.Get("/{id}", getSurveyHandler(db))
			// 	r.Post("/{id}/submit", submitSurveyHandler(db))
			// 	r.Get("/responses", getSurveyResponsesHandler(db))
			// })

			// Notification routes
			// r.Route("/notifications", func(r chi.Router) {
			// 	r.Get("/", getNotificationsHandler(db))
			// 	r.Get("/unread-count", getUnreadNotificationsCountHandler(db))
			// 	r.Put("/{id}/read", markNotificationReadHandler(db))
			// 	r.Put("/read-all", markAllNotificationsReadHandler(db))
			// 	r.Delete("/{id}", deleteNotificationHandler(db))
			// })
		})
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
