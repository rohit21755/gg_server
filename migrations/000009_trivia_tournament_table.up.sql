-- Trivia/Tournaments
CREATE TABLE trivia_tournaments (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    tournament_type VARCHAR(50) DEFAULT 'weekly' CHECK (tournament_type IN ('weekly', 'monthly', 'special')),
    questions JSONB NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    duration_minutes INTEGER DEFAULT 10,
    max_participants INTEGER,
    entry_fee_xp INTEGER DEFAULT 0,
    rewards JSONB DEFAULT '{}'::jsonb,
    status VARCHAR(20) DEFAULT 'upcoming' CHECK (status IN ('upcoming', 'active', 'completed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trivia Participants
CREATE TABLE trivia_participants (
    id SERIAL PRIMARY KEY,
    trivia_id INTEGER,
    user_id INTEGER,
    score INTEGER DEFAULT 0,
    correct_answers INTEGER DEFAULT 0,
    time_taken_seconds INTEGER,
    rank INTEGER,
    reward_claimed BOOLEAN DEFAULT false,
    participated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign keys and unique constraint
ALTER TABLE trivia_participants 
ADD CONSTRAINT fk_trivia_participants_trivia 
FOREIGN KEY (trivia_id) REFERENCES trivia_tournaments(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_trivia_participants_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE trivia_participants 
ADD UNIQUE (trivia_id, user_id);