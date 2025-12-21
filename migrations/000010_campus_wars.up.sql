-- Campus Wars/State vs State
CREATE TABLE campus_wars (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    war_type VARCHAR(50) CHECK (war_type IN ('campus_vs_campus', 'state_vs_state', 'region_vs_region')),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    status VARCHAR(20) DEFAULT 'upcoming' CHECK (status IN ('upcoming', 'active', 'completed')),
    metrics JSONB DEFAULT '{"xp": true, "submissions": true, "referrals": true}'::jsonb,
    rewards JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- War Participants (Colleges/States)
CREATE TABLE war_participants (
    id SERIAL PRIMARY KEY,
    war_id INTEGER,
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN ('college', 'state')),
    entity_id INTEGER NOT NULL,
    total_xp INTEGER DEFAULT 0,
    total_submissions INTEGER DEFAULT 0,
    total_referrals INTEGER DEFAULT 0,
    rank INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- War Leaderboard Snapshots
CREATE TABLE war_leaderboard_snapshots (
    id SERIAL PRIMARY KEY,
    war_id INTEGER,
    snapshot_date DATE NOT NULL,
    leaderboard_data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign keys and constraints
ALTER TABLE war_participants 
ADD CONSTRAINT fk_war_participants_war 
FOREIGN KEY (war_id) REFERENCES campus_wars(id) ON DELETE CASCADE;

ALTER TABLE war_leaderboard_snapshots 
ADD CONSTRAINT fk_war_leaderboard_snapshots_war 
FOREIGN KEY (war_id) REFERENCES campus_wars(id) ON DELETE CASCADE;

ALTER TABLE war_participants 
ADD UNIQUE (war_id, entity_type, entity_id);

ALTER TABLE war_leaderboard_snapshots 
ADD UNIQUE (war_id, snapshot_date);