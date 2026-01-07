package store

import (
	"time"

	"gorm.io/gorm"
)

// SocialPost represents a post in the social feed
type SocialPost struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	MediaURLs   *string   `gorm:"type:text" json:"media_urls,omitempty"` // JSON array of URLs
	PostType    string    `gorm:"size:50;default:'text'" json:"post_type"` // text, image, video, achievement, etc.
	IsPublic    bool      `gorm:"default:true" json:"is_public"`
	LikesCount  int       `gorm:"default:0" json:"likes_count"`
	CommentsCount int     `gorm:"default:0" json:"comments_count"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (SocialPost) TableName() string { return "social_posts" }

// PostLike represents a like on a post
type PostLike struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    uint      `gorm:"not null;index" json:"post_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (PostLike) TableName() string { return "post_likes" }

// PostComment represents a comment on a post
type PostComment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    uint      `gorm:"not null;index" json:"post_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (PostComment) TableName() string { return "post_comments" }

// GORM helper functions
func CreateSocialPost(db *gorm.DB, post *SocialPost) error {
	return db.Create(post).Error
}

func GetSocialPost(db *gorm.DB, id uint) (*SocialPost, error) {
	var post SocialPost
	if err := db.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func GetSocialFeed(db *gorm.DB, userID uint, limit int, offset int) ([]SocialPost, error) {
	var posts []SocialPost
	query := db.Where("is_public = ? OR user_id = ?", true, userID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset)
	if err := query.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func LikePost(db *gorm.DB, postID, userID uint) error {
	// Check if already liked
	var existing PostLike
	err := db.Where("post_id = ? AND user_id = ?", postID, userID).First(&existing).Error
	if err == nil {
		return nil // Already liked
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Create like
	like := PostLike{PostID: postID, UserID: userID}
	if err := db.Create(&like).Error; err != nil {
		return err
	}

	// Update post likes count
	return db.Model(&SocialPost{}).Where("id = ?", postID).
		UpdateColumn("likes_count", gorm.Expr("likes_count + 1")).Error
}

func UnlikePost(db *gorm.DB, postID, userID uint) error {
	result := db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&PostLike{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		// Update post likes count
		return db.Model(&SocialPost{}).Where("id = ?", postID).
			UpdateColumn("likes_count", gorm.Expr("GREATEST(likes_count - 1, 0)")).Error
	}
	return nil
}

func CreatePostComment(db *gorm.DB, comment *PostComment) error {
	if err := db.Create(comment).Error; err != nil {
		return err
	}
	// Update post comments count
	return db.Model(&SocialPost{}).Where("id = ?", comment.PostID).
		UpdateColumn("comments_count", gorm.Expr("comments_count + 1")).Error
}

func GetPostComments(db *gorm.DB, postID uint, limit int) ([]PostComment, error) {
	var comments []PostComment
	query := db.Where("post_id = ?", postID).
		Where("deleted_at IS NULL").
		Order("created_at ASC").
		Preload("User")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func DeletePost(db *gorm.DB, postID, userID uint) error {
	now := time.Now()
	return db.Model(&SocialPost{}).
		Where("id = ? AND user_id = ?", postID, userID).
		Update("deleted_at", now).Error
}

