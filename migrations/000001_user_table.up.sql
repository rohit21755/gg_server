-- Users Table (must be first)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) UNIQUE NOT NULL DEFAULT gen_random_uuid()::text,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL DEFAULT 'ca' CHECK (role IN ('admin', 'state_lead', 'ca')),
    college_id INTEGER,
    state_id INTEGER,
    referral_code VARCHAR(20) UNIQUE NOT NULL,
    referred_by INTEGER,
    xp INTEGER NOT NULL DEFAULT 0,
    level_id INTEGER DEFAULT 1,
    streak_count INTEGER DEFAULT 0,
    last_login_date DATE,
    total_submissions INTEGER DEFAULT 0,
    approved_submissions INTEGER DEFAULT 0,
    win_rate DECIMAL(5,2) DEFAULT 0,
    profile_skin_id INTEGER,
    avatar_url TEXT,
    resume_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User Sessions (for WS server)
CREATE TABLE user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    device_id VARCHAR(255),
    platform VARCHAR(50),
    last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);