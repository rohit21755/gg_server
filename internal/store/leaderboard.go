package store

import (
	"time"

	"gorm.io/gorm"
)

type Leaderboard struct {
	ID            uint       `gorm:"primaryKey"`
	Name          string     `gorm:"size:100;not null"`
	LeaderboardType string   `gorm:"size:50;not null;check:leaderboard_type IN ('global', 'state', 'college', 'weekly', 'monthly', 'campaign', 'war')"`
	EntityType    *string    `gorm:"size:50;check:entity_type IN ('user', 'college', 'state')"`
	EntityID      *int       `gorm:"type:integer"`
	PeriodStart   *time.Time `gorm:"type:date"`
	PeriodEnd     *time.Time `gorm:"type:date"`
	Metrics       *string    `gorm:"type:jsonb;default:'{\"xp\": true}'"`
	IsActive      bool       `gorm:"default:true"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
}

func (Leaderboard) TableName() string {
	return "leaderboards"
}

type LeaderboardEntry struct {
	ID              uint       `gorm:"primaryKey"`
	LeaderboardID   *int       `gorm:"index;constraint:OnDelete:CASCADE"`
	UserID          *int       `gorm:"index;constraint:OnDelete:CASCADE"`
	CollegeID       *int       `gorm:"index"`
	StateID         *int       `gorm:"index"`
	XP              int        `gorm:"default:0"`
	SubmissionsCount int       `gorm:"default:0"`
	ReferralsCount  int        `gorm:"default:0"`
	WinRate         float64    `gorm:"type:decimal(5,2);default:0"`
	Rank            *int       `gorm:"type:integer"`
	PreviousRank    *int       `gorm:"type:integer"`
	Trend           *string    `gorm:"size:10;check:trend IN ('up', 'down', 'stable', 'new')"`
	SnapshotDate    time.Time  `gorm:"type:date;not null"`
	CreatedAt       time.Time  `gorm:"autoCreateTime"`

	// Relations
	Leaderboard *Leaderboard `gorm:"foreignKey:LeaderboardID"`
	User        *User        `gorm:"foreignKey:UserID"`
	College     *College     `gorm:"foreignKey:CollegeID"`
	State       *State       `gorm:"foreignKey:StateID"`
}

func (LeaderboardEntry) TableName() string {
	return "leaderboard_entries"
}

func CreateLeaderboard(db *gorm.DB, leaderboard *Leaderboard) error {
	return db.Create(leaderboard).Error
}

func GetLeaderboardByID(db *gorm.DB, id uint) (*Leaderboard, error) {
	var leaderboard Leaderboard
	if err := db.First(&leaderboard, id).Error; err != nil {
		return nil, err
	}
	return &leaderboard, nil
}

func CreateLeaderboardEntry(db *gorm.DB, entry *LeaderboardEntry) error {
	return db.Create(entry).Error
}

func GetLeaderboardEntryByID(db *gorm.DB, id uint) (*LeaderboardEntry, error) {
	var entry LeaderboardEntry
	if err := db.First(&entry, id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}
