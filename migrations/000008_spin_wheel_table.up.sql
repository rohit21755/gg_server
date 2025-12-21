-- Spin Wheel
CREATE TABLE spin_wheels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    wheel_type VARCHAR(50) DEFAULT 'weekly' CHECK (wheel_type IN ('weekly', 'daily', 'special')),
    is_active BOOLEAN DEFAULT true,
    spins_per_user INTEGER DEFAULT 1,
    reset_frequency VARCHAR(20) DEFAULT 'weekly',
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    min_activity_level INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Spin Wheel Items
CREATE TABLE spin_wheel_items (
    id SERIAL PRIMARY KEY,
    spin_wheel_id INTEGER,
    item_type VARCHAR(50) NOT NULL CHECK (item_type IN ('xp', 'coins', 'badge', 'physical', 'discount')),
    item_value INTEGER NOT NULL,
    item_label VARCHAR(100) NOT NULL,
    probability DECIMAL(5,4) NOT NULL,
    max_quantity INTEGER,
    current_quantity INTEGER,
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0
);

-- User Spins
CREATE TABLE user_spins (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    spin_wheel_id INTEGER,
    spin_wheel_item_id INTEGER,
    earned_value INTEGER NOT NULL,
    spun_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Mystery Box
CREATE TABLE mystery_boxes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    cost_xp INTEGER NOT NULL,
    cost_coins INTEGER DEFAULT 0,
    contents JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Mystery Box Redemptions
CREATE TABLE mystery_box_redemptions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    mystery_box_id INTEGER,
    reward_type VARCHAR(50) NOT NULL,
    reward_value INTEGER NOT NULL,
    redeemed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign keys
ALTER TABLE spin_wheel_items 
ADD CONSTRAINT fk_spin_wheel_items_wheel 
FOREIGN KEY (spin_wheel_id) REFERENCES spin_wheels(id) ON DELETE CASCADE;

ALTER TABLE user_spins 
ADD CONSTRAINT fk_user_spins_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_user_spins_wheel 
FOREIGN KEY (spin_wheel_id) REFERENCES spin_wheels(id),
ADD CONSTRAINT fk_user_spins_item 
FOREIGN KEY (spin_wheel_item_id) REFERENCES spin_wheel_items(id);

ALTER TABLE mystery_box_redemptions 
ADD CONSTRAINT fk_mystery_box_redemptions_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_mystery_box_redemptions_box 
FOREIGN KEY (mystery_box_id) REFERENCES mystery_boxes(id);