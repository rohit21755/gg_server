package store

import (
	"time"

	"gorm.io/gorm"
)

type Level struct {
	ID               uint      `gorm:"primaryKey"`
	Name             string    `gorm:"size:50;unique;not null"`
	RankOrder        int       `gorm:"unique;not null"`
	MinXP            int       `gorm:"not null"`
	MaxXP            *int      `gorm:"type:integer"`
	BadgeURL         *string   `gorm:"type:text"`
	Description      *string   `gorm:"type:text"`
	UnlockableFeatures *string `gorm:"type:jsonb;default:'[]'"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}

func (Level) TableName() string {
	return "levels"
}

type Badge struct {
	ID               uint       `gorm:"primaryKey"`
	Name             string     `gorm:"size:100;not null"`
	Description      *string    `gorm:"type:text"`
	BadgeType        string     `gorm:"size:50;not null"`
	Category         *string    `gorm:"size:50;check:category IN ('referral', 'submission', 'campaign', 'streak', 'special')"`
	ImageURL         string     `gorm:"type:text;not null"`
	XPReward         int        `gorm:"default:0"`
	CriteriaType     string     `gorm:"size:50;not null"`
	CriteriaValue    int        `gorm:"not null"`
	IsSecret         bool       `gorm:"default:false"`
	IsLimitedEdition bool       `gorm:"default:false"`
	AvailableUntil   *time.Time `gorm:"type:timestamp"`
	CreatedAt        time.Time  `gorm:"autoCreateTime"`
}

func (Badge) TableName() string {
	return "badges"
}

type ProfileSkin struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:100;not null"`
	Description *string   `gorm:"type:text"`
	PreviewURL  string    `gorm:"type:text;not null"`
	UnlockMethod string   `gorm:"size:50;check:unlock_method IN ('xp', 'campaign', 'purchase', 'special')"`
	XPCost      int       `gorm:"default:0"`
	CampaignID  *int      `gorm:"index"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	// Relations
	Campaign *Campaign `gorm:"foreignKey:CampaignID"`
}

func (ProfileSkin) TableName() string {
	return "profile_skins"
}

func CreateLevel(db *gorm.DB, level *Level) error {
	return db.Create(level).Error
}

func GetLevelByID(db *gorm.DB, id uint) (*Level, error) {
	var level Level
	if err := db.First(&level, id).Error; err != nil {
		return nil, err
	}
	return &level, nil
}

func CreateBadge(db *gorm.DB, badge *Badge) error {
	return db.Create(badge).Error
}

func GetBadgeByID(db *gorm.DB, id uint) (*Badge, error) {
	var badge Badge
	if err := db.First(&badge, id).Error; err != nil {
		return nil, err
	}
	return &badge, nil
}

func CreateProfileSkin(db *gorm.DB, skin *ProfileSkin) error {
	return db.Create(skin).Error
}

func GetProfileSkinByID(db *gorm.DB, id uint) (*ProfileSkin, error) {
	var skin ProfileSkin
	if err := db.First(&skin, id).Error; err != nil {
		return nil, err
	}
	return &skin, nil
}
