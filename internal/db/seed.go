package db

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/rohit21755/gg_server.git/internal/store"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed populates the database with initial data
func Seed(database *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Seed in order of dependencies
	if err := seedStates(database); err != nil {
		return fmt.Errorf("failed to seed states: %w", err)
	}

	if err := seedColleges(database); err != nil {
		return fmt.Errorf("failed to seed colleges: %w", err)
	}

	if err := seedUsers(database); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	if err := seedBadges(database); err != nil {
		return fmt.Errorf("failed to seed badges: %w", err)
	}

	if err := seedCampaigns(database); err != nil {
		return fmt.Errorf("failed to seed campaigns: %w", err)
	}

	if err := seedTasks(database); err != nil {
		return fmt.Errorf("failed to seed tasks: %w", err)
	}

	if err := seedRewards(database); err != nil {
		return fmt.Errorf("failed to seed rewards: %w", err)
	}

	if err := seedSpinWheel(database); err != nil {
		return fmt.Errorf("failed to seed spin wheel: %w", err)
	}

	if err := seedMysteryBoxes(database); err != nil {
		return fmt.Errorf("failed to seed mystery boxes: %w", err)
	}

	if err := seedSecretCodes(database); err != nil {
		return fmt.Errorf("failed to seed secret codes: %w", err)
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

// seedStates seeds Indian states
func seedStates(db *gorm.DB) error {
	states := []store.State{
		{Name: "Maharashtra", Code: "MH"},
		{Name: "Karnataka", Code: "KA"},
		{Name: "Tamil Nadu", Code: "TN"},
		{Name: "Delhi", Code: "DL"},
		{Name: "Gujarat", Code: "GJ"},
		{Name: "Rajasthan", Code: "RJ"},
		{Name: "West Bengal", Code: "WB"},
		{Name: "Uttar Pradesh", Code: "UP"},
		{Name: "Telangana", Code: "TG"},
		{Name: "Kerala", Code: "KL"},
		{Name: "Punjab", Code: "PB"},
		{Name: "Haryana", Code: "HR"},
		{Name: "Madhya Pradesh", Code: "MP"},
		{Name: "Bihar", Code: "BR"},
		{Name: "Andhra Pradesh", Code: "AP"},
	}

	for _, state := range states {
		var existing store.State
		if err := db.Where("code = ?", state.Code).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := store.CreateState(db, &state); err != nil {
					return fmt.Errorf("failed to create state %s: %w", state.Name, err)
				}
				log.Printf("Created state: %s", state.Name)
			} else {
				return err
			}
		}
	}
	return nil
}

// seedColleges seeds colleges across different states
func seedColleges(db *gorm.DB) error {
	// Get state IDs
	var states []store.State
	if err := db.Find(&states).Error; err != nil {
		return err
	}

	stateMap := make(map[string]uint)
	for _, state := range states {
		stateMap[state.Code] = state.ID
	}

	colleges := []struct {
		name    string
		state   string
		code    string
		totalXP int
	}{
		{"Indian Institute of Technology Bombay", "MH", "IITB", 50000},
		{"Indian Institute of Technology Delhi", "DL", "IITD", 45000},
		{"Indian Institute of Technology Madras", "TN", "IITM", 40000},
		{"Indian Institute of Technology Bangalore", "KA", "IITB", 38000},
		{"National Institute of Technology Surathkal", "KA", "NITSK", 25000},
		{"Vellore Institute of Technology", "TN", "VIT", 35000},
		{"Birla Institute of Technology and Science", "RJ", "BITS", 30000},
		{"Manipal Institute of Technology", "KA", "MIT", 28000},
		{"SRM Institute of Science and Technology", "TN", "SRM", 22000},
		{"Amity University", "UP", "AU", 20000},
		{"Symbiosis International University", "MH", "SIU", 18000},
		{"Christ University", "KA", "CU", 15000},
		{"Jadavpur University", "WB", "JU", 12000},
		{"Delhi Technological University", "DL", "DTU", 10000},
		{"National Institute of Technology Calicut", "KL", "NITC", 8000},
	}

	for _, c := range colleges {
		var existing store.College
		if err := db.Where("code = ?", c.code).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				stateID := int(stateMap[c.state])
				college := &store.College{
					Name:     c.name,
					StateID:  &stateID,
					Code:     &c.code,
					TotalXP:  c.totalXP,
					IsActive: true,
				}
				if err := store.CreateCollege(db, college); err != nil {
					return fmt.Errorf("failed to create college %s: %w", c.name, err)
				}
				log.Printf("Created college: %s", c.name)
			} else {
				return err
			}
		}
	}
	return nil
}

// seedUsers seeds admin and regular user
func seedUsers(db *gorm.DB) error {
	// Get first college and state for users
	var college store.College
	if err := db.First(&college).Error; err != nil {
		return fmt.Errorf("no colleges found, seed colleges first: %w", err)
	}

	var state store.State
	if err := db.First(&state).Error; err != nil {
		return fmt.Errorf("no states found, seed states first: %w", err)
	}

	collegeID := int(college.ID)
	stateID := int(state.ID)

	// Get first level (Rookie)
	var level store.Level
	if err := db.Where("rank_order = ?", 1).First(&level).Error; err != nil {
		return fmt.Errorf("no levels found: %w", err)
	}
	levelID := int(level.ID)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	users := []struct {
		email     string
		firstName string
		lastName  string
		role      string
		xp        int
		levelID   *int
	}{
		{
			email:     "admin@campusambassador.com",
			firstName: "Admin",
			lastName:  "User",
			role:      "admin",
			xp:        100000,
			levelID:   nil, // Will set to highest level
		},
		{
			email:     "user@campusambassador.com",
			firstName: "John",
			lastName:  "Doe",
			role:      "ca",
			xp:        2500,
			levelID:   &levelID,
		},
	}

	for _, u := range users {
		var existing store.User
		if err := db.Where("email = ?", u.email).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Generate referral code
				referralCode := generateReferralCode()

				user := &store.User{
					UUID:         uuid.New().String(),
					Email:        u.email,
					PasswordHash: string(hashedPassword),
					FirstName:    u.firstName,
					LastName:     u.lastName,
					Role:         u.role,
					CollegeID:    &collegeID,
					StateID:      &stateID,
					ReferralCode: referralCode,
					XP:           u.xp,
					LevelID:      u.levelID,
					IsActive:     true,
				}

				// Set admin to highest level
				if u.role == "admin" {
					var highestLevel store.Level
					if err := db.Order("rank_order DESC").First(&highestLevel).Error; err == nil {
						highestLevelID := int(highestLevel.ID)
						user.LevelID = &highestLevelID
					}
				}

				if err := store.CreateUser(db, user); err != nil {
					return fmt.Errorf("failed to create user %s: %w", u.email, err)
				}
				log.Printf("Created user: %s (%s)", u.email, u.role)
			} else {
				return err
			}
		}
	}
	return nil
}

// seedBadges seeds various badges
func seedBadges(db *gorm.DB) error {
	badges := []store.Badge{
		{
			Name:          "First Steps",
			Description:   stringPtr("Complete your first task"),
			BadgeType:     "achievement",
			Category:      stringPtr("submission"),
			ImageURL:      "https://example.com/badges/first-steps.png",
			XPReward:      100,
			CriteriaType:  "submission_count",
			CriteriaValue: 1,
			IsSecret:      false,
		},
		{
			Name:          "Referral Master",
			Description:   stringPtr("Refer 10 friends"),
			BadgeType:     "achievement",
			Category:      stringPtr("referral"),
			ImageURL:      "https://example.com/badges/referral-master.png",
			XPReward:      500,
			CriteriaType:  "referral_count",
			CriteriaValue: 10,
			IsSecret:      false,
		},
		{
			Name:          "Campaign Champion",
			Description:   stringPtr("Complete 5 campaigns"),
			BadgeType:     "achievement",
			Category:      stringPtr("campaign"),
			ImageURL:      "https://example.com/badges/campaign-champion.png",
			XPReward:      1000,
			CriteriaType:  "campaign_count",
			CriteriaValue: 5,
			IsSecret:      false,
		},
		{
			Name:          "Streak Warrior",
			Description:   stringPtr("Maintain a 30-day streak"),
			BadgeType:     "achievement",
			Category:      stringPtr("streak"),
			ImageURL:      "https://example.com/badges/streak-warrior.png",
			XPReward:      750,
			CriteriaType:  "streak_days",
			CriteriaValue: 30,
			IsSecret:      false,
		},
		{
			Name:          "Hidden Gem",
			Description:   stringPtr("Secret achievement badge"),
			BadgeType:     "special",
			Category:      stringPtr("special"),
			ImageURL:      "https://example.com/badges/hidden-gem.png",
			XPReward:      2000,
			CriteriaType:  "special",
			CriteriaValue: 1,
			IsSecret:      true,
		},
		{
			Name:          "Early Bird",
			Description:   stringPtr("Join in the first month"),
			BadgeType:     "limited",
			Category:      stringPtr("special"),
			ImageURL:      "https://example.com/badges/early-bird.png",
			XPReward:      300,
			CriteriaType:  "join_date",
			CriteriaValue: 1,
			IsSecret:      false,
			IsLimitedEdition: true,
		},
	}

	for _, badge := range badges {
		var existing store.Badge
		if err := db.Where("name = ?", badge.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&badge).Error; err != nil {
					return fmt.Errorf("failed to create badge %s: %w", badge.Name, err)
				}
				log.Printf("Created badge: %s", badge.Name)
			} else {
				return err
			}
		}
	}
	return nil
}

// seedCampaigns seeds campaigns
func seedCampaigns(db *gorm.DB) error {
	// Get admin user
	var admin store.User
	if err := db.Where("role = ?", "admin").First(&admin).Error; err != nil {
		return fmt.Errorf("admin user not found: %w", err)
	}

	adminID := int(admin.ID)
	now := time.Now()

	campaigns := []store.Campaign{
		{
			UUID:                uuid.New().String(),
			Title:               "Summer Brand Campaign 2024",
			Description:         stringPtr("Engage with our summer product launch and earn exciting rewards"),
			CampaignType:        "brand_specific",
			Category:            stringPtr("online"),
			BannerImageURL:      stringPtr("https://example.com/campaigns/summer-2024.jpg"),
			StartDate:           now.AddDate(0, 0, -10),
			EndDate:             now.AddDate(0, 1, 0),
			MaxParticipants:     intPtr(1000),
			CurrentParticipants: 150,
			Status:              "active",
			Priority:            "high",
			CreatedBy:           &adminID,
			IsLimitedEdition:    false,
			IsGGLed:             false,
		},
		{
			UUID:                uuid.New().String(),
			Title:               "Flash Challenge: Social Media Blitz",
			Description:         stringPtr("Quick 24-hour challenge to boost social media presence"),
			CampaignType:        "flash",
			Category:            stringPtr("online"),
			BannerImageURL:      stringPtr("https://example.com/campaigns/flash-social.jpg"),
			StartDate:           now,
			EndDate:             now.AddDate(0, 0, 1),
			MaxParticipants:     intPtr(500),
			CurrentParticipants: 75,
			Status:              "active",
			Priority:            "flash",
			CreatedBy:           &adminID,
			IsLimitedEdition:    true,
			IsGGLed:             true,
		},
		{
			UUID:                uuid.New().String(),
			Title:               "Weekly Vibe Challenge",
			Description:         stringPtr("Weekly creative content challenge"),
			CampaignType:        "weekly_vibe",
			Category:            stringPtr("online"),
			BannerImageURL:      stringPtr("https://example.com/campaigns/weekly-vibe.jpg"),
			StartDate:           now.AddDate(0, 0, -7),
			EndDate:             now.AddDate(0, 0, 7),
			MaxParticipants:     nil,
			CurrentParticipants: 200,
			Status:              "active",
			Priority:            "medium",
			CreatedBy:           &adminID,
			IsLimitedEdition:    false,
			IsGGLed:             true,
		},
		{
			UUID:                uuid.New().String(),
			Title:               "Thematic Campaign: Sustainability",
			Description:         stringPtr("Promote sustainability and environmental awareness"),
			CampaignType:        "thematic",
			Category:            stringPtr("offline"),
			BannerImageURL:      stringPtr("https://example.com/campaigns/sustainability.jpg"),
			StartDate:           now.AddDate(0, -1, 0),
			EndDate:             now.AddDate(0, 1, 0),
			MaxParticipants:     intPtr(2000),
			CurrentParticipants: 450,
			Status:              "active",
			Priority:            "high",
			CreatedBy:           &adminID,
			IsLimitedEdition:    false,
			IsGGLed:             false,
		},
	}

	for _, campaign := range campaigns {
		var existing store.Campaign
		if err := db.Where("title = ?", campaign.Title).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&campaign).Error; err != nil {
					return fmt.Errorf("failed to create campaign %s: %w", campaign.Title, err)
				}
				log.Printf("Created campaign: %s", campaign.Title)
			} else {
				return err
			}
		}
	}
	return nil
}

// seedTasks seeds tasks
func seedTasks(db *gorm.DB) error {
	// Get admin user
	var admin store.User
	if err := db.Where("role = ?", "admin").First(&admin).Error; err != nil {
		return fmt.Errorf("admin user not found: %w", err)
	}

	adminID := int(admin.ID)

	// Get first campaign
	var campaign store.Campaign
	if err := db.First(&campaign).Error; err != nil {
		return fmt.Errorf("no campaigns found: %w", err)
	}

	campaignID := int(campaign.ID)
	durationHours := 72

	tasks := []store.Task{
		{
			UUID:                   uuid.New().String(),
			CampaignID:             &campaignID,
			Title:                  "Create Instagram Post",
			Description:            "Create an engaging Instagram post about our product and tag us",
			TaskType:               "online",
			ProofType:              "screenshot",
			XPReward:               500,
			CoinReward:             100,
			DurationHours:           &durationHours,
			Priority:               "medium",
			AssignmentType:         stringPtr("individual"),
			MaxSubmissions:         1,
			IsActive:               true,
			SubmissionInstructions: stringPtr("Upload a screenshot of your Instagram post"),
			CreatedBy:              &adminID,
		},
		{
			UUID:                   uuid.New().String(),
			CampaignID:             &campaignID,
			Title:                  "Share on LinkedIn",
			Description:            "Share our campaign announcement on LinkedIn with your network",
			TaskType:               "online",
			ProofType:              "url",
			XPReward:               300,
			CoinReward:             50,
			DurationHours:           &durationHours,
			Priority:               "low",
			AssignmentType:         stringPtr("individual"),
			MaxSubmissions:         1,
			IsActive:               true,
			SubmissionInstructions: stringPtr("Provide the LinkedIn post URL"),
			CreatedBy:              &adminID,
		},
		{
			UUID:                   uuid.New().String(),
			CampaignID:             nil, // Standalone task
			Title:                  "Campus Event Organization",
			Description:            "Organize a campus event with at least 50 participants",
			TaskType:               "offline",
			ProofType:              "pdf",
			XPReward:               2000,
			CoinReward:             500,
			DurationHours:           nil,
			Priority:               "high",
			AssignmentType:         stringPtr("college"),
			MaxSubmissions:         1,
			IsActive:               true,
			SubmissionInstructions: stringPtr("Upload event report as PDF"),
			CreatedBy:              &adminID,
		},
		{
			UUID:                   uuid.New().String(),
			CampaignID:             &campaignID,
			Title:                  "Create YouTube Video",
			Description:            "Create a 2-5 minute YouTube video review",
			TaskType:               "online",
			ProofType:              "video",
			XPReward:               1500,
			CoinReward:             300,
			DurationHours:           intPtr(168), // 7 days
			Priority:               "high",
			AssignmentType:         stringPtr("individual"),
			MaxSubmissions:         1,
			IsActive:               true,
			SubmissionInstructions: stringPtr("Provide YouTube video URL"),
			CreatedBy:              &adminID,
		},
	}

	for _, task := range tasks {
		var existing store.Task
		if err := db.Where("title = ?", task.Title).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&task).Error; err != nil {
					return fmt.Errorf("failed to create task %s: %w", task.Title, err)
				}
				log.Printf("Created task: %s", task.Title)
			} else {
				return err
			}
		}
	}
	return nil
}

// seedRewards seeds rewards store
func seedRewards(db *gorm.DB) error {
	rewards := []store.RewardStore{
		{
			Name:             "Amazon Gift Card ₹500",
			Description:      stringPtr("Redeemable Amazon gift card worth ₹500"),
			RewardType:       "gift_card",
			Category:         stringPtr("digital"),
			ImageURL:         stringPtr("https://example.com/rewards/amazon-500.jpg"),
			XPCost:           5000,
			CoinCost:         0,
			CashCost:         nil,
			QuantityAvailable: intPtr(100),
			QuantitySold:     0,
			IsFeatured:       true,
			IsActive:         true,
			ValidityDays:     intPtr(365),
		},
		{
			Name:             "Branded T-Shirt",
			Description:      stringPtr("Premium quality branded t-shirt"),
			RewardType:       "physical",
			Category:         stringPtr("merchandise"),
			ImageURL:         stringPtr("https://example.com/rewards/tshirt.jpg"),
			XPCost:           3000,
			CoinCost:         500,
			CashCost:         nil,
			QuantityAvailable: intPtr(50),
			QuantitySold:     5,
			IsFeatured:       true,
			IsActive:         true,
			ValidityDays:     nil,
		},
		{
			Name:             "XP Boost 2x (7 days)",
			Description:      stringPtr("Double your XP earnings for 7 days"),
			RewardType:       "xp_boost",
			Category:         stringPtr("digital"),
			ImageURL:         stringPtr("https://example.com/rewards/xp-boost.jpg"),
			XPCost:           2000,
			CoinCost:         0,
			CashCost:         nil,
			QuantityAvailable: nil,
			QuantitySold:     0,
			IsFeatured:       false,
			IsActive:         true,
			ValidityDays:     intPtr(7),
		},
		{
			Name:             "Premium Profile Skin",
			Description:      stringPtr("Exclusive profile skin design"),
			RewardType:       "profile_skin",
			Category:         stringPtr("digital"),
			ImageURL:         stringPtr("https://example.com/rewards/skin-premium.jpg"),
			XPCost:           1500,
			CoinCost:         200,
			CashCost:         nil,
			QuantityAvailable: nil,
			QuantitySold:     0,
			IsFeatured:       false,
			IsActive:         true,
			ValidityDays:     nil,
		},
		{
			Name:             "Certificate of Excellence",
			Description:      stringPtr("Digital certificate for outstanding performance"),
			RewardType:       "certificate",
			Category:         stringPtr("digital"),
			ImageURL:         stringPtr("https://example.com/rewards/certificate.jpg"),
			XPCost:           1000,
			CoinCost:         0,
			CashCost:         nil,
			QuantityAvailable: nil,
			QuantitySold:     0,
			IsFeatured:       false,
			IsActive:         true,
			ValidityDays:     nil,
		},
	}

	for _, reward := range rewards {
		var existing store.RewardStore
		if err := db.Where("name = ?", reward.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&reward).Error; err != nil {
					return fmt.Errorf("failed to create reward %s: %w", reward.Name, err)
				}
				log.Printf("Created reward: %s", reward.Name)
			} else {
				return err
			}
		}
	}
	return nil
}

// seedSpinWheel seeds spin wheel
func seedSpinWheel(db *gorm.DB) error {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 0, 7)

	wheel := store.SpinWheel{
		Name:            "Weekly Spin Wheel",
		Description:     stringPtr("Spin to win rewards every week"),
		WheelType:       "weekly",
		IsActive:        true,
		SpinsPerUser:    3,
		ResetFrequency:  "weekly",
		StartDate:       &startDate,
		EndDate:         &endDate,
		MinActivityLevel: 0,
	}

	var existing store.SpinWheel
	if err := db.Where("wheel_type = ?", "weekly").First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&wheel).Error; err != nil {
				return fmt.Errorf("failed to create spin wheel: %w", err)
			}
			log.Printf("Created spin wheel: %s", wheel.Name)

			// Create spin wheel items
			wheelID := intPtr(int(wheel.ID))
			items := []store.SpinWheelItem{
				{SpinWheelID: wheelID, ItemType: "xp", ItemValue: 100, ItemLabel: "100 XP", Probability: 0.30, MaxQuantity: nil, CurrentQuantity: nil, IsActive: true, SortOrder: 1},
				{SpinWheelID: wheelID, ItemType: "xp", ItemValue: 250, ItemLabel: "250 XP", Probability: 0.20, MaxQuantity: nil, CurrentQuantity: nil, IsActive: true, SortOrder: 2},
				{SpinWheelID: wheelID, ItemType: "xp", ItemValue: 500, ItemLabel: "500 XP", Probability: 0.15, MaxQuantity: nil, CurrentQuantity: nil, IsActive: true, SortOrder: 3},
				{SpinWheelID: wheelID, ItemType: "coins", ItemValue: 50, ItemLabel: "50 Coins", Probability: 0.20, MaxQuantity: nil, CurrentQuantity: nil, IsActive: true, SortOrder: 4},
				{SpinWheelID: wheelID, ItemType: "coins", ItemValue: 100, ItemLabel: "100 Coins", Probability: 0.10, MaxQuantity: nil, CurrentQuantity: nil, IsActive: true, SortOrder: 5},
				{SpinWheelID: wheelID, ItemType: "badge", ItemValue: 1, ItemLabel: "Mystery Badge", Probability: 0.04, MaxQuantity: intPtr(10), CurrentQuantity: intPtr(10), IsActive: true, SortOrder: 6},
				{SpinWheelID: wheelID, ItemType: "xp", ItemValue: 1000, ItemLabel: "1000 XP", Probability: 0.01, MaxQuantity: nil, CurrentQuantity: nil, IsActive: true, SortOrder: 7},
			}

			for _, item := range items {
				if err := db.Create(&item).Error; err != nil {
					return fmt.Errorf("failed to create spin wheel item: %w", err)
				}
			}
			log.Printf("Created %d spin wheel items", len(items))
		} else {
			return err
		}
	}
	return nil
}

// seedMysteryBoxes seeds mystery boxes
func seedMysteryBoxes(db *gorm.DB) error {
	boxes := []store.MysteryBox{
		{
			Name:        "Standard Mystery Box",
			Description: stringPtr("Contains random rewards"),
			CostXP:      1000,
			CostCoins:   0,
			Contents:    `{"xp": [100, 500], "coins": [50, 200], "badge_chance": 0.1}`,
			IsActive:    true,
		},
		{
			Name:        "Premium Mystery Box",
			Description: stringPtr("Higher chance of rare rewards"),
			CostXP:      2500,
			CostCoins:   100,
			Contents:    `{"xp": [500, 1500], "coins": [200, 500], "badge_chance": 0.3, "profile_skin_chance": 0.1}`,
			IsActive:    true,
		},
	}

	for _, box := range boxes {
		var existing store.MysteryBox
		if err := db.Where("name = ?", box.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&box).Error; err != nil {
					return fmt.Errorf("failed to create mystery box %s: %w", box.Name, err)
				}
				log.Printf("Created mystery box: %s", box.Name)
			} else {
				return err
			}
		}
	}
	return nil
}

// seedSecretCodes seeds secret codes
func seedSecretCodes(db *gorm.DB) error {
	now := time.Now()
	
	// Get admin user for CreatedBy
	var admin store.User
	if err := db.Where("role = ?", "admin").First(&admin).Error; err != nil {
		return fmt.Errorf("admin user not found: %w", err)
	}
	adminID := int(admin.ID)

	codes := []store.SecretCode{
		{
			Code:              "WELCOME2024",
			Description:       stringPtr("Welcome bonus code"),
			XPReward:          500,
			CoinReward:        0,
			BadgeID:           nil,
			MaxRedemptions:    1000,
			CurrentRedemptions: 0,
			ValidFrom:         now,
			ValidUntil:        now.AddDate(1, 0, 0),
			DistributionChannel: stringPtr("email"),
			IsActive:          true,
			CreatedBy:         &adminID,
		},
		{
			Code:              "FIRST100",
			Description:       stringPtr("First 100 users code"),
			XPReward:          0,
			CoinReward:        100,
			BadgeID:           nil,
			MaxRedemptions:    100,
			CurrentRedemptions: 25,
			ValidFrom:         now,
			ValidUntil:        now.AddDate(0, 6, 0),
			DistributionChannel: stringPtr("special"),
			IsActive:          true,
			CreatedBy:         &adminID,
		},
		{
			Code:              "FLASH50",
			Description:       stringPtr("Flash challenge bonus"),
			XPReward:          50,
			CoinReward:        0,
			BadgeID:           nil,
			MaxRedemptions:    0, // Unlimited
			CurrentRedemptions: 0,
			ValidFrom:         now,
			ValidUntil:        now.AddDate(0, 0, 7),
			DistributionChannel: stringPtr("campaign"),
			IsActive:          true,
			CreatedBy:         &adminID,
		},
	}

	for _, code := range codes {
		var existing store.SecretCode
		if err := db.Where("code = ?", code.Code).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&code).Error; err != nil {
					return fmt.Errorf("failed to create secret code %s: %w", code.Code, err)
				}
				log.Printf("Created secret code: %s", code.Code)
			} else {
				return err
			}
		}
	}
	return nil
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func generateReferralCode() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
