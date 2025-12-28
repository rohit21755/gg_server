package store

import (
	"time"

	"gorm.io/gorm"
)

type Certificate struct {
	ID                uint       `gorm:"primaryKey"`
	UserID            *int       `gorm:"index;constraint:OnDelete:CASCADE"`
	CertificateType   *string    `gorm:"size:50;check:certificate_type IN ('achievement', 'completion', 'winner', 'participation')"`
	Title             string     `gorm:"size:200;not null"`
	Description       *string    `gorm:"type:text"`
	IssuingAuthority  string     `gorm:"size:200;default:'Grove Growth'"`
	IssueDate         time.Time  `gorm:"type:date;not null"`
	ExpiryDate        *time.Time `gorm:"type:date"`
	CertificateURL    string     `gorm:"type:text;not null"`
	TemplateID        *int       `gorm:"type:integer"`
	Metadata          *string    `gorm:"type:jsonb;default:'{}'"`
	CreatedAt         time.Time  `gorm:"autoCreateTime"`

	// Relations
	User *User `gorm:"foreignKey:UserID"`
}

func (Certificate) TableName() string {
	return "certificates"
}

func CreateCertificate(db *gorm.DB, certificate *Certificate) error {
	return db.Create(certificate).Error
}

func GetCertificateByID(db *gorm.DB, id uint) (*Certificate, error) {
	var certificate Certificate
	if err := db.First(&certificate, id).Error; err != nil {
		return nil, err
	}
	return &certificate, nil
}
