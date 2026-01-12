-- Insert default badges
-- Using WHERE NOT EXISTS to avoid duplicates if badges already exist
INSERT INTO badges (name, description, badge_type, category, image_url, xp_reward, criteria_type, criteria_value, is_secret, is_limited_edition)
SELECT 'First Steps', 'Complete your first task', 'achievement', 'submission', 'https://example.com/badges/first-steps.png', 100, 'submission_count', 1, false, false
WHERE NOT EXISTS (SELECT 1 FROM badges WHERE name = 'First Steps');

INSERT INTO badges (name, description, badge_type, category, image_url, xp_reward, criteria_type, criteria_value, is_secret, is_limited_edition)
SELECT 'Referral Master', 'Refer 10 friends', 'achievement', 'referral', 'https://example.com/badges/referral-master.png', 500, 'referral_count', 10, false, false
WHERE NOT EXISTS (SELECT 1 FROM badges WHERE name = 'Referral Master');

INSERT INTO badges (name, description, badge_type, category, image_url, xp_reward, criteria_type, criteria_value, is_secret, is_limited_edition)
SELECT 'Campaign Champion', 'Complete 5 campaigns', 'achievement', 'campaign', 'https://example.com/badges/campaign-champion.png', 1000, 'campaign_count', 5, false, false
WHERE NOT EXISTS (SELECT 1 FROM badges WHERE name = 'Campaign Champion');

INSERT INTO badges (name, description, badge_type, category, image_url, xp_reward, criteria_type, criteria_value, is_secret, is_limited_edition)
SELECT 'Streak Warrior', 'Maintain a 30-day streak', 'achievement', 'streak', 'https://example.com/badges/streak-warrior.png', 750, 'streak_days', 30, false, false
WHERE NOT EXISTS (SELECT 1 FROM badges WHERE name = 'Streak Warrior');

INSERT INTO badges (name, description, badge_type, category, image_url, xp_reward, criteria_type, criteria_value, is_secret, is_limited_edition)
SELECT 'Hidden Gem', 'Secret achievement badge', 'special', 'special', 'https://example.com/badges/hidden-gem.png', 2000, 'special', 1, true, false
WHERE NOT EXISTS (SELECT 1 FROM badges WHERE name = 'Hidden Gem');

INSERT INTO badges (name, description, badge_type, category, image_url, xp_reward, criteria_type, criteria_value, is_secret, is_limited_edition)
SELECT 'Early Bird', 'Join in the first month', 'limited', 'special', 'https://example.com/badges/early-bird.png', 300, 'join_date', 1, false, true
WHERE NOT EXISTS (SELECT 1 FROM badges WHERE name = 'Early Bird');
