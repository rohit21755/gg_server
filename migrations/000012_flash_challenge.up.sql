-- Flash Challenges
CREATE TABLE flash_challenges (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    challenge_type VARCHAR(50) CHECK (challenge_type IN ('meme', 'reel', 'video', 'qr_scan', 'content')),
    duration_hours INTEGER NOT NULL,
    xp_reward INTEGER NOT NULL,
    max_participants INTEGER,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'active', 'completed')),
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign key
ALTER TABLE flash_challenges 
ADD CONSTRAINT fk_flash_challenges_created_by 
FOREIGN KEY (created_by) REFERENCES users(id);