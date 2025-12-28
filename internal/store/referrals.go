package store

import (
	"time"

	"gorm.io/gorm"
)

type Referral struct {
	ID                  uint      `gorm:"primaryKey"`
	ReferrerID          *uint     `gorm:"index;constraint:OnDelete:CASCADE"` // FK -> users.id (cascade)
	ReferredEmail       string    `gorm:"size:255;not null;uniqueIndex:idx_referrer_email"`
	ReferredUserID      *uint     `gorm:"index"` // FK -> users.id
	Status              string    `gorm:"size:20;default:'pending';check:status IN ('pending','joined','completed_task','converted')"`
	XPAwarded           int       `gorm:"default:0"`
	XPAwardedToReferred int       `gorm:"default:0"`
	ConversionStage     int       `gorm:"default:1"`
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`

	// Relations
	Referrer     *User `gorm:"foreignKey:ReferrerID"`     // User who referred
	ReferredUser *User `gorm:"foreignKey:ReferredUserID"` // User who signed up
}

// Table name overrides default naming
func (Referral) TableName() string {
	return "referrals"
}

func CreateReferral(db *gorm.DB, referral *Referral) error {
	return db.Create(referral).Error
}
