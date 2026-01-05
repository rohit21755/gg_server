package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rohit21755/gg_server.git/internal/store"
	"gorm.io/gorm"
)

// Get Notifications
func getNotificationsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Parse query parameters
		notificationType := r.URL.Query().Get("type")
		isRead := r.URL.Query().Get("is_read")
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
		query := db.Model(&store.Notification{}).Where("user_id = ?", user.ID)

		if notificationType != "" {
			query = query.Where("notification_type = ?", notificationType)
		}

		if isRead != "" {
			query = query.Where("is_read = ?", isRead == "true")
		}

		// Get total count
		var totalCount int64
		query.Count(&totalCount)

		// Get notifications
		var notifications []store.Notification
		result := query.
			Order("sent_at DESC, created_at DESC").
			Offset(offset).
			Limit(limit).
			Find(&notifications)

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		// Format response
		var responseNotifications []map[string]interface{}
		for _, notification := range notifications {
			notificationData := map[string]interface{}{
				"id":                notification.ID,
				"notification_type": notification.NotificationType,
				"title":             notification.Title,
				"message":           notification.Message,
				"is_read":           notification.IsRead,
				"is_actionable":     notification.IsActionable,
				"action_url":        notification.ActionURL,
				"sent_at":           notification.SentAt,
				"read_at":           notification.ReadAt,
				"created_at":        notification.CreatedAt,
			}

			// Parse data JSON
			if notification.Data != nil && *notification.Data != "" {
				var data interface{}
				if err := json.Unmarshal([]byte(*notification.Data), &data); err == nil {
					notificationData["data"] = data
				}
			}

			responseNotifications = append(responseNotifications, notificationData)
		}

		// Get unread count
		var unreadCount int64
		db.Model(&store.Notification{}).
			Where("user_id = ? AND is_read = ?", user.ID, false).
			Count(&unreadCount)

		response := map[string]interface{}{
			"notifications": responseNotifications,
			"unread_count":  unreadCount,
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

// Get Unread Notifications Count
func getUnreadNotificationsCountHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		var unreadCount int64
		db.Model(&store.Notification{}).
			Where("user_id = ? AND is_read = ?", user.ID, false).
			Count(&unreadCount)

		response := map[string]interface{}{
			"unread_count": unreadCount,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Mark Notification as Read
func markNotificationReadHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		notificationIDStr := chi.URLParam(r, "id")
		notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid notification ID"))
			return
		}

		var notification store.Notification
		result := db.First(&notification, notificationID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("notification not found"))
			} else {
				internalServerError(w, r, result.Error)
			}
			return
		}

		// Check ownership
		if notification.UserID == nil || *notification.UserID != int(user.ID) {
			unauthorizedResponse(w, r, errors.New("not authorized to update this notification"))
			return
		}

		// Mark as read
		if !notification.IsRead {
			now := time.Now()
			notification.IsRead = true
			notification.ReadAt = &now

			if err := db.Save(&notification).Error; err != nil {
				internalServerError(w, r, err)
				return
			}
		}

		response := map[string]interface{}{
			"message": "Notification marked as read",
			"notification": map[string]interface{}{
				"id":      notification.ID,
				"is_read": notification.IsRead,
				"read_at": notification.ReadAt,
			},
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Mark All Notifications as Read
func markAllNotificationsReadHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		// Get count of unread notifications before update
		var unreadCountBefore int64
		db.Model(&store.Notification{}).
			Where("user_id = ? AND is_read = ?", user.ID, false).
			Count(&unreadCountBefore)

		// Mark all as read
		now := time.Now()
		result := db.Model(&store.Notification{}).
			Where("user_id = ? AND is_read = ?", user.ID, false).
			Updates(map[string]interface{}{
				"is_read": true,
				"read_at": now,
			})

		if result.Error != nil {
			internalServerError(w, r, result.Error)
			return
		}

		response := map[string]interface{}{
			"message":               "All notifications marked as read",
			"notifications_updated": result.RowsAffected,
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}

// Delete Notification
func deleteNotificationHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			unauthorizedResponse(w, r, errors.New("user not found in context"))
			return
		}

		notificationIDStr := chi.URLParam(r, "id")
		notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid notification ID"))
			return
		}

		var notification store.Notification
		result := db.First(&notification, notificationID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				notFoundResponse(w, r, errors.New("notification not found"))
			} else {
				internalServerError(w, r, result.Error)
			}
			return
		}

		// Check ownership
		if notification.UserID == nil || *notification.UserID != int(user.ID) {
			unauthorizedResponse(w, r, errors.New("not authorized to delete this notification"))
			return
		}

		// Delete notification
		if err := db.Delete(&notification).Error; err != nil {
			internalServerError(w, r, err)
			return
		}

		response := map[string]interface{}{
			"message": "Notification deleted successfully",
		}

		if err := jsonResponse(w, http.StatusOK, response); err != nil {
			internalServerError(w, r, err)
		}
	}
}
