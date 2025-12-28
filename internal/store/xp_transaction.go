package store

import (
	"time"

	"gorm.io/gorm"
)

type XPTransaction struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	TransactionType string    `gorm:"size:50;not null;check:transaction_type IN ('task_completion', 'referral', 'streak', 'spin_wheel', 'mystery_box', 'quiz', 'battle_win', 'redemption', 'correction', 'bonus')"`
	Amount          int       `gorm:"not null"`
	BalanceAfter    int       `gorm:"not null"`
	SourceID        *int      `gorm:"type:integer"`
	SourceType      *string   `gorm:"size:50"`
	Description     *string   `gorm:"type:text"`
	Metadata        *string   `gorm:"type:jsonb;default:'{}'"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`

	// Relations
	User *User `gorm:"foreignKey:UserID"`
}

func (XPTransaction) TableName() string {
	return "xp_transactions"
}

func CreateXPTransaction(db *gorm.DB, transaction *XPTransaction) error {
	return db.Create(transaction).Error
}

func GetXPTransactionByID(db *gorm.DB, id uint) (*XPTransaction, error) {
	var transaction XPTransaction
	if err := db.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}
