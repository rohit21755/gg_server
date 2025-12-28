package store

import (
	"time"

	"gorm.io/gorm"
)

type Submission struct {
	ID               uint       `gorm:"primaryKey"`
	UUID             string     `gorm:"type:varchar(36);unique;not null"`
	TaskID           *int       `gorm:"index"`
	UserID           *int       `gorm:"index"`
	CampaignID       *int       `gorm:"index"`
	ProofType        string     `gorm:"size:50;not null"`
	ProofURL         string     `gorm:"type:text;not null"`
	ProofText        *string    `gorm:"type:text"`
	Status           string     `gorm:"size:20;default:'pending';check:status IN ('draft', 'pending', 'under_review', 'approved', 'rejected', 'needs_revision')"`
	SubmissionStage  string     `gorm:"size:20;default:'initial';check:submission_stage IN ('initial', 'resubmission')"`
	SubmittedAt      time.Time  `gorm:"autoCreateTime"`
	ReviewedAt       *time.Time `gorm:"type:timestamp"`
	ReviewedBy       *int       `gorm:"index"`
	ReviewComments   *string    `gorm:"type:text"`
	XPAwarded        int        `gorm:"default:0"`
	CoinsAwarded     int        `gorm:"default:0"`
	IsWinner         bool       `gorm:"default:false"`
	Score            *float64   `gorm:"type:decimal(5,2)"`
	RevisionDeadline *time.Time `gorm:"type:timestamp"`
	CreatedAt        time.Time  `gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime"`

	// Relations
	Task     *Task     `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE"`
	User     *User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Campaign *Campaign `gorm:"foreignKey:CampaignID"`
	Reviewer *User     `gorm:"foreignKey:ReviewedBy"`
}

func (Submission) TableName() string {
	return "submissions"
}

type SubmissionMedia struct {
	ID           uint      `gorm:"primaryKey"`
	SubmissionID *int      `gorm:"index"`
	MediaType    string    `gorm:"size:50;not null;check:media_type IN ('image', 'video', 'document', 'link')"`
	MediaURL     string    `gorm:"type:text;not null"`
	ThumbnailURL *string   `gorm:"type:text"`
	FileName     *string   `gorm:"size:255"`
	FileSize     *int      `gorm:"type:integer"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	// Relations
	Submission *Submission `gorm:"foreignKey:SubmissionID;constraint:OnDelete:CASCADE"`
}

func (SubmissionMedia) TableName() string {
	return "submission_media"
}

func CreateSubmission(db *gorm.DB, submission *Submission) error {
	return db.Create(submission).Error
}

func GetSubmissionByID(db *gorm.DB, id uint) (*Submission, error) {
	var submission Submission
	if err := db.First(&submission, id).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

func CreateSubmissionMedia(db *gorm.DB, media *SubmissionMedia) error {
	return db.Create(media).Error
}

func GetSubmissionMediaByID(db *gorm.DB, id uint) (*SubmissionMedia, error) {
	var media SubmissionMedia
	if err := db.First(&media, id).Error; err != nil {
		return nil, err
	}
	return &media, nil
}

func UpdateSubmission(db *gorm.DB, submission *Submission) error {
	return db.Save(submission).Error
}

func DeleteSubmission(db *gorm.DB, id uint) error {
	return db.Delete(&Submission{}, id).Error
}
func GetSubmissionsByUserAndTasks(db *gorm.DB, userID uint, taskIDs []uint) ([]Submission, error) {
	var submissions []Submission
	if err := db.Where("user_id = ? AND task_id IN (?)", userID, taskIDs).Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func GetTaskAssignmentsByUser(db *gorm.DB, userID uint) ([]TaskAssignment, error) {
	var assignments []TaskAssignment
	if err := db.Where("assignee_id = ?", userID).Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}
