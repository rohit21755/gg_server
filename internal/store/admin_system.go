package store

import (
	"time"

	"gorm.io/gorm"
)

type AdminAction struct {
	ID           uint      `gorm:"primaryKey"`
	AdminID      *int      `gorm:"index"`
	ActionType   string    `gorm:"size:100;not null"`
	ResourceType string    `gorm:"size:50;not null"`
	ResourceID   *int      `gorm:"type:integer"`
	Changes      *string   `gorm:"type:jsonb"`
	IPAddress    *string   `gorm:"type:inet"`
	UserAgent    *string   `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	// Relations
	Admin *User `gorm:"foreignKey:AdminID;constraint:OnDelete:SET NULL"`
}

func (AdminAction) TableName() string {
	return "admin_actions"
}

type SystemConfig struct {
	ID          uint       `gorm:"primaryKey"`
	ConfigKey   string     `gorm:"size:100;unique;not null"`
	ConfigValue string     `gorm:"type:jsonb;not null"`
	Description *string    `gorm:"type:text"`
	IsPublic    bool       `gorm:"default:false"`
	UpdatedBy   *int       `gorm:"index"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`

	// Relations
	Updater *User `gorm:"foreignKey:UpdatedBy"`
}

func (SystemConfig) TableName() string {
	return "system_config"
}

type ScheduledJob struct {
	ID           uint       `gorm:"primaryKey"`
	JobType      string     `gorm:"size:100;not null"`
	JobData      *string    `gorm:"type:jsonb"`
	ScheduledFor time.Time  `gorm:"not null"`
	Status       string     `gorm:"size:20;default:'pending';check:status IN ('pending', 'running', 'completed', 'failed', 'cancelled')"`
	Result       *string    `gorm:"type:jsonb"`
	Attempts     int        `gorm:"default:0"`
	MaxAttempts  int        `gorm:"default:3"`
	ErrorMessage *string    `gorm:"type:text"`
	StartedAt    *time.Time `gorm:"type:timestamp"`
	CompletedAt  *time.Time `gorm:"type:timestamp"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
}

func (ScheduledJob) TableName() string {
	return "scheduled_jobs"
}

func CreateAdminAction(db *gorm.DB, action *AdminAction) error {
	return db.Create(action).Error
}

func GetAdminActionByID(db *gorm.DB, id uint) (*AdminAction, error) {
	var action AdminAction
	if err := db.First(&action, id).Error; err != nil {
		return nil, err
	}
	return &action, nil
}

func CreateSystemConfig(db *gorm.DB, config *SystemConfig) error {
	return db.Create(config).Error
}

func GetSystemConfigByID(db *gorm.DB, id uint) (*SystemConfig, error) {
	var config SystemConfig
	if err := db.First(&config, id).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func CreateScheduledJob(db *gorm.DB, job *ScheduledJob) error {
	return db.Create(job).Error
}

func GetScheduledJobByID(db *gorm.DB, id uint) (*ScheduledJob, error) {
	var job ScheduledJob
	if err := db.First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}
