package store

import (
	"time"

	"gorm.io/gorm"
)

type State struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;unique;not null"`
	Code      string    `gorm:"size:10;unique;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (State) TableName() string {
	return "states"
}

type College struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:200;not null"`
	StateID   *int      `gorm:"index"`
	Code      *string   `gorm:"size:50;unique"`
	TotalCAs  int       `gorm:"default:0"`
	TotalXP   int       `gorm:"default:0"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relations
	State *State `gorm:"foreignKey:StateID"`
}

func (College) TableName() string {
	return "colleges"
}

func CreateState(db *gorm.DB, state *State) error {
	return db.Create(state).Error
}

func GetStateByID(db *gorm.DB, id uint) (*State, error) {
	var state State
	if err := db.First(&state, id).Error; err != nil {
		return nil, err
	}
	return &state, nil
}

func CreateCollege(db *gorm.DB, college *College) error {
	return db.Create(college).Error
}

func GetCollegeByID(db *gorm.DB, id uint) (*College, error) {
	var college College
	if err := db.First(&college, id).Error; err != nil {
		return nil, err
	}
	return &college, nil
}
