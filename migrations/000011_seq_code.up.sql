-- Secret Code Hunts
CREATE TABLE secret_codes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    xp_reward INTEGER NOT NULL,
    coin_reward INTEGER DEFAULT 0,
    badge_id INTEGER,
    max_redemptions INTEGER DEFAULT 1,
    current_redemptions INTEGER DEFAULT 0,
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    distribution_channel VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Secret Code Redemptions
CREATE TABLE secret_code_redemptions (
    id SERIAL PRIMARY KEY,
    secret_code_id INTEGER,
    user_id INTEGER,
    redeemed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address INET
);

-- Add foreign keys and constraints
ALTER TABLE secret_codes 
ADD CONSTRAINT fk_secret_codes_badge 
FOREIGN KEY (badge_id) REFERENCES badges(id),
ADD CONSTRAINT fk_secret_codes_created_by 
FOREIGN KEY (created_by) REFERENCES users(id);

ALTER TABLE secret_code_redemptions 
ADD CONSTRAINT fk_secret_code_redemptions_code 
FOREIGN KEY (secret_code_id) REFERENCES secret_codes(id),
ADD CONSTRAINT fk_secret_code_redemptions_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE secret_code_redemptions 
ADD UNIQUE (secret_code_id, user_id);