package store

import (
	"time"

	"gorm.io/gorm"
)

type BadgeBingo struct {
	ID           uint       `gorm:"primaryKey"`
	Title        string     `gorm:"size:200;not null"`
	Description  *string    `gorm:"type:text"`
	BingoCard    string     `gorm:"type:jsonb;not null"`
	BundleRewards *string   `gorm:"type:jsonb;default:'{}'"`
	StartDate    time.Time  `gorm:"not null"`
	EndDate      time.Time  `gorm:"not null"`
	IsActive     bool       `gorm:"default:true"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
}

func (BadgeBingo) TableName() string {
	return "badge_bingo"
}

type UserBingoProgress struct {
	ID               uint      `gorm:"primaryKey"`
	UserID           *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	BingoID          *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	CompletedCells   *string   `gorm:"type:jsonb;default:'[]'"`
	CompletedRows    int       `gorm:"default:0"`
	CompletedColumns int       `gorm:"default:0"`
	CompletedDiagonals int     `gorm:"default:0"`
	RewardsClaimed   *string   `gorm:"type:jsonb;default:'[]'"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`

	// Relations
	User  *User      `gorm:"foreignKey:UserID"`
	Bingo *BadgeBingo `gorm:"foreignKey:BingoID"`
}

func (UserBingoProgress) TableName() string {
	return "user_bingo_progress"
}

func CreateBadgeBingo(db *gorm.DB, bingo *BadgeBingo) error {
	return db.Create(bingo).Error
}

func GetBadgeBingoByID(db *gorm.DB, id uint) (*BadgeBingo, error) {
	var bingo BadgeBingo
	if err := db.First(&bingo, id).Error; err != nil {
		return nil, err
	}
	return &bingo, nil
}

func CreateUserBingoProgress(db *gorm.DB, progress *UserBingoProgress) error {
	return db.Create(progress).Error
}

func GetUserBingoProgressByID(db *gorm.DB, id uint) (*UserBingoProgress, error) {
	var progress UserBingoProgress
	if err := db.First(&progress, id).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}
