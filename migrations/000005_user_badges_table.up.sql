-- User Badges (can only be created after users and badges exist)
CREATE TABLE user_badges (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    badge_id INTEGER REFERENCES badges(id),
    earned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, badge_id)
);

-- Also add foreign key for profile_skin_id in users
ALTER TABLE users 
ADD CONSTRAINT fk_users_profile_skin 
FOREIGN KEY (profile_skin_id) REFERENCES profile_skins(id);