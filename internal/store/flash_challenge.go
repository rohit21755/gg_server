package store

import (
	"time"

	"gorm.io/gorm"
)

type FlashChallenge struct {
	ID             uint       `gorm:"primaryKey"`
	Title          string     `gorm:"size:200;not null"`
	Description    *string    `gorm:"type:text"`
	ChallengeType  *string    `gorm:"size:50;check:challenge_type IN ('meme', 'reel', 'video', 'qr_scan', 'content')"`
	DurationHours  int        `gorm:"not null"`
	XPReward       int        `gorm:"not null"`
	MaxParticipants *int      `gorm:"type:integer"`
	StartTime      time.Time  `gorm:"not null"`
	EndTime        time.Time  `gorm:"not null"`
	Status         string     `gorm:"size:20;default:'scheduled';check:status IN ('scheduled', 'active', 'completed')"`
	CreatedBy      *int       `gorm:"index"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`

	// Relations
	Creator *User `gorm:"foreignKey:CreatedBy"`
}

func (FlashChallenge) TableName() string {
	return "flash_challenges"
}

func CreateFlashChallenge(db *gorm.DB, challenge *FlashChallenge) error {
	return db.Create(challenge).Error
}

func GetFlashChallengeByID(db *gorm.DB, id uint) (*FlashChallenge, error) {
	var challenge FlashChallenge
	if err := db.First(&challenge, id).Error; err != nil {
		return nil, err
	}
	return &challenge, nil
}
