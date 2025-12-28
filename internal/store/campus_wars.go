package store

import (
	"time"

	"gorm.io/gorm"
)

type CampusWar struct {
	ID        uint       `gorm:"primaryKey"`
	Name      string     `gorm:"size:200;not null"`
	Description *string  `gorm:"type:text"`
	WarType   *string    `gorm:"size:50;check:war_type IN ('campus_vs_campus', 'state_vs_state', 'region_vs_region')"`
	StartDate time.Time  `gorm:"not null"`
	EndDate   time.Time  `gorm:"not null"`
	Status    string     `gorm:"size:20;default:'upcoming';check:status IN ('upcoming', 'active', 'completed')"`
	Metrics   *string    `gorm:"type:jsonb;default:'{\"xp\": true, \"submissions\": true, \"referrals\": true}'"`
	Rewards   *string    `gorm:"type:jsonb;default:'{}'"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
}

func (CampusWar) TableName() string {
	return "campus_wars"
}

type WarParticipant struct {
	ID              uint      `gorm:"primaryKey"`
	WarID           *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	EntityType      string    `gorm:"size:50;not null;check:entity_type IN ('college', 'state')"`
	EntityID        int       `gorm:"not null"`
	TotalXP         int       `gorm:"default:0"`
	TotalSubmissions int      `gorm:"default:0"`
	TotalReferrals  int       `gorm:"default:0"`
	Rank            *int      `gorm:"type:integer"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`

	// Relations
	War *CampusWar `gorm:"foreignKey:WarID"`
}

func (WarParticipant) TableName() string {
	return "war_participants"
}

type WarLeaderboardSnapshot struct {
	ID              uint      `gorm:"primaryKey"`
	WarID           *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	SnapshotDate    time.Time `gorm:"type:date;not null"`
	LeaderboardData string    `gorm:"type:jsonb;not null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`

	// Relations
	War *CampusWar `gorm:"foreignKey:WarID"`
}

func (WarLeaderboardSnapshot) TableName() string {
	return "war_leaderboard_snapshots"
}

func CreateCampusWar(db *gorm.DB, war *CampusWar) error {
	return db.Create(war).Error
}

func GetCampusWarByID(db *gorm.DB, id uint) (*CampusWar, error) {
	var war CampusWar
	if err := db.First(&war, id).Error; err != nil {
		return nil, err
	}
	return &war, nil
}

func CreateWarParticipant(db *gorm.DB, participant *WarParticipant) error {
	return db.Create(participant).Error
}

func GetWarParticipantByID(db *gorm.DB, id uint) (*WarParticipant, error) {
	var participant WarParticipant
	if err := db.First(&participant, id).Error; err != nil {
		return nil, err
	}
	return &participant, nil
}

func CreateWarLeaderboardSnapshot(db *gorm.DB, snapshot *WarLeaderboardSnapshot) error {
	return db.Create(snapshot).Error
}

func GetWarLeaderboardSnapshotByID(db *gorm.DB, id uint) (*WarLeaderboardSnapshot, error) {
	var snapshot WarLeaderboardSnapshot
	if err := db.First(&snapshot, id).Error; err != nil {
		return nil, err
	}
	return &snapshot, nil
}
