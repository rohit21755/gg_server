package store

import (
	"time"

	"gorm.io/gorm"
)

type SecretCode struct {
	ID               uint       `gorm:"primaryKey"`
	Code             string     `gorm:"size:50;unique;not null"`
	Description      *string    `gorm:"type:text"`
	XPReward         int        `gorm:"not null"`
	CoinReward       int        `gorm:"default:0"`
	BadgeID          *int       `gorm:"index"`
	MaxRedemptions   int        `gorm:"default:1"`
	CurrentRedemptions int      `gorm:"default:0"`
	ValidFrom        time.Time  `gorm:"not null"`
	ValidUntil       time.Time  `gorm:"not null"`
	DistributionChannel *string `gorm:"size:50"`
	IsActive         bool       `gorm:"default:true"`
	CreatedBy        *int       `gorm:"index"`
	CreatedAt        time.Time  `gorm:"autoCreateTime"`

	// Relations
	Badge     *Badge `gorm:"foreignKey:BadgeID"`
	Creator   *User  `gorm:"foreignKey:CreatedBy"`
}

func (SecretCode) TableName() string {
	return "secret_codes"
}

type SecretCodeRedemption struct {
	ID           uint      `gorm:"primaryKey"`
	SecretCodeID *int      `gorm:"index"`
	UserID       *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	RedeemedAt   time.Time `gorm:"autoCreateTime"`
	IPAddress    *string   `gorm:"type:inet"`

	// Relations
	SecretCode *SecretCode `gorm:"foreignKey:SecretCodeID"`
	User       *User       `gorm:"foreignKey:UserID"`
}

func (SecretCodeRedemption) TableName() string {
	return "secret_code_redemptions"
}

func CreateSecretCode(db *gorm.DB, code *SecretCode) error {
	return db.Create(code).Error
}

func GetSecretCodeByID(db *gorm.DB, id uint) (*SecretCode, error) {
	var code SecretCode
	if err := db.First(&code, id).Error; err != nil {
		return nil, err
	}
	return &code, nil
}

func CreateSecretCodeRedemption(db *gorm.DB, redemption *SecretCodeRedemption) error {
	return db.Create(redemption).Error
}

func GetSecretCodeRedemptionByID(db *gorm.DB, id uint) (*SecretCodeRedemption, error) {
	var redemption SecretCodeRedemption
	if err := db.First(&redemption, id).Error; err != nil {
		return nil, err
	}
	return &redemption, nil
}
