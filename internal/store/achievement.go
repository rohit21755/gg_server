package store

import (
	"time"

	"gorm.io/gorm"
)

// Achievement represents an achievement definition
type Achievement struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	IconURL     *string   `gorm:"type:text" json:"icon_url,omitempty"`
	Category    string    `gorm:"size:50" json:"category"` // submission, streak, xp, referral, etc.
	XP          int       `gorm:"default:0" json:"xp"`
	Coins       int       `gorm:"default:0" json:"coins"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Achievement) TableName() string { return "achievements" }

// UserAchievement represents a user's earned achievement
type UserAchievement struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	AchievementID uint      `gorm:"not null;index" json:"achievement_id"`
	EarnedAt      time.Time `gorm:"autoCreateTime" json:"earned_at"`
}

func (UserAchievement) TableName() string { return "user_achievements" }

// GORM helper functions
func GetAchievement(db *gorm.DB, id uint) (*Achievement, error) {
	var achievement Achievement
	if err := db.First(&achievement, id).Error; err != nil {
		return nil, err
	}
	return &achievement, nil
}

func GetAllAchievements(db *gorm.DB) ([]Achievement, error) {
	var achievements []Achievement
	if err := db.Where("is_active = ?", true).Find(&achievements).Error; err != nil {
		return nil, err
	}
	return achievements, nil
}

func GetUserAchievements(db *gorm.DB, userID uint) ([]UserAchievement, error) {
	var userAchievements []UserAchievement
	if err := db.Where("user_id = ?", userID).
		Preload("Achievement").
		Order("earned_at DESC").
		Find(&userAchievements).Error; err != nil {
		return nil, err
	}
	return userAchievements, nil
}

func AwardAchievement(db *gorm.DB, userID, achievementID uint) error {
	// Check if already awarded
	var existing UserAchievement
	err := db.Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		First(&existing).Error
	if err == nil {
		return nil // Already awarded
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Award achievement
	userAchievement := UserAchievement{
		UserID:        userID,
		AchievementID: achievementID,
	}
	return db.Create(&userAchievement).Error
}

func CreateAchievement(db *gorm.DB, achievement *Achievement) error {
	return db.Create(achievement).Error
}

func UpdateAchievement(db *gorm.DB, achievement *Achievement) error {
	return db.Save(achievement).Error
}

func DeleteAchievement(db *gorm.DB, id uint) error {
	return db.Delete(&Achievement{}, id).Error
}

