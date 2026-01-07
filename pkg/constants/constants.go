package constants

// User roles
const (
	RoleCA      = "ca"
	RoleAdmin   = "admin"
	RoleModerator = "moderator"
)

// Task statuses
const (
	TaskStatusDraft     = "draft"
	TaskStatusActive    = "active"
	TaskStatusCompleted = "completed"
	TaskStatusCancelled = "cancelled"
)

// Submission statuses
const (
	SubmissionStatusPending  = "pending"
	SubmissionStatusApproved = "approved"
	SubmissionStatusRejected = "rejected"
	SubmissionStatusAppealed = "appealed"
)

// Campaign statuses
const (
	CampaignStatusDraft     = "draft"
	CampaignStatusActive   = "active"
	CampaignStatusCompleted = "completed"
	CampaignStatusCancelled = "cancelled"
)

// Reward types
const (
	RewardTypeDigital = "digital"
	RewardTypePhysical = "physical"
	RewardTypeCash     = "cash"
	RewardTypeCoins    = "coins"
)

// Wallet transaction types
const (
	WalletTxCredit  = "credit"
	WalletTxDebit   = "debit"
	WalletTxTransfer = "transfer"
)

// Notification types
const (
	NotificationTypeTask        = "task"
	NotificationTypeSubmission  = "submission"
	NotificationTypeAchievement = "achievement"
	NotificationTypeReward     = "reward"
	NotificationTypeCampaign   = "campaign"
	NotificationTypeSystem     = "system"
)

// Post types
const (
	PostTypeText       = "text"
	PostTypeImage      = "image"
	PostTypeVideo      = "video"
	PostTypeAchievement = "achievement"
)

// Achievement categories
const (
	AchievementCategorySubmission = "submission"
	AchievementCategoryStreak     = "streak"
	AchievementCategoryXP         = "xp"
	AchievementCategoryReferral   = "referral"
	AchievementCategorySocial     = "social"
)

// Quest types
const (
	QuestTypeDaily   = "daily"
	QuestTypeWeekly  = "weekly"
	QuestTypeMonthly = "monthly"
)

// File upload limits
const (
	MaxFileSize      = 10 * 1024 * 1024 // 10MB
	MaxImageSize     = 5 * 1024 * 1024  // 5MB
	MaxVideoSize     = 50 * 1024 * 1024 // 50MB
	AllowedImageExts = "jpg,jpeg,png,gif,webp"
	AllowedVideoExts = "mp4,webm,mov"
	AllowedDocExts   = "pdf,doc,docx"
)

