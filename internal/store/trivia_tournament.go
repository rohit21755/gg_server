package store

import (
	"time"

	"gorm.io/gorm"
)

type TriviaTournament struct {
	ID              uint       `gorm:"primaryKey"`
	Title           string     `gorm:"size:200;not null"`
	Description     *string    `gorm:"type:text"`
	TournamentType  string     `gorm:"size:50;default:'weekly';check:tournament_type IN ('weekly', 'monthly', 'special')"`
	Questions       string     `gorm:"type:jsonb;not null"`
	StartDate       time.Time  `gorm:"not null"`
	EndDate         time.Time  `gorm:"not null"`
	DurationMinutes int        `gorm:"default:10"`
	MaxParticipants *int       `gorm:"type:integer"`
	EntryFeeXP      int        `gorm:"default:0"`
	Rewards         *string    `gorm:"type:jsonb;default:'{}'"`
	Status          string     `gorm:"size:20;default:'upcoming';check:status IN ('upcoming', 'active', 'completed')"`
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
}

func (TriviaTournament) TableName() string {
	return "trivia_tournaments"
}

type TriviaParticipant struct {
	ID               uint      `gorm:"primaryKey"`
	TriviaID         *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	UserID           *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	Score            int       `gorm:"default:0"`
	CorrectAnswers   int       `gorm:"default:0"`
	TimeTakenSeconds *int      `gorm:"type:integer"`
	Rank             *int      `gorm:"type:integer"`
	RewardClaimed    bool      `gorm:"default:false"`
	ParticipatedAt   time.Time `gorm:"autoCreateTime"`

	// Relations
	Trivia *TriviaTournament `gorm:"foreignKey:TriviaID"`
	User   *User             `gorm:"foreignKey:UserID"`
}

func (TriviaParticipant) TableName() string {
	return "trivia_participants"
}

func CreateTriviaTournament(db *gorm.DB, tournament *TriviaTournament) error {
	return db.Create(tournament).Error
}

func GetTriviaTournamentByID(db *gorm.DB, id uint) (*TriviaTournament, error) {
	var tournament TriviaTournament
	if err := db.First(&tournament, id).Error; err != nil {
		return nil, err
	}
	return &tournament, nil
}

func CreateTriviaParticipant(db *gorm.DB, participant *TriviaParticipant) error {
	return db.Create(participant).Error
}

func GetTriviaParticipantByID(db *gorm.DB, id uint) (*TriviaParticipant, error) {
	var participant TriviaParticipant
	if err := db.First(&participant, id).Error; err != nil {
		return nil, err
	}
	return &participant, nil
}
