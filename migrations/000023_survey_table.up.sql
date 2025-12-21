-- Surveys
CREATE TABLE surveys (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    survey_type VARCHAR(50) CHECK (survey_type IN ('feedback', 'quiz', 'research', 'poll')),
    questions JSONB NOT NULL,
    xp_reward INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Survey Responses
CREATE TABLE survey_responses (
    id SERIAL PRIMARY KEY,
    survey_id INTEGER,
    user_id INTEGER,
    responses JSONB NOT NULL,
    completion_percentage INTEGER DEFAULT 100,
    xp_awarded INTEGER DEFAULT 0,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign keys and constraints
ALTER TABLE surveys 
ADD CONSTRAINT fk_surveys_created_by 
FOREIGN KEY (created_by) REFERENCES users(id);

ALTER TABLE survey_responses 
ADD CONSTRAINT fk_survey_responses_survey 
FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_survey_responses_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE survey_responses 
ADD UNIQUE (survey_id, user_id);