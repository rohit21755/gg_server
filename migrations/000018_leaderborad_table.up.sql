-- Leaderboards
CREATE TABLE leaderboards (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    leaderboard_type VARCHAR(50) NOT NULL CHECK (leaderboard_type IN ('global', 'state', 'college', 'weekly', 'monthly', 'campaign', 'war')),
    entity_type VARCHAR(50) CHECK (entity_type IN ('user', 'college', 'state')),
    entity_id INTEGER,
    period_start DATE,
    period_end DATE,
    metrics JSONB DEFAULT '{"xp": true}'::jsonb,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Leaderboard Entries
CREATE TABLE leaderboard_entries (
    id SERIAL PRIMARY KEY,
    leaderboard_id INTEGER,
    user_id INTEGER,
    college_id INTEGER,
    state_id INTEGER,
    xp INTEGER DEFAULT 0,
    submissions_count INTEGER DEFAULT 0,
    referrals_count INTEGER DEFAULT 0,
    win_rate DECIMAL(5,2) DEFAULT 0,
    rank INTEGER,
    previous_rank INTEGER,
    trend VARCHAR(10) CHECK (trend IN ('up', 'down', 'stable', 'new')),
    snapshot_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign keys and constraints
ALTER TABLE leaderboards 
ADD UNIQUE (name, leaderboard_type, period_start, period_end);

ALTER TABLE leaderboard_entries 
ADD CONSTRAINT fk_leaderboard_entries_leaderboard 
FOREIGN KEY (leaderboard_id) REFERENCES leaderboards(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_leaderboard_entries_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_leaderboard_entries_college 
FOREIGN KEY (college_id) REFERENCES colleges(id),
ADD CONSTRAINT fk_leaderboard_entries_state 
FOREIGN KEY (state_id) REFERENCES states(id);

-- Add unique constraints with conditional
CREATE UNIQUE INDEX unique_user_leaderboard_entry 
ON leaderboard_entries (leaderboard_id, user_id, snapshot_date) 
WHERE user_id IS NOT NULL;

CREATE UNIQUE INDEX unique_college_leaderboard_entry 
ON leaderboard_entries (leaderboard_id, college_id, snapshot_date) 
WHERE college_id IS NOT NULL;

CREATE UNIQUE INDEX unique_state_leaderboard_entry 
ON leaderboard_entries (leaderboard_id, state_id, snapshot_date) 
WHERE state_id IS NOT NULL;