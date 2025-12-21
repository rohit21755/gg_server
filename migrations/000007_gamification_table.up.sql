-- Streak Tracking
CREATE TABLE user_streaks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE,
    streak_type VARCHAR(50) NOT NULL CHECK (streak_type IN ('daily_login', 'weekly_task', 'campaign')),
    current_streak INTEGER DEFAULT 0,
    longest_streak INTEGER DEFAULT 0,
    last_activity_date DATE NOT NULL,
    total_days INTEGER DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Streak Log
CREATE TABLE streak_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    streak_type VARCHAR(50) NOT NULL,
    activity_date DATE NOT NULL,
    earned_xp INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, streak_type, activity_date)
);

-- Add foreign keys
ALTER TABLE user_streaks 
ADD CONSTRAINT fk_user_streaks_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE streak_logs 
ADD CONSTRAINT fk_streak_logs_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;