package store

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog represents admin action logging
type AuditLog struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminID     uint      `gorm:"not null;index" json:"admin_id"`
	Action      string    `gorm:"size:100;not null" json:"action"`     // create, update, delete, etc.
	EntityType  string    `gorm:"size:50;not null" json:"entity_type"` // user, task, campaign, etc.
	EntityID    *uint     `json:"entity_id,omitempty"`
	Description string    `gorm:"type:text" json:"description"`
	IPAddress   *string   `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent   *string   `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_logs" }

// EmailTemplate represents email templates
type EmailTemplate struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:255;unique;not null" json:"name"`
	Subject   string    `gorm:"size:500;not null" json:"subject"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	Variables *string   `gorm:"type:text" json:"variables,omitempty"` // JSON array of variable names
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (EmailTemplate) TableName() string { return "email_templates" }

// UserEmailPreferences represents user email notification preferences
type UserEmailPreferences struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	MarketingEmails   bool      `gorm:"default:true" json:"marketing_emails"`
	TaskNotifications bool      `gorm:"default:true" json:"task_notifications"`
	AchievementEmails bool      `gorm:"default:true" json:"achievement_emails"`
	WeeklyDigest      bool      `gorm:"default:true" json:"weekly_digest"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserEmailPreferences) TableName() string { return "user_email_preferences" }

// GORM helper functions
func CreateAuditLog(db *gorm.DB, log *AuditLog) error {
	return db.Create(log).Error
}

func GetAuditLogs(db *gorm.DB, limit int, offset int) ([]AuditLog, error) {
	var logs []AuditLog
	query := db.Order("created_at DESC").Limit(limit).Offset(offset)
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func GetEmailTemplate(db *gorm.DB, name string) (*EmailTemplate, error) {
	var template EmailTemplate
	if err := db.Where("name = ? AND is_active = ?", name, true).First(&template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func GetUserEmailPreferences(db *gorm.DB, userID uint) (*UserEmailPreferences, error) {
	var prefs UserEmailPreferences
	err := db.Where("user_id = ?", userID).First(&prefs).Error
	if err == gorm.ErrRecordNotFound {
		// Create default preferences
		prefs = UserEmailPreferences{
			UserID:            userID,
			MarketingEmails:   true,
			TaskNotifications: true,
			AchievementEmails: true,
			WeeklyDigest:      true,
		}
		if err := db.Create(&prefs).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &prefs, nil
}

func UpdateUserEmailPreferences(db *gorm.DB, prefs *UserEmailPreferences) error {
	return db.Save(prefs).Error
}
