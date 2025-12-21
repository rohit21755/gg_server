-- Submissions Table
CREATE TABLE submissions (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) UNIQUE NOT NULL DEFAULT gen_random_uuid()::text,
    task_id INTEGER,
    user_id INTEGER,
    campaign_id INTEGER,
    proof_type VARCHAR(50) NOT NULL,
    proof_url TEXT NOT NULL,
    proof_text TEXT,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('draft', 'pending', 'under_review', 'approved', 'rejected', 'needs_revision')),
    submission_stage VARCHAR(20) DEFAULT 'initial' CHECK (submission_stage IN ('initial', 'resubmission')),
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP,
    reviewed_by INTEGER,
    review_comments TEXT,
    xp_awarded INTEGER DEFAULT 0,
    coins_awarded INTEGER DEFAULT 0,
    is_winner BOOLEAN DEFAULT false,
    score DECIMAL(5,2),
    revision_deadline TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Submission Media
CREATE TABLE submission_media (
    id SERIAL PRIMARY KEY,
    submission_id INTEGER,
    media_type VARCHAR(50) NOT NULL CHECK (media_type IN ('image', 'video', 'document', 'link')),
    media_url TEXT NOT NULL,
    thumbnail_url TEXT,
    file_name VARCHAR(255),
    file_size INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign keys
ALTER TABLE submissions 
ADD CONSTRAINT fk_submissions_task 
FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_submissions_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_submissions_campaign 
FOREIGN KEY (campaign_id) REFERENCES campaigns(id),
ADD CONSTRAINT fk_submissions_reviewed_by 
FOREIGN KEY (reviewed_by) REFERENCES users(id);

ALTER TABLE submission_media 
ADD CONSTRAINT fk_submission_media_submission 
FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE;