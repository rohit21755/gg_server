-- Users indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_referral_code ON users(referral_code);
CREATE INDEX idx_users_college_id ON users(college_id);
CREATE INDEX idx_users_state_id ON users(state_id);
CREATE INDEX idx_users_level_id ON users(level_id);
CREATE INDEX idx_users_xp ON users(xp DESC);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- User sessions
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);
CREATE INDEX idx_user_sessions_token ON user_sessions(session_token);

-- Colleges & States
CREATE INDEX idx_colleges_state_id ON colleges(state_id);
CREATE INDEX idx_colleges_code ON colleges(code);
CREATE INDEX idx_states_code ON states(code);

-- Levels & Badges
CREATE INDEX idx_levels_rank_order ON levels(rank_order);
CREATE INDEX idx_badges_category ON badges(category);
CREATE INDEX idx_user_badges_user_id ON user_badges(user_id);
CREATE INDEX idx_user_badges_badge_id ON user_badges(badge_id);

-- Campaigns & Tasks
CREATE INDEX idx_campaigns_status ON campaigns(status);
CREATE INDEX idx_campaigns_start_date ON campaigns(start_date);
CREATE INDEX idx_campaigns_end_date ON campaigns(end_date);
CREATE INDEX idx_campaigns_type ON campaigns(campaign_type);
CREATE INDEX idx_campaigns_created_by ON campaigns(created_by);

CREATE INDEX idx_tasks_campaign_id ON tasks(campaign_id);
CREATE INDEX idx_tasks_status ON tasks(is_active);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_created_by ON tasks(created_by);
CREATE INDEX idx_tasks_created_at ON tasks(created_at DESC);

CREATE INDEX idx_task_assignments_task_id ON task_assignments(task_id);
CREATE INDEX idx_task_assignments_assignee ON task_assignments(assignee_type, assignee_id);
CREATE INDEX idx_task_assignments_status ON task_assignments(status);

-- Submissions
CREATE INDEX idx_submissions_user_id ON submissions(user_id);
CREATE INDEX idx_submissions_task_id ON submissions(task_id);
CREATE INDEX idx_submissions_status ON submissions(status);
CREATE INDEX idx_submissions_submitted_at ON submissions(submitted_at DESC);
CREATE INDEX idx_submissions_reviewed_by ON submissions(reviewed_by);
CREATE INDEX idx_submissions_campaign_id ON submissions(campaign_id);

CREATE INDEX idx_submission_media_submission_id ON submission_media(submission_id);

-- Gamification
CREATE INDEX idx_user_streaks_user_id ON user_streaks(user_id);
CREATE INDEX idx_user_streaks_type ON user_streaks(streak_type);
CREATE INDEX idx_streak_logs_user_date ON streak_logs(user_id, activity_date);

-- Spin Wheel & Mystery Box
CREATE INDEX idx_spin_wheel_items_wheel_id ON spin_wheel_items(spin_wheel_id);
CREATE INDEX idx_user_spins_user_id ON user_spins(user_id);
CREATE INDEX idx_user_spins_spun_at ON user_spins(spun_at DESC);
CREATE INDEX idx_mystery_box_redemptions_user_id ON mystery_box_redemptions(user_id);

-- Trivia
CREATE INDEX idx_trivia_tournaments_status ON trivia_tournaments(status);
CREATE INDEX idx_trivia_tournaments_end_date ON trivia_tournaments(end_date);
CREATE INDEX idx_trivia_participants_user_id ON trivia_participants(user_id);
CREATE INDEX idx_trivia_participants_trivia_id ON trivia_participants(trivia_id);

-- Campus Wars
CREATE INDEX idx_campus_wars_status ON campus_wars(status);
CREATE INDEX idx_campus_wars_end_date ON campus_wars(end_date);
CREATE INDEX idx_war_participants_war_id ON war_participants(war_id);
CREATE INDEX idx_war_participants_entity ON war_participants(entity_type, entity_id);

-- Secret Codes
CREATE INDEX idx_secret_codes_code ON secret_codes(code);
CREATE INDEX idx_secret_codes_valid_until ON secret_codes(valid_until);
CREATE INDEX idx_secret_code_redemptions_user_id ON secret_code_redemptions(user_id);

-- Flash Challenges
CREATE INDEX idx_flash_challenges_status ON flash_challenges(status);
CREATE INDEX idx_flash_challenges_end_time ON flash_challenges(end_time);
CREATE INDEX idx_flash_challenges_created_by ON flash_challenges(created_by);

-- Badge Bingo
CREATE INDEX idx_user_bingo_progress_user_id ON user_bingo_progress(user_id);
CREATE INDEX idx_user_bingo_progress_bingo_id ON user_bingo_progress(bingo_id);

-- Content Battles
CREATE INDEX idx_content_battles_status ON content_battles(status);
CREATE INDEX idx_content_battles_voting_end ON content_battles(voting_end);
CREATE INDEX idx_battle_submissions_battle_id ON battle_submissions(battle_id);
CREATE INDEX idx_battle_submissions_user_id ON battle_submissions(user_id);
CREATE INDEX idx_battle_votes_voter_id ON battle_votes(voter_id);
CREATE INDEX idx_battle_votes_submission_id ON battle_votes(submission_id);

-- Quests
CREATE INDEX idx_user_quests_user_id ON user_quests(user_id);
CREATE INDEX idx_user_quests_status ON user_quests(status);

-- Rewards
CREATE INDEX idx_rewards_store_is_active ON rewards_store(is_active);
CREATE INDEX idx_rewards_store_category ON rewards_store(category);
CREATE INDEX idx_user_rewards_user_id ON user_rewards(user_id);
CREATE INDEX idx_user_rewards_status ON user_rewards(status);

-- Certificates
CREATE INDEX idx_certificates_user_id ON certificates(user_id);
CREATE INDEX idx_certificates_issue_date ON certificates(issue_date DESC);

-- Leaderboards
CREATE INDEX idx_leaderboards_type ON leaderboards(leaderboard_type);
CREATE INDEX idx_leaderboards_period ON leaderboards(period_start, period_end);
CREATE INDEX idx_leaderboard_entries_leaderboard_id ON leaderboard_entries(leaderboard_id);
CREATE INDEX idx_leaderboard_entries_user_id ON leaderboard_entries(user_id);
CREATE INDEX idx_leaderboard_entries_snapshot_date ON leaderboard_entries(snapshot_date);
CREATE INDEX idx_leaderboard_entries_rank ON leaderboard_entries(rank);

-- XP Transactions
CREATE INDEX idx_xp_transactions_user_id ON xp_transactions(user_id);
CREATE INDEX idx_xp_transactions_created_at ON xp_transactions(created_at DESC);
CREATE INDEX idx_xp_transactions_type ON xp_transactions(transaction_type);

-- Activity Logs
CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX idx_activity_logs_created_at ON activity_logs(created_at DESC);
CREATE INDEX idx_activity_logs_type ON activity_logs(activity_type);

-- Referrals
CREATE INDEX idx_referrals_referrer_id ON referrals(referrer_id);
CREATE INDEX idx_referrals_referred_email ON referrals(referred_email);
CREATE INDEX idx_referrals_status ON referrals(status);

-- Notifications
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_sent_at ON notifications(sent_at DESC);
CREATE INDEX idx_notifications_type ON notifications(notification_type);

-- Surveys
CREATE INDEX idx_surveys_is_active ON surveys(is_active);
CREATE INDEX idx_surveys_end_date ON surveys(end_date);
CREATE INDEX idx_survey_responses_user_id ON survey_responses(user_id);
CREATE INDEX idx_survey_responses_survey_id ON survey_responses(survey_id);

-- Admin tables
CREATE INDEX idx_admin_actions_admin_id ON admin_actions(admin_id);
CREATE INDEX idx_admin_actions_created_at ON admin_actions(created_at DESC);
CREATE INDEX idx_scheduled_jobs_status ON scheduled_jobs(status);
CREATE INDEX idx_scheduled_jobs_scheduled_for ON scheduled_jobs(scheduled_for);