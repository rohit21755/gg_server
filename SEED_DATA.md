# Seed Data Documentation

This document describes the seed data that gets populated when you run the seed command.

## Running the Seed

```bash
make seed
```

Or directly:
```bash
go run cmd/seed/main.go
```

## What Gets Seeded

### 1. States (15 Indian States)
- Maharashtra, Karnataka, Tamil Nadu, Delhi, Gujarat
- Rajasthan, West Bengal, Uttar Pradesh, Telangana, Kerala
- Punjab, Haryana, Madhya Pradesh, Bihar, Andhra Pradesh

### 2. Colleges (15 Colleges)
Colleges across different states including:
- IITs (Bombay, Delhi, Madras, Bangalore)
- NITs (Surathkal, Calicut)
- VIT, BITS, Manipal, SRM, Amity, Symbiosis, Christ, Jadavpur, DTU

Each college has:
- Name, State association, Unique code
- Total XP and CA count

### 3. Users (2 Users)

#### Admin User
- **Email**: `admin@campusambassador.com`
- **Password**: `password123`
- **Role**: `admin`
- **XP**: 100,000
- **Level**: Highest level (Champion)

#### Regular User
- **Email**: `user@campusambassador.com`
- **Password**: `password123`
- **Role**: `ca` (Campus Ambassador)
- **XP**: 2,500
- **Level**: Rookie (Level 1)

Both users are:
- Associated with the first college and state
- Have unique referral codes
- Active accounts

### 4. Badges (6 Badges)
- **First Steps**: Complete first task (100 XP)
- **Referral Master**: Refer 10 friends (500 XP)
- **Campaign Champion**: Complete 5 campaigns (1000 XP)
- **Streak Warrior**: 30-day streak (750 XP)
- **Hidden Gem**: Secret achievement (2000 XP)
- **Early Bird**: Limited edition badge (300 XP)

### 5. Campaigns (4 Campaigns)
- **Summer Brand Campaign 2024**: Active brand-specific campaign
- **Flash Challenge: Social Media Blitz**: 24-hour flash challenge
- **Weekly Vibe Challenge**: Weekly creative content challenge
- **Thematic Campaign: Sustainability**: Sustainability awareness campaign

### 6. Tasks (4 Tasks)
- **Create Instagram Post**: 500 XP, 100 coins
- **Share on LinkedIn**: 300 XP, 50 coins
- **Campus Event Organization**: 2000 XP, 500 coins
- **Create YouTube Video**: 1500 XP, 300 coins

### 7. Rewards Store (5 Rewards)
- **Amazon Gift Card ₹500**: 5000 XP
- **Branded T-Shirt**: 3000 XP + 500 coins
- **XP Boost 2x (7 days)**: 2000 XP
- **Premium Profile Skin**: 1500 XP + 200 coins
- **Certificate of Excellence**: 1000 XP

### 8. Spin Wheel
- **Weekly Spin Wheel**: Active weekly spin wheel
- **7 Items**: XP rewards (100, 250, 500, 1000), Coins (50, 100), Mystery Badge
- **3 spins per user per week**

### 9. Mystery Boxes (2 Boxes)
- **Standard Mystery Box**: 1000 XP cost
- **Premium Mystery Box**: 2500 XP + 100 coins cost

### 10. Secret Codes (3 Codes)
- **WELCOME2024**: 500 XP (1000 uses, valid 1 year)
- **FIRST100**: 100 coins (100 uses, 25 already used, valid 6 months)
- **FLASH50**: 50 XP (unlimited, valid 7 days)

## Notes

- Seed data is idempotent - running it multiple times won't create duplicates
- Existing records are skipped (checked by unique fields)
- All seed data uses realistic values and relationships
- Users can log in with the credentials above for testing

## Testing with Seed Data

After seeding, you can:

1. **Login as Admin**:
   ```
   POST /api/v1/auth/login
   {
     "email": "admin@campusambassador.com",
     "password": "password123"
   }
   ```

2. **Login as User**:
   ```
   POST /api/v1/auth/login
   {
     "email": "user@campusambassador.com",
     "password": "password123"
   }
   ```

3. **Redeem Secret Code**:
   ```
   POST /api/v1/secret-codes/redeem/WELCOME2024
   ```

4. **View Campaigns**:
   ```
   GET /api/v1/campaigns
   ```

5. **View Tasks**:
   ```
   GET /api/v1/tasks
   ```

## Resetting Seed Data

To reset and re-seed:
1. Drop and recreate your database
2. Run migrations: `make migrate-up`
3. Run seed: `make seed`
