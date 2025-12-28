package store

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID              uint       `gorm:"primaryKey"`
	UserID          *int       `gorm:"index;constraint:OnDelete:CASCADE"`
	NotificationType string    `gorm:"size:50;not null;check:notification_type IN ('task_assigned', 'submission_status', 'reward_unlocked', 'level_up', 'streak_update', 'new_challenge', 'winner_announcement', 'system')"`
	Title           string     `gorm:"size:200;not null"`
	Message         string     `gorm:"type:text;not null"`
	Data            *string    `gorm:"type:jsonb;default:'{}'"`
	IsRead          bool       `gorm:"default:false"`
	IsActionable    bool       `gorm:"default:false"`
	ActionURL       *string    `gorm:"type:text"`
	ScheduledFor    *time.Time `gorm:"type:timestamp"`
	SentAt          time.Time  `gorm:"autoCreateTime"`
	ReadAt          *time.Time `gorm:"type:timestamp"`
	CreatedAt       time.Time  `gorm:"autoCreateTime"`

	// Relations
	User *User `gorm:"foreignKey:UserID"`
}

func (Notification) TableName() string {
	return "notifications"
}

func CreateNotification(db *gorm.DB, notification *Notification) error {
	return db.Create(notification).Error
}

func GetNotificationByID(db *gorm.DB, id uint) (*Notification, error) {
	var notification Notification
	if err := db.First(&notification, id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}
