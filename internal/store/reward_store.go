package store

import (
	"time"

	"gorm.io/gorm"
)

type RewardStore struct {
	ID              uint       `gorm:"primaryKey"`
	Name            string     `gorm:"size:200;not null"`
	Description     *string    `gorm:"type:text"`
	RewardType      string     `gorm:"size:50;not null;check:reward_type IN ('physical', 'digital', 'badge', 'xp_boost', 'profile_skin', 'certificate', 'gift_card')"`
	Category        *string    `gorm:"size:50"`
	ImageURL        *string    `gorm:"type:text"`
	XPCost          int        `gorm:"default:0"`
	CoinCost        int        `gorm:"default:0"`
	CashCost        *float64   `gorm:"type:decimal(10,2)"`
	QuantityAvailable *int     `gorm:"type:integer"`
	QuantitySold    int        `gorm:"default:0"`
	IsFeatured      bool       `gorm:"default:false"`
	IsActive        bool       `gorm:"default:true"`
	ValidityDays    *int       `gorm:"type:integer"`
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime"`
}

func (RewardStore) TableName() string {
	return "rewards_store"
}

type UserReward struct {
	ID             uint       `gorm:"primaryKey"`
	UserID         *int       `gorm:"index;constraint:OnDelete:CASCADE"`
	RewardID       *int       `gorm:"index"`
	RedemptionCode *string    `gorm:"size:100"`
	Status         string     `gorm:"size:20;default:'pending';check:status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled', 'expired')"`
	XPPaid         int        `gorm:"default:0"`
	CoinsPaid      int        `gorm:"default:0"`
	CashPaid       float64    `gorm:"type:decimal(10,2);default:0"`
	ShippingAddress *string   `gorm:"type:jsonb"`
	TrackingNumber *string    `gorm:"size:100"`
	ClaimedAt      time.Time  `gorm:"autoCreateTime"`
	DeliveredAt    *time.Time `gorm:"type:timestamp"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`

	// Relations
	User   *User        `gorm:"foreignKey:UserID"`
	Reward *RewardStore `gorm:"foreignKey:RewardID"`
}

func (UserReward) TableName() string {
	return "user_rewards"
}

func CreateRewardStore(db *gorm.DB, reward *RewardStore) error {
	return db.Create(reward).Error
}

func GetRewardStoreByID(db *gorm.DB, id uint) (*RewardStore, error) {
	var reward RewardStore
	if err := db.First(&reward, id).Error; err != nil {
		return nil, err
	}
	return &reward, nil
}

func CreateUserReward(db *gorm.DB, userReward *UserReward) error {
	return db.Create(userReward).Error
}

func GetUserRewardByID(db *gorm.DB, id uint) (*UserReward, error) {
	var userReward UserReward
	if err := db.First(&userReward, id).Error; err != nil {
		return nil, err
	}
	return &userReward, nil
}
