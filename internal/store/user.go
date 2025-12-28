package store

import (
	"time"

	"gorm.io/gorm"
)

// User represents the `users` table.
type User struct {
	ID                  uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID                string     `gorm:"type:varchar(36);unique;not null;default:gen_random_uuid()::text" json:"uuid"`
	Email               string     `gorm:"size:255;unique;not null" json:"email"`
	PasswordHash        string     `gorm:"size:255;not null" json:"-"`
	FirstName           string     `gorm:"size:100;not null" json:"first_name"`
	LastName            string     `gorm:"size:100;not null" json:"last_name"`
	Phone               *string    `gorm:"size:20" json:"phone,omitempty"`
	Role                string     `gorm:"size:20;not null;default:'ca'" json:"role"`
	CollegeID           *int       `json:"college_id,omitempty"`
	StateID             *int       `json:"state_id,omitempty"`
	ReferralCode        string     `gorm:"size:20;unique;not null" json:"referral_code"`
	ReferredBy          *int       `json:"referred_by,omitempty"`
	XP                  int        `gorm:"default:0" json:"xp"`
	LevelID             *int       `gorm:"default:1" json:"level_id,omitempty"`
	StreakCount         int        `gorm:"default:0" json:"streak_count"`
	LastLoginDate       *time.Time `json:"last_login_date,omitempty"`
	TotalSubmissions    int        `gorm:"default:0" json:"total_submissions"`
	ApprovedSubmissions int        `gorm:"default:0" json:"approved_submissions"`
	WinRate             float64    `gorm:"type:decimal(5,2);default:0" json:"win_rate"`
	ProfileSkinID       *int       `json:"profile_skin_id,omitempty"`
	AvatarURL           *string    `gorm:"type:text" json:"avatar_url,omitempty"`
	ResumeURL           *string    `gorm:"type:text" json:"resume_url,omitempty"`
	IsActive            bool       `gorm:"default:true" json:"is_active"`
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string { return "users" }

// UserSession represents the `user_sessions` table.
type UserSession struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       *int      `json:"user_id,omitempty"`
	SessionToken string    `gorm:"size:255;unique;not null" json:"session_token"`
	DeviceID     *string   `gorm:"size:255" json:"device_id,omitempty"`
	Platform     *string   `gorm:"size:50" json:"platform,omitempty"`
	LastActive   time.Time `gorm:"autoUpdateTime" json:"last_active"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (UserSession) TableName() string { return "user_sessions" }

// --- GORM helper functions ---

func CreateUser(db *gorm.DB, u *User) error {
	return db.Create(u).Error
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var u User
	if err := db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByUUID(db *gorm.DB, uuid string) (*User, error) {
	var u User
	if err := db.Where("uuid = ?", uuid).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var u User
	if err := db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUser(db *gorm.DB, u *User) error {
	return db.Save(u).Error
}

func DeleteUser(db *gorm.DB, id uint) error {
	return db.Delete(&User{}, id).Error
}

func CreateSession(db *gorm.DB, s *UserSession) error {
	return db.Create(s).Error
}

func GetSessionByToken(db *gorm.DB, token string) (*UserSession, error) {
	var s UserSession
	if err := db.Where("session_token = ?", token).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func DeleteSessionByToken(db *gorm.DB, token string) error {
	return db.Where("session_token = ?", token).Delete(&UserSession{}).Error
}

func DeleteExpiredSessions(db *gorm.DB) error {
	return db.Where("expires_at < ?", time.Now()).Delete(&UserSession{}).Error
}

func GetUserByReferralCode(db *gorm.DB, referralCode string) (*User, error) {
	var u User
	if err := db.Where("referral_code = ?", referralCode).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
func GetUserStreak(db *gorm.DB, userID uint, streakType string) (*UserStreak, error) {
	var streak UserStreak
	userIDInt := int(userID)
	if err := db.Where("user_id = ? AND streak_type = ?", userIDInt, streakType).First(&streak).Error; err != nil {
		return nil, err
	}
	return &streak, nil
}

func UpdateUserStreak(db *gorm.DB, streak *UserStreak) error {
	if streak.ID == 0 {
		return gorm.ErrRecordNotFound
	}
	result := db.Save(streak)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func GetUserWithRelations(db *gorm.DB, userID uint) (*User, error) {
	var user User
	if err := db.Preload("College").Preload("State").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserBadgeCount(db *gorm.DB, userID uint) (int, error) {
	var count int64
	if err := db.Model(&UserBadge{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func GetUserCertificates(db *gorm.DB, userID uint) ([]Certificate, error) {
	var certificates []Certificate
	if err := db.Where("user_id = ?", userID).Find(&certificates).Error; err != nil {
		return nil, err
	}
	return certificates, nil
}
