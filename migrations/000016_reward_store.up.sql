-- Rewards Store
CREATE TABLE rewards_store (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    reward_type VARCHAR(50) NOT NULL CHECK (reward_type IN ('physical', 'digital', 'badge', 'xp_boost', 'profile_skin', 'certificate', 'gift_card')),
    category VARCHAR(50),
    image_url TEXT,
    xp_cost INTEGER DEFAULT 0,
    coin_cost INTEGER DEFAULT 0,
    cash_cost DECIMAL(10,2),
    quantity_available INTEGER,
    quantity_sold INTEGER DEFAULT 0,
    is_featured BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    validity_days INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User Rewards/Redemptions
CREATE TABLE user_rewards (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    reward_id INTEGER,
    redemption_code VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled', 'expired')),
    xp_paid INTEGER DEFAULT 0,
    coins_paid INTEGER DEFAULT 0,
    cash_paid DECIMAL(10,2) DEFAULT 0,
    shipping_address JSONB,
    tracking_number VARCHAR(100),
    claimed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delivered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign keys
ALTER TABLE user_rewards 
ADD CONSTRAINT fk_user_rewards_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_user_rewards_reward 
FOREIGN KEY (reward_id) REFERENCES rewards_store(id);