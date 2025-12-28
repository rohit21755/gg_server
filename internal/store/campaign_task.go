package store

import (
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	ID                  uint      `gorm:"primaryKey"`
	UUID                string    `gorm:"type:varchar(36);unique;not null"`
	Title               string    `gorm:"size:200;not null"`
	Description         *string   `gorm:"type:text"`
	CampaignType        string    `gorm:"size:50;not null;check:campaign_type IN ('brand_specific', 'thematic', 'seasonal', 'gg_led', 'flash', 'weekly_vibe', 'limited_edition')"`
	Category            *string   `gorm:"size:50;check:category IN ('solo', 'group', 'online', 'offline')"`
	BannerImageURL      *string   `gorm:"type:text"`
	StartDate           time.Time `gorm:"not null"`
	EndDate             time.Time `gorm:"not null"`
	MaxParticipants     *int      `gorm:"type:integer"`
	CurrentParticipants int       `gorm:"default:0"`
	Status              string    `gorm:"size:20;default:'draft';check:status IN ('draft', 'active', 'paused', 'completed', 'cancelled')"`
	Priority            string    `gorm:"size:20;default:'medium';check:priority IN ('low', 'medium', 'high')"`
	CreatedBy           *int      `gorm:"index"`
	IsLimitedEdition    bool      `gorm:"default:false"`
	IsGGLed             bool      `gorm:"default:false"`
	Metadata            *string   `gorm:"type:jsonb;default:'{}'"`
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`

	// Relations
	Creator *User `gorm:"foreignKey:CreatedBy"`
}

func (Campaign) TableName() string {
	return "campaigns"
}

type Task struct {
	ID                     uint      `gorm:"primaryKey"`
	UUID                   string    `gorm:"type:varchar(36);unique;not null"`
	CampaignID             *int      `gorm:"index"`
	Title                  string    `gorm:"size:200;not null"`
	Description            string    `gorm:"type:text;not null"`
	TaskType               string    `gorm:"size:50;not null;check:task_type IN ('solo', 'group', 'online', 'offline')"`
	ProofType              string    `gorm:"size:50;not null;check:proof_type IN ('screenshot', 'url', 'pdf', 'video', 'text')"`
	XPReward               int       `gorm:"not null;default:0"`
	CoinReward             int       `gorm:"default:0"`
	DurationHours          *int      `gorm:"type:integer"`
	Priority               string    `gorm:"size:20;default:'medium';check:priority IN ('low', 'medium', 'high', 'flash')"`
	AssignmentType         *string   `gorm:"size:50;check:assignment_type IN ('role', 'college', 'state', 'individual')"`
	AssignmentTarget       *string   `gorm:"type:jsonb;default:'{}'"`
	MaxSubmissions         int       `gorm:"default:1"`
	IsActive               bool      `gorm:"default:true"`
	SubmissionInstructions *string   `gorm:"type:text"`
	CreatedBy              *int      `gorm:"index"`
	CreatedAt              time.Time `gorm:"autoCreateTime"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime"`

	// Relations
	Campaign *Campaign `gorm:"foreignKey:CampaignID;constraint:OnDelete:CASCADE"`
	Creator  *User     `gorm:"foreignKey:CreatedBy"`
}

func (Task) TableName() string {
	return "tasks"
}

type TaskAssignment struct {
	ID           uint      `gorm:"primaryKey"`
	TaskID       *int      `gorm:"index"`
	AssigneeType string    `gorm:"size:50;not null;check:assignee_type IN ('user', 'role', 'college', 'state')"`
	AssigneeID   *int      `gorm:"index"`
	AssigneeRole *string   `gorm:"size:20"`
	AssignedBy   *int      `gorm:"index"`
	AssignedAt   time.Time `gorm:"autoCreateTime"`
	Status       string    `gorm:"size:20;default:'assigned';check:status IN ('assigned', 'accepted', 'declined', 'completed')"`

	// Relations
	Task     *Task `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE"`
	Assigner *User `gorm:"foreignKey:AssignedBy"`
}

func (TaskAssignment) TableName() string {
	return "task_assignments"
}

func CreateCampaign(db *gorm.DB, campaign *Campaign) error {
	return db.Create(campaign).Error
}

func GetCampaignByID(db *gorm.DB, id uint) (*Campaign, error) {
	var campaign Campaign
	if err := db.First(&campaign, id).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

func CreateTask(db *gorm.DB, task *Task) error {
	return db.Create(task).Error
}

func GetTaskByID(db *gorm.DB, id uint) (*Task, error) {
	var task Task
	if err := db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func CreateTaskAssignment(db *gorm.DB, assignment *TaskAssignment) error {
	return db.Create(assignment).Error
}

func GetTaskAssignmentByID(db *gorm.DB, id uint) (*TaskAssignment, error) {
	var assignment TaskAssignment
	if err := db.First(&assignment, id).Error; err != nil {
		return nil, err
	}
	return &assignment, nil
}

func UpdateCampaign(db *gorm.DB, campaign *Campaign) error {
	return db.Save(campaign).Error
}

func GetSubmissionsByUserAndTask(db *gorm.DB, userID uint, taskID uint) ([]Submission, error) {
	var submissions []Submission
	if err := db.Where("user_id = ? AND task_id = ?", userID, taskID).Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}
