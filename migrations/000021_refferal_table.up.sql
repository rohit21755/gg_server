-- Referrals
CREATE TABLE referrals (
    id SERIAL PRIMARY KEY,
    referrer_id INTEGER,
    referred_email VARCHAR(255) NOT NULL,
    referred_user_id INTEGER,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'joined', 'completed_task', 'converted')),
    xp_awarded INTEGER DEFAULT 0,
    xp_awarded_to_referred INTEGER DEFAULT 0,
    conversion_stage INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign keys and constraints
ALTER TABLE referrals 
ADD CONSTRAINT fk_referrals_referrer 
FOREIGN KEY (referrer_id) REFERENCES users(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_referrals_referred_user 
FOREIGN KEY (referred_user_id) REFERENCES users(id);

ALTER TABLE referrals 
ADD UNIQUE (referrer_id, referred_email);