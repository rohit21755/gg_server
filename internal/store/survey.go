package store

import (
	"time"

	"gorm.io/gorm"
)

type Survey struct {
	ID         uint       `gorm:"primaryKey"`
	Title      string     `gorm:"size:200;not null"`
	Description *string   `gorm:"type:text"`
	SurveyType *string    `gorm:"size:50;check:survey_type IN ('feedback', 'quiz', 'research', 'poll')"`
	Questions  string     `gorm:"type:jsonb;not null"`
	XPReward   int        `gorm:"default:0"`
	IsActive   bool       `gorm:"default:true"`
	StartDate  *time.Time `gorm:"type:timestamp"`
	EndDate    *time.Time `gorm:"type:timestamp"`
	CreatedBy  *int       `gorm:"index"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`

	// Relations
	Creator *User `gorm:"foreignKey:CreatedBy"`
}

func (Survey) TableName() string {
	return "surveys"
}

type SurveyResponse struct {
	ID                 uint      `gorm:"primaryKey"`
	SurveyID           *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	UserID             *int      `gorm:"index;constraint:OnDelete:CASCADE"`
	Responses          string    `gorm:"type:jsonb;not null"`
	CompletionPercentage int     `gorm:"default:100"`
	XPAwarded          int       `gorm:"default:0"`
	SubmittedAt        time.Time `gorm:"autoCreateTime"`

	// Relations
	Survey *Survey `gorm:"foreignKey:SurveyID"`
	User   *User   `gorm:"foreignKey:UserID"`
}

func (SurveyResponse) TableName() string {
	return "survey_responses"
}

func CreateSurvey(db *gorm.DB, survey *Survey) error {
	return db.Create(survey).Error
}

func GetSurveyByID(db *gorm.DB, id uint) (*Survey, error) {
	var survey Survey
	if err := db.First(&survey, id).Error; err != nil {
		return nil, err
	}
	return &survey, nil
}

func CreateSurveyResponse(db *gorm.DB, response *SurveyResponse) error {
	return db.Create(response).Error
}

func GetSurveyResponseByID(db *gorm.DB, id uint) (*SurveyResponse, error) {
	var response SurveyResponse
	if err := db.First(&response, id).Error; err != nil {
		return nil, err
	}
	return &response, nil
}
