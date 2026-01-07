package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// Get activity feed
func getActivityFeedHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

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

		posts, err := store.GetSocialFeed(db, uint(user.ID), limit, offset)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch feed")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"feed": posts,
		})
	}
}

// Create social post
func createPostHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		var req struct {
			Content   string `json:"content"`
			MediaURLs string `json:"media_urls,omitempty"`
			PostType  string `json:"post_type"`
			IsPublic  bool   `json:"is_public"`
		}

		if err := readJSON(w, r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		post := &store.SocialPost{
			UserID:    uint(user.ID),
			Content:   req.Content,
			PostType:  req.PostType,
			IsPublic:  req.IsPublic,
		}

		if req.MediaURLs != "" {
			post.MediaURLs = &req.MediaURLs
		}

		if err := store.CreateSocialPost(db, post); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to create post")
			return
		}

		writeJSON(w, http.StatusCreated, post)
	}
}

// Like post
func likePostHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		postIDStr := chi.URLParam(r, "id")
		postID, err := strconv.ParseUint(postIDStr, 10, 32)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid post ID")
			return
		}

		if err := store.LikePost(db, uint(postID), uint(user.ID)); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to like post")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"message": "post liked"})
	}
}

// Unlike post
func unlikePostHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		postIDStr := chi.URLParam(r, "id")
		postID, err := strconv.ParseUint(postIDStr, 10, 32)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid post ID")
			return
		}

		if err := store.UnlikePost(db, uint(postID), uint(user.ID)); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to unlike post")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"message": "post unliked"})
	}
}

// Comment on post
func commentPostHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		postIDStr := chi.URLParam(r, "id")
		postID, err := strconv.ParseUint(postIDStr, 10, 32)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid post ID")
			return
		}

		var req struct {
			Content string `json:"content"`
		}

		if err := readJSON(w, r, &req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		comment := &store.PostComment{
			PostID:  uint(postID),
			UserID:  uint(user.ID),
			Content: req.Content,
		}

		if err := store.CreatePostComment(db, comment); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to create comment")
			return
		}

		writeJSON(w, http.StatusCreated, comment)
	}
}

// Get post comments
func getPostCommentsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := chi.URLParam(r, "id")
		postID, err := strconv.ParseUint(postIDStr, 10, 32)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid post ID")
			return
		}

		limit := 50
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		comments, err := store.GetPostComments(db, uint(postID), limit)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch comments")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"comments": comments,
		})
	}
}

// Get global activity feed
func getGlobalActivityFeedHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 50
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		activities, err := store.GetGlobalActivityLogs(db, limit)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch activities")
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"activities": activities,
		})
	}
}

