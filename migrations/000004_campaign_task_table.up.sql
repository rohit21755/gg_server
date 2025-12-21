-- Campaigns Table
CREATE TABLE campaigns (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) UNIQUE NOT NULL DEFAULT gen_random_uuid()::text,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    campaign_type VARCHAR(50) NOT NULL CHECK (campaign_type IN ('brand_specific', 'thematic', 'seasonal', 'gg_led', 'flash', 'weekly_vibe', 'limited_edition')),
    category VARCHAR(50) CHECK (category IN ('solo', 'group', 'online', 'offline')),
    banner_image_url TEXT,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    max_participants INTEGER,
    current_participants INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'paused', 'completed', 'cancelled')),
    priority VARCHAR(20) DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
    created_by INTEGER,
    is_limited_edition BOOLEAN DEFAULT false,
    is_gg_led BOOLEAN DEFAULT false,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tasks Table
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) UNIQUE NOT NULL DEFAULT gen_random_uuid()::text,
    campaign_id INTEGER,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    task_type VARCHAR(50) NOT NULL CHECK (task_type IN ('solo', 'group', 'online', 'offline')),
    proof_type VARCHAR(50) NOT NULL CHECK (proof_type IN ('screenshot', 'url', 'pdf', 'video', 'text')),
    xp_reward INTEGER NOT NULL DEFAULT 0,
    coin_reward INTEGER DEFAULT 0,
    duration_hours INTEGER,
    priority VARCHAR(20) DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high', 'flash')),
    assignment_type VARCHAR(50) CHECK (assignment_type IN ('role', 'college', 'state', 'individual')),
    assignment_target JSONB DEFAULT '{}'::jsonb,
    max_submissions INTEGER DEFAULT 1,
    is_active BOOLEAN DEFAULT true,
    submission_instructions TEXT,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Task Assignment by Role
CREATE TABLE task_assignments (
    id SERIAL PRIMARY KEY,
    task_id INTEGER,
    assignee_type VARCHAR(50) NOT NULL CHECK (assignee_type IN ('user', 'role', 'college', 'state')),
    assignee_id INTEGER,
    assignee_role VARCHAR(20),
    assigned_by INTEGER,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'assigned' CHECK (status IN ('assigned', 'accepted', 'declined', 'completed'))
);

-- Now add foreign keys
ALTER TABLE campaigns 
ADD CONSTRAINT fk_campaigns_created_by 
FOREIGN KEY (created_by) REFERENCES users(id);

ALTER TABLE tasks 
ADD CONSTRAINT fk_tasks_campaign 
FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_tasks_created_by 
FOREIGN KEY (created_by) REFERENCES users(id);

ALTER TABLE task_assignments 
ADD CONSTRAINT fk_task_assignments_task 
FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_task_assignments_assigner 
FOREIGN KEY (assigned_by) REFERENCES users(id);

-- Add foreign key to profile_skins now that campaigns exists
ALTER TABLE profile_skins 
ADD CONSTRAINT fk_profile_skins_campaign 
FOREIGN KEY (campaign_id) REFERENCES campaigns(id);