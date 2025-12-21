-- Levels Table
CREATE TABLE levels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    rank_order INTEGER NOT NULL UNIQUE,
    min_xp INTEGER NOT NULL,
    max_xp INTEGER,
    badge_url TEXT,
    description TEXT,
    unlockable_features JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Default levels: Rookie->Performer->Pro->Ace->Legend->Champion
INSERT INTO levels (name, rank_order, min_xp, description) VALUES
('Rookie', 1, 0, 'Just starting out'),
('Performer', 2, 1000, 'Consistent performer'),
('Pro', 3, 5000, 'Pro level ambassador'),
('Ace', 4, 15000, 'Top tier performer'),
('Legend', 5, 30000, 'Ambassador legend'),
('Champion', 6, 50000, 'Ultimate champion');

-- Badges Table
CREATE TABLE badges (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    badge_type VARCHAR(50) NOT NULL,
    category VARCHAR(50) CHECK (category IN ('referral', 'submission', 'campaign', 'streak', 'special')),
    image_url TEXT NOT NULL,
    xp_reward INTEGER DEFAULT 0,
    criteria_type VARCHAR(50) NOT NULL,
    criteria_value INTEGER NOT NULL,
    is_secret BOOLEAN DEFAULT false,
    is_limited_edition BOOLEAN DEFAULT false,
    available_until TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Profile Skins (without foreign key initially)
CREATE TABLE profile_skins (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    preview_url TEXT NOT NULL,
    unlock_method VARCHAR(50) CHECK (unlock_method IN ('xp', 'campaign', 'purchase', 'special')),
    xp_cost INTEGER DEFAULT 0,
    campaign_id INTEGER,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add level foreign key to users
ALTER TABLE users 
ADD CONSTRAINT fk_users_level 
FOREIGN KEY (level_id) REFERENCES levels(id);