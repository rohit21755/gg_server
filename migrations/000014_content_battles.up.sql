-- Content Battles
CREATE TABLE content_battles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    battle_type VARCHAR(50) CHECK (battle_type IN ('meme', 'video', 'reel', 'post')),
    theme VARCHAR(200),
    submission_deadline TIMESTAMP NOT NULL,
    voting_start TIMESTAMP NOT NULL,
    voting_end TIMESTAMP NOT NULL,
    max_participants INTEGER,
    rewards JSONB DEFAULT '{}'::jsonb,
    status VARCHAR(20) DEFAULT 'upcoming' CHECK (status IN ('upcoming', 'submissions', 'voting', 'completed')),
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Battle Submissions
CREATE TABLE battle_submissions (
    id SERIAL PRIMARY KEY,
    battle_id INTEGER,
    user_id INTEGER,
    title VARCHAR(200),
    description TEXT,
    media_url TEXT NOT NULL,
    thumbnail_url TEXT,
    vote_count INTEGER DEFAULT 0,
    rank INTEGER,
    is_winner BOOLEAN DEFAULT false,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Battle Votes
CREATE TABLE battle_votes (
    id SERIAL PRIMARY KEY,
    battle_id INTEGER,
    submission_id INTEGER,
    voter_id INTEGER,
    voted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign keys and constraints
ALTER TABLE content_battles 
ADD CONSTRAINT fk_content_battles_created_by 
FOREIGN KEY (created_by) REFERENCES users(id);

ALTER TABLE battle_submissions 
ADD CONSTRAINT fk_battle_submissions_battle 
FOREIGN KEY (battle_id) REFERENCES content_battles(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_battle_submissions_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE battle_votes 
ADD CONSTRAINT fk_battle_votes_battle 
FOREIGN KEY (battle_id) REFERENCES content_battles(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_battle_votes_submission 
FOREIGN KEY (submission_id) REFERENCES battle_submissions(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_battle_votes_voter 
FOREIGN KEY (voter_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE battle_submissions 
ADD UNIQUE (battle_id, user_id);

ALTER TABLE battle_votes 
ADD UNIQUE (battle_id, voter_id, submission_id);