package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rohit21755/gg_server.git/internal/env"
	"github.com/rohit21755/gg_server.git/internal/store"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Request/Response structs
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Phone      string `json:"phone"`
	CollegeID  uint   `json:"college_id"`
	ReferralID string `json:"referral_id"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// JWT Claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Generate JWT token
func generateToken(user *store.User) (string, string, error) {
	// Access token
	accessTokenClaims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "campus-ambassador",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(env.Get("JWT_SECRET", "1234567890")))
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshTokenClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "campus-ambassador-refresh",
		Subject:   fmt.Sprintf("%d", user.ID),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(env.Get("JWT_REFRESH", "1234567890")))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// Login Handler
func loginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Validate request
		if err := Validate.Struct(req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Find user by email
		user, err := store.GetUserByEmail(db, req.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				unauthorizedResponse(w, r, errors.New("invalid credentials"))
			} else {
				internalServerError(w, r, err)
			}
			return
		}

		// Check password
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			unauthorizedResponse(w, r, errors.New("invalid credentials"))
			return
		}

		// Check if user is active
		if !user.IsActive {
			unauthorizedResponse(w, r, errors.New("account is deactivated"))
			return
		}

		// Generate tokens
		accessToken, refreshToken, err := generateToken(user)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		// Update last login
		now := time.Now()
		user.LastLoginDate = &now
		if err := store.UpdateUser(db, user); err != nil {
			internalServerError(w, r, err)
			return
		}

		// Create session
		userID := int(user.ID)
		deviceID := r.Header.Get("X-Device-ID")
		platform := r.Header.Get("X-Platform")
		session := store.UserSession{
			UserID:       &userID,
			SessionToken: accessToken,
			DeviceID:     stringPtr(deviceID),
			Platform:     stringPtr(platform),
			LastActive:   time.Now(),
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}
		if err := store.CreateSession(db, &session); err != nil {
			internalServerError(w, r, err)
			return
		}

		// Update streak
		if err := updateUserStreak(db, user.ID, "daily_login"); err != nil {
			// Log but don't fail login
			fmt.Printf("Error updating streak: %v\n", err)
		}

		// Response
		response := map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    24 * 3600,
			"user": map[string]interface{}{
				"id":         user.ID,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"role":       user.Role,
				"xp":         user.XP,
				"college_id": user.CollegeID,
			},
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Register Handler
func registerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Validate request
		if err := Validate.Struct(req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Check if user already exists
		_, err := store.GetUserByEmail(db, req.Email)
		if err == nil {
			conflictResponse(w, r, errors.New("user with this email already exists"))
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		// Generate referral code
		referralCode := generateReferralCode()

		// Check referral
		var referrerID *uint
		if req.ReferralID != "" {
			referrer, err := store.GetUserByReferralCode(db, req.ReferralID)
			if err == nil {
				referrerID = &referrer.ID
			}
		}

		// Convert types for User struct
		var collegeID *int
		if req.CollegeID != 0 {
			id := int(req.CollegeID)
			collegeID = &id
		}
		var referredBy *int
		if referrerID != nil {
			id := int(*referrerID)
			referredBy = &id
		}
		levelID := 1

		// Create user
		user := &store.User{
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			Phone:        stringPtr(req.Phone),
			CollegeID:    collegeID,
			ReferralCode: referralCode,
			ReferredBy:   referredBy,
			Role:         "ca",
			XP:           100,      // Starting XP
			LevelID:      &levelID, // Rookie level
			IsActive:     true,
		}

		if err := store.CreateUser(db, user); err != nil {
			internalServerError(w, r, err)
			return
		}

		// Create referral record if referred
		if referrerID != nil {
			referral := &store.Referral{
				ReferrerID:     referrerID,
				ReferredEmail:  req.Email,
				ReferredUserID: &user.ID,
				Status:         "joined",
			}
			if err := store.CreateReferral(db, referral); err != nil {
				fmt.Printf("Error creating referral: %v\n", err)
			}

			// Award XP to referrer
			userID := int(*referrerID)
			sourceID := int(referral.ID)
			xpTransaction := &store.XPTransaction{
				UserID:          &userID,
				TransactionType: "referral",
				Amount:          500,
				SourceType:      stringPtr("referral"),
				SourceID:        &sourceID,
				Description:     stringPtr("Referral bonus"),
			}
			if err := store.CreateXPTransaction(db, xpTransaction); err != nil {
				fmt.Printf("Error awarding referral XP: %v\n", err)
			}
		}

		// Generate tokens
		accessToken, refreshToken, err := generateToken(user)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		// Create session
		userID := int(user.ID)
		deviceID := r.Header.Get("X-Device-ID")
		platform := r.Header.Get("X-Platform")
		session := store.UserSession{
			UserID:       &userID,
			SessionToken: accessToken,
			DeviceID:     stringPtr(deviceID),
			Platform:     stringPtr(platform),
			LastActive:   time.Now(),
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}
		if err := store.CreateSession(db, &session); err != nil {
			internalServerError(w, r, err)
			return
		}

		// Response
		response := map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    24 * 3600,
			"user": map[string]interface{}{
				"id":            user.ID,
				"email":         user.Email,
				"first_name":    user.FirstName,
				"last_name":     user.LastName,
				"role":          user.Role,
				"xp":            user.XP,
				"referral_code": user.ReferralCode,
				"college_id":    user.CollegeID,
			},
		}

		if err := jsonResponse(w, http.StatusCreated, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Refresh Token Handler
func refreshTokenHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RefreshTokenRequest
		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Validate refresh token
		token, err := jwt.ParseWithClaims(req.RefreshToken, &jwt.RegisteredClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(env.Get("JWT_REFRESH", "1234567890")), nil
			})

		if err != nil || !token.Valid {
			unauthorizedResponse(w, r, errors.New("invalid refresh token"))
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			unauthorizedResponse(w, r, errors.New("invalid token claims"))
			return
		}

		// Get user
		userID, err := claims.GetSubject()
		if err != nil {
			unauthorizedResponse(w, r, err)
			return
		}

		var userIDUint uint
		fmt.Sscanf(userID, "%d", &userIDUint)

		user, err := store.GetUserByID(db, userIDUint)
		if err != nil {
			unauthorizedResponse(w, r, errors.New("user not found"))
			return
		}

		// Generate new tokens
		accessToken, refreshToken, err := generateToken(user)
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		// Create new session
		userIDInt := int(userIDUint)
		deviceID := r.Header.Get("X-Device-ID")
		platform := r.Header.Get("X-Platform")
		session := store.UserSession{
			UserID:       &userIDInt,
			SessionToken: accessToken,
			DeviceID:     stringPtr(deviceID),
			Platform:     stringPtr(platform),
			LastActive:   time.Now(),
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}
		if err := store.CreateSession(db, &session); err != nil {
			internalServerError(w, r, err)
			return
		}

		response := map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    24 * 3600,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Logout Handler
func logoutHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, http.StatusUnauthorized, "authorization header is required")
			return
		}

		// Extract token
		token := extractToken(authHeader)
		if token == "" {
			unauthorizedResponse(w, r, errors.New("invalid token format"))
			return
		}

		// Delete session
		if err := store.DeleteSessionByToken(db, token); err != nil {
			internalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Successfully logged out"})
	}
}

// Forgot Password Handler
func forgotPasswordHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ForgotPasswordRequest
		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Find user
		_, err := store.GetUserByEmail(db, req.Email)
		if err != nil {
			// Don't reveal if user exists or not
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "If an account exists with this email, you will receive a password reset link",
			})
			return
		}

		// Generate reset token
		// resetToken := generateSecureToken(32)
		// resetExpiry := time.Now().Add(1 * time.Hour)

		// Store reset token (you'd create a password_reset_tokens table)
		// For now, we'll just return success

		// Send email (in production)
		// sendPasswordResetEmail(user.Email, resetToken)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Password reset email sent",
		})
	}
}

// Reset Password Handler
func resetPasswordHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ResetPasswordRequest
		if err := readJSON(w, r, &req); err != nil {
			badRequestResponse(w, r, err)
			return
		}

		// Validate token (check in password_reset_tokens table)
		// For now, we'll skip token validation

		// Get user from token (you'd look up by token)
		// userID := getUserIdFromResetToken(req.Token)

		// Hash new password
		// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		// if err != nil {
		// 	internalServerError(w, r, err)
		// 	return
		// }

		// Update password
		// user.PasswordHash = string(hashedPassword)
		// store.UpdateUser(db, user)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Password reset successful",
		})
	}
}

// Verify Email Handler
func verifyEmailHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")
		if token == "" {
			badRequestResponse(w, r, errors.New("token is required"))
			return
		}

		// Validate email verification token
		// Update user's email_verified status

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Email verified successfully",
		})
	}
}

// Helper functions
func generateReferralCode() string {
	// Generate a random 8-character referral code
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateSecureToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func extractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}

func updateUserStreak(db *gorm.DB, userID uint, streakType string) error {
	// Get or create streak
	streak, err := store.GetUserStreak(db, userID, streakType)
	if err != nil {
		// Create new streak
		userIDInt := int(userID)
		streak = &store.UserStreak{
			UserID:           &userIDInt,
			StreakType:       streakType,
			CurrentStreak:    1,
			LongestStreak:    1,
			LastActivityDate: time.Now().Truncate(24 * time.Hour),
			TotalDays:        1,
		}
		return store.CreateUserStreak(db, streak)
	}

	// Check if last activity was yesterday
	lastActivity := streak.LastActivityDate
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	// If already logged in today, don't update
	if lastActivity.Equal(today) {
		return nil
	}

	if lastActivity.Equal(yesterday) {
		// Continue streak
		streak.CurrentStreak++
		if streak.CurrentStreak > streak.LongestStreak {
			streak.LongestStreak = streak.CurrentStreak
		}
	} else if lastActivity.Before(yesterday) {
		// Break streak, start new one
		streak.CurrentStreak = 1
	}

	streak.LastActivityDate = today
	streak.TotalDays++

	return store.UpdateUserStreak(db, streak)
}

// stringPtr returns a pointer to the string if it's not empty, otherwise returns nil
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intPtr(i int) *int {
	return &i
}
