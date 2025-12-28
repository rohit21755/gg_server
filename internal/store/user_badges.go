package store

import (
	"time"

	"gorm.io/gorm"
)

type UserBadge struct {
	ID       uint      `gorm:"primaryKey"`
	UserID   int       `gorm:"index;constraint:OnDelete:CASCADE"`
	BadgeID  int       `gorm:"index"`
	EarnedAt time.Time `gorm:"autoCreateTime"`

	// Relations
	User  *User  `gorm:"foreignKey:UserID"`
	Badge *Badge `gorm:"foreignKey:BadgeID"`
}

func (UserBadge) TableName() string {
	return "user_badges"
}

func CreateUserBadge(db *gorm.DB, userBadge *UserBadge) error {
	return db.Create(userBadge).Error
}

func GetUserBadgeByID(db *gorm.DB, id uint) (*UserBadge, error) {
	var userBadge UserBadge
	if err := db.First(&userBadge, id).Error; err != nil {
		return nil, err
	}
	return &userBadge, nil
}

func GetUserBadges(db *gorm.DB, userID uint) ([]UserBadge, error) {
	var userBadges []UserBadge
	userIDInt := int(userID)
	if err := db.Where("user_id = ?", userIDInt).Find(&userBadges).Error; err != nil {
		return nil, err
	}
	return userBadges, nil
}
