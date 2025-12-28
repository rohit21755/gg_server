package store

import (
	"time"

	"gorm.io/gorm"
)

type UserStreak struct {
	ID             uint      `gorm:"primaryKey"`
	UserID         *int      `gorm:"unique;index;constraint:OnDelete:CASCADE"`
	StreakType     string    `gorm:"size:50;not null;check:streak_type IN ('daily_login', 'weekly_task', 'campaign')"`
	CurrentStreak  int       `gorm:"default:0"`
	LongestStreak  int       `gorm:"default:0"`
	LastActivityDate time.Time `gorm:"type:date;not null"`
	TotalDays      int       `gorm:"default:0"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`

	// Relations
	User *User `gorm:"foreignKey:UserID"`
}

func (UserStreak) TableName() string {
	return "user_streaks"
}

type StreakLog struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	StreakType   string    `gorm:"size:50;not null"`
	ActivityDate time.Time `gorm:"type:date;not null"`
	EarnedXP     int       `gorm:"default:0"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	// Relations
	User *User `gorm:"foreignKey:UserID"`
}

func (StreakLog) TableName() string {
	return "streak_logs"
}

func CreateUserStreak(db *gorm.DB, streak *UserStreak) error {
	return db.Create(streak).Error
}

func GetUserStreakByID(db *gorm.DB, id uint) (*UserStreak, error) {
	var streak UserStreak
	if err := db.First(&streak, id).Error; err != nil {
		return nil, err
	}
	return &streak, nil
}

func CreateStreakLog(db *gorm.DB, log *StreakLog) error {
	return db.Create(log).Error
}

func GetStreakLogByID(db *gorm.DB, id uint) (*StreakLog, error) {
	var log StreakLog
	if err := db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}
