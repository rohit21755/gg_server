-- Quest Lines
CREATE TABLE quests (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    quest_type VARCHAR(50) CHECK (quest_type IN ('multi_step', 'achievement', 'special')),
    steps JSONB NOT NULL,
    rewards JSONB DEFAULT '{}'::jsonb,
    time_limit_days INTEGER,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User Quest Progress
CREATE TABLE user_quests (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    quest_id INTEGER,
    current_step INTEGER DEFAULT 1,
    progress_data JSONB DEFAULT '{}'::jsonb,
    status VARCHAR(20) DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed', 'failed', 'abandoned')),
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign keys and constraints
ALTER TABLE user_quests 
ADD CONSTRAINT fk_user_quests_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_user_quests_quest 
FOREIGN KEY (quest_id) REFERENCES quests(id);

ALTER TABLE user_quests 
ADD UNIQUE (user_id, quest_id);