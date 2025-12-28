package store

import (
	"time"

	"gorm.io/gorm"
)

type ActivityLog struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	ActivityType string    `gorm:"size:100;not null"`
	ActivityData *string   `gorm:"type:jsonb;default:'{}'"`
	IPAddress    *string   `gorm:"type:inet"`
	UserAgent    *string   `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	// Relations
	User *User `gorm:"foreignKey:UserID"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}

func CreateActivityLog(db *gorm.DB, log *ActivityLog) error {
	return db.Create(log).Error
}

func GetActivityLogByID(db *gorm.DB, id uint) (*ActivityLog, error) {
	var log ActivityLog
	if err := db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}
