package store

import (
	"time"

	"gorm.io/gorm"
)

// UserWallet represents the user's wallet/currency balance
type UserWallet struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	Coins     int       `gorm:"default:0" json:"coins"`
	Cash      float64   `gorm:"type:decimal(10,2);default:0" json:"cash"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserWallet) TableName() string { return "user_wallets" }

// WalletTransaction represents wallet transaction history
type WalletTransaction struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Type        string    `gorm:"size:50;not null" json:"type"` // credit, debit, transfer
	Amount      float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency    string    `gorm:"size:10;not null;default:'coins'" json:"currency"` // coins, cash
	Description string    `gorm:"type:text" json:"description"`
	ReferenceID *uint     `json:"reference_id,omitempty"` // ID of related entity (reward, etc.)
	ReferenceType *string `gorm:"size:50" json:"reference_type,omitempty"` // reward, transfer, etc.
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (WalletTransaction) TableName() string { return "wallet_transactions" }

// GORM helper functions
func GetUserWallet(db *gorm.DB, userID uint) (*UserWallet, error) {
	var wallet UserWallet
	if err := db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create wallet if it doesn't exist
			wallet = UserWallet{UserID: userID, Coins: 0, Cash: 0}
			if err := db.Create(&wallet).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &wallet, nil
}

func UpdateWallet(db *gorm.DB, wallet *UserWallet) error {
	return db.Save(wallet).Error
}

func CreateWalletTransaction(db *gorm.DB, tx *WalletTransaction) error {
	return db.Create(tx).Error
}

func GetWalletTransactions(db *gorm.DB, userID uint, limit int) ([]WalletTransaction, error) {
	var transactions []WalletTransaction
	query := db.Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

