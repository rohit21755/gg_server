package store

import (
	"time"

	"gorm.io/gorm"
)

type Quest struct {
	ID            uint       `gorm:"primaryKey"`
	Title         string     `gorm:"size:200;not null"`
	Description   *string    `gorm:"type:text"`
	QuestType     *string    `gorm:"size:50;check:quest_type IN ('multi_step', 'achievement', 'special')"`
	Steps         string     `gorm:"type:jsonb;not null"`
	Rewards       *string    `gorm:"type:jsonb;default:'{}'"`
	TimeLimitDays *int       `gorm:"type:integer"`
	IsActive      bool       `gorm:"default:true"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
}

func (Quest) TableName() string {
	return "quests"
}

type UserQuest struct {
	ID           uint       `gorm:"primaryKey"`
	UserID       *int       `gorm:"index;constraint:OnDelete:CASCADE"`
	QuestID      *int       `gorm:"index"`
	CurrentStep  int        `gorm:"default:1"`
	ProgressData *string    `gorm:"type:jsonb;default:'{}'"`
	Status       string     `gorm:"size:20;default:'in_progress';check:status IN ('in_progress', 'completed', 'failed', 'abandoned')"`
	StartedAt    time.Time  `gorm:"autoCreateTime"`
	CompletedAt  *time.Time `gorm:"type:timestamp"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`

	// Relations
	User  *User  `gorm:"foreignKey:UserID"`
	Quest *Quest `gorm:"foreignKey:QuestID"`
}

func (UserQuest) TableName() string {
	return "user_quests"
}

func CreateQuest(db *gorm.DB, quest *Quest) error {
	return db.Create(quest).Error
}

func GetQuestByID(db *gorm.DB, id uint) (*Quest, error) {
	var quest Quest
	if err := db.First(&quest, id).Error; err != nil {
		return nil, err
	}
	return &quest, nil
}

func CreateUserQuest(db *gorm.DB, userQuest *UserQuest) error {
	return db.Create(userQuest).Error
}

func GetUserQuestByID(db *gorm.DB, id uint) (*UserQuest, error) {
	var userQuest UserQuest
	if err := db.First(&userQuest, id).Error; err != nil {
		return nil, err
	}
	return &userQuest, nil
}
