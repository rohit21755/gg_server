package store

import (
	"time"

	"gorm.io/gorm"
)

type SpinWheel struct {
	ID              uint       `gorm:"primaryKey"`
	Name            string     `gorm:"size:100;not null"`
	Description     *string    `gorm:"type:text"`
	WheelType       string     `gorm:"size:50;default:'weekly';check:wheel_type IN ('weekly', 'daily', 'special')"`
	IsActive        bool       `gorm:"default:true"`
	SpinsPerUser    int        `gorm:"default:1"`
	ResetFrequency  string     `gorm:"size:20;default:'weekly'"`
	StartDate       *time.Time `gorm:"type:timestamp"`
	EndDate         *time.Time `gorm:"type:timestamp"`
	MinActivityLevel int       `gorm:"default:0"`
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
}

func (SpinWheel) TableName() string {
	return "spin_wheels"
}

type SpinWheelItem struct {
	ID             uint      `gorm:"primaryKey"`
	SpinWheelID    *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	ItemType       string    `gorm:"size:50;not null;check:item_type IN ('xp', 'coins', 'badge', 'physical', 'discount')"`
	ItemValue      int       `gorm:"not null"`
	ItemLabel      string    `gorm:"size:100;not null"`
	Probability    float64   `gorm:"type:decimal(5,4);not null"`
	MaxQuantity    *int      `gorm:"type:integer"`
	CurrentQuantity *int     `gorm:"type:integer"`
	IsActive       bool      `gorm:"default:true"`
	SortOrder      int       `gorm:"default:0"`

	// Relations
	SpinWheel *SpinWheel `gorm:"foreignKey:SpinWheelID"`
}

func (SpinWheelItem) TableName() string {
	return "spin_wheel_items"
}

type UserSpin struct {
	ID             uint      `gorm:"primaryKey"`
	UserID         *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	SpinWheelID    *int      `gorm:"index"`
	SpinWheelItemID *int     `gorm:"index"`
	EarnedValue    int       `gorm:"not null"`
	SpunAt         time.Time `gorm:"autoCreateTime"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`

	// Relations
	User         *User         `gorm:"foreignKey:UserID"`
	SpinWheel    *SpinWheel    `gorm:"foreignKey:SpinWheelID"`
	SpinWheelItem *SpinWheelItem `gorm:"foreignKey:SpinWheelItemID"`
}

func (UserSpin) TableName() string {
	return "user_spins"
}

type MysteryBox struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:100;not null"`
	Description *string   `gorm:"type:text"`
	CostXP      int       `gorm:"not null"`
	CostCoins   int       `gorm:"default:0"`
	Contents    string    `gorm:"type:jsonb;not null"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (MysteryBox) TableName() string {
	return "mystery_boxes"
}

type MysteryBoxRedemption struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	MysteryBoxID *int     `gorm:"index"`
	RewardType  string    `gorm:"size:50;not null"`
	RewardValue int       `gorm:"not null"`
	RedeemedAt  time.Time `gorm:"autoCreateTime"`

	// Relations
	User       *User       `gorm:"foreignKey:UserID"`
	MysteryBox *MysteryBox `gorm:"foreignKey:MysteryBoxID"`
}

func (MysteryBoxRedemption) TableName() string {
	return "mystery_box_redemptions"
}

func CreateSpinWheel(db *gorm.DB, wheel *SpinWheel) error {
	return db.Create(wheel).Error
}

func GetSpinWheelByID(db *gorm.DB, id uint) (*SpinWheel, error) {
	var wheel SpinWheel
	if err := db.First(&wheel, id).Error; err != nil {
		return nil, err
	}
	return &wheel, nil
}

func CreateSpinWheelItem(db *gorm.DB, item *SpinWheelItem) error {
	return db.Create(item).Error
}

func GetSpinWheelItemByID(db *gorm.DB, id uint) (*SpinWheelItem, error) {
	var item SpinWheelItem
	if err := db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateUserSpin(db *gorm.DB, spin *UserSpin) error {
	return db.Create(spin).Error
}

func GetUserSpinByID(db *gorm.DB, id uint) (*UserSpin, error) {
	var spin UserSpin
	if err := db.First(&spin, id).Error; err != nil {
		return nil, err
	}
	return &spin, nil
}

func CreateMysteryBox(db *gorm.DB, box *MysteryBox) error {
	return db.Create(box).Error
}

func GetMysteryBoxByID(db *gorm.DB, id uint) (*MysteryBox, error) {
	var box MysteryBox
	if err := db.First(&box, id).Error; err != nil {
		return nil, err
	}
	return &box, nil
}

func CreateMysteryBoxRedemption(db *gorm.DB, redemption *MysteryBoxRedemption) error {
	return db.Create(redemption).Error
}

func GetMysteryBoxRedemptionByID(db *gorm.DB, id uint) (*MysteryBoxRedemption, error) {
	var redemption MysteryBoxRedemption
	if err := db.First(&redemption, id).Error; err != nil {
		return nil, err
	}
	return &redemption, nil
}

func GetUserSpinsToday(db *gorm.DB, userID uint, wheelID uint) (int, error) {
	var count int64
	userIDInt := int(userID)
	wheelIDInt := int(wheelID)
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	
	if err := db.Model(&UserSpin{}).
		Where("user_id = ? AND spin_wheel_id = ? AND spun_at >= ? AND spun_at < ?", 
			userIDInt, wheelIDInt, today, tomorrow).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
