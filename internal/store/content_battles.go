package store

import (
	"time"

	"gorm.io/gorm"
)

type ContentBattle struct {
	ID               uint       `gorm:"primaryKey"`
	Title            string     `gorm:"size:200;not null"`
	Description      *string    `gorm:"type:text"`
	BattleType       *string    `gorm:"size:50;check:battle_type IN ('meme', 'video', 'reel', 'post')"`
	Theme            *string    `gorm:"size:200"`
	SubmissionDeadline time.Time `gorm:"not null"`
	VotingStart      time.Time  `gorm:"not null"`
	VotingEnd        time.Time  `gorm:"not null"`
	MaxParticipants  *int       `gorm:"type:integer"`
	Rewards          *string    `gorm:"type:jsonb;default:'{}'"`
	Status           string     `gorm:"size:20;default:'upcoming';check:status IN ('upcoming', 'submissions', 'voting', 'completed')"`
	CreatedBy        *int       `gorm:"index"`
	CreatedAt        time.Time  `gorm:"autoCreateTime"`

	// Relations
	Creator *User `gorm:"foreignKey:CreatedBy"`
}

func (ContentBattle) TableName() string {
	return "content_battles"
}

type BattleSubmission struct {
	ID           uint      `gorm:"primaryKey"`
	BattleID     *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	UserID       *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	Title        *string   `gorm:"size:200"`
	Description  *string   `gorm:"type:text"`
	MediaURL     string    `gorm:"type:text;not null"`
	ThumbnailURL *string   `gorm:"type:text"`
	VoteCount    int       `gorm:"default:0"`
	Rank         *int      `gorm:"type:integer"`
	IsWinner     bool      `gorm:"default:false"`
	SubmittedAt  time.Time `gorm:"autoCreateTime"`

	// Relations
	Battle *ContentBattle `gorm:"foreignKey:BattleID"`
	User   *User          `gorm:"foreignKey:UserID"`
}

func (BattleSubmission) TableName() string {
	return "battle_submissions"
}

type BattleVote struct {
	ID          uint      `gorm:"primaryKey"`
	BattleID    *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	SubmissionID *int     `gorm:"index;constraint:OnDelete:CASCADE"`
	VoterID     *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	VotedAt     time.Time `gorm:"autoCreateTime"`

	// Relations
	Battle     *ContentBattle  `gorm:"foreignKey:BattleID"`
	Submission *BattleSubmission `gorm:"foreignKey:SubmissionID"`
	Voter      *User           `gorm:"foreignKey:VoterID"`
}

func (BattleVote) TableName() string {
	return "battle_votes"
}

func CreateContentBattle(db *gorm.DB, battle *ContentBattle) error {
	return db.Create(battle).Error
}

func GetContentBattleByID(db *gorm.DB, id uint) (*ContentBattle, error) {
	var battle ContentBattle
	if err := db.First(&battle, id).Error; err != nil {
		return nil, err
	}
	return &battle, nil
}

func CreateBattleSubmission(db *gorm.DB, submission *BattleSubmission) error {
	return db.Create(submission).Error
}

func GetBattleSubmissionByID(db *gorm.DB, id uint) (*BattleSubmission, error) {
	var submission BattleSubmission
	if err := db.First(&submission, id).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

func CreateBattleVote(db *gorm.DB, vote *BattleVote) error {
	return db.Create(vote).Error
}

func GetBattleVoteByID(db *gorm.DB, id uint) (*BattleVote, error) {
	var vote BattleVote
	if err := db.First(&vote, id).Error; err != nil {
		return nil, err
	}
	return &vote, nil
}
