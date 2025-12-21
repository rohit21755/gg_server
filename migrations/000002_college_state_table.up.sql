-- States Table
CREATE TABLE states (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    code VARCHAR(10) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Colleges Table
CREATE TABLE colleges (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    state_id INTEGER,
    code VARCHAR(50) UNIQUE,
    total_cas INTEGER DEFAULT 0,
    total_xp INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Now add foreign keys to users table
ALTER TABLE users 
ADD CONSTRAINT fk_users_college 
FOREIGN KEY (college_id) REFERENCES colleges(id),
ADD CONSTRAINT fk_users_state 
FOREIGN KEY (state_id) REFERENCES states(id),
ADD CONSTRAINT fk_users_referred_by 
FOREIGN KEY (referred_by) REFERENCES users(id);

ALTER TABLE user_sessions 
ADD CONSTRAINT fk_user_sessions_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;