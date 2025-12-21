-- Badge Bingo
CREATE TABLE badge_bingo (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    bingo_card JSONB NOT NULL,
    bundle_rewards JSONB DEFAULT '{}'::jsonb,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User Bingo Progress
CREATE TABLE user_bingo_progress (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    bingo_id INTEGER,
    completed_cells JSONB DEFAULT '[]'::jsonb,
    completed_rows INTEGER DEFAULT 0,
    completed_columns INTEGER DEFAULT 0,
    completed_diagonals INTEGER DEFAULT 0,
    rewards_claimed JSONB DEFAULT '[]'::jsonb,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign keys and constraints
ALTER TABLE user_bingo_progress 
ADD CONSTRAINT fk_user_bingo_progress_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_user_bingo_progress_bingo 
FOREIGN KEY (bingo_id) REFERENCES badge_bingo(id) ON DELETE CASCADE;

ALTER TABLE user_bingo_progress 
ADD UNIQUE (user_id, bingo_id);