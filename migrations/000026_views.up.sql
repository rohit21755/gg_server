-- User Dashboard View
CREATE OR REPLACE VIEW user_dashboard AS
SELECT 
    u.id,
    u.email,
    u.first_name,
    u.last_name,
    u.xp,
    l.name as level_name,
    u.streak_count,
    c.name as college_name,
    s.name as state_name,
    COUNT(DISTINCT ub.badge_id) as badge_count,
    COUNT(DISTINCT ta.task_id) as active_tasks,
    COUNT(DISTINCT sub.id) as pending_submissions
FROM users u
LEFT JOIN levels l ON u.level_id = l.id
LEFT JOIN colleges c ON u.college_id = c.id
LEFT JOIN states s ON u.state_id = s.id
LEFT JOIN user_badges ub ON u.id = ub.user_id
LEFT JOIN task_assignments ta ON u.id = ta.assignee_id AND ta.assignee_type = 'user' AND ta.status = 'assigned'
LEFT JOIN tasks t ON ta.task_id = t.id AND t.is_active = true
LEFT JOIN submissions sub ON u.id = sub.user_id AND sub.status IN ('pending', 'under_review')
GROUP BY u.id, l.name, c.name, s.name;

-- Leaderboard Summary View
CREATE OR REPLACE VIEW leaderboard_summary AS
SELECT 
    le.snapshot_date,
    le.leaderboard_id,
    lb.name as leaderboard_name,
    lb.leaderboard_type,
    u.id as user_id,
    u.first_name,
    u.last_name,
    c.name as college_name,
    s.name as state_name,
    le.xp,
    le.submissions_count,
    le.referrals_count,
    le.rank,
    le.previous_rank,
    le.trend
FROM leaderboard_entries le
JOIN leaderboards lb ON le.leaderboard_id = lb.id
JOIN users u ON le.user_id = u.id
LEFT JOIN colleges c ON u.college_id = c.id
LEFT JOIN states s ON u.state_id = s.id
WHERE le.user_id IS NOT NULL;

-- Campaign Performance View
CREATE OR REPLACE VIEW campaign_performance AS
SELECT 
    c.id as campaign_id,
    c.title,
    c.campaign_type,
    c.status,
    COUNT(DISTINCT sub.user_id) as participants,
    COUNT(DISTINCT sub.id) as total_submissions,
    COUNT(DISTINCT CASE WHEN sub.status = 'approved' THEN sub.id END) as approved_submissions,
    SUM(CASE WHEN sub.status = 'approved' THEN sub.xp_awarded ELSE 0 END) as total_xp_awarded,
    AVG(CASE WHEN sub.status = 'approved' THEN sub.score END) as average_score
FROM campaigns c
LEFT JOIN tasks t ON c.id = t.campaign_id
LEFT JOIN submissions sub ON t.id = sub.task_id
GROUP BY c.id, c.title, c.campaign_type, c.status;

-- User Stats View
CREATE OR REPLACE VIEW user_stats AS
SELECT 
    u.id as user_id,
    u.first_name || ' ' || u.last_name as full_name,
    u.xp,
    l.name as current_level,
    c.name as college,
    COUNT(DISTINCT ub.badge_id) as badge_count,
    COUNT(DISTINCT sub.id) as total_submissions,
    COUNT(DISTINCT CASE WHEN sub.status = 'approved' THEN sub.id END) as approved_submissions,
    COUNT(DISTINCT r.id) as rewards_redeemed,
    u.streak_count as current_streak,
    us.longest_streak as longest_streak,
    COUNT(DISTINCT ref.id) as total_referrals
FROM users u
LEFT JOIN levels l ON u.level_id = l.id
LEFT JOIN colleges c ON u.college_id = c.id
LEFT JOIN user_badges ub ON u.id = ub.user_id
LEFT JOIN submissions sub ON u.id = sub.user_id
LEFT JOIN user_rewards r ON u.id = r.user_id AND r.status IN ('delivered', 'processing')
LEFT JOIN user_streaks us ON u.id = us.user_id AND us.streak_type = 'daily_login'
LEFT JOIN referrals ref ON u.id = ref.referrer_id
GROUP BY u.id, l.name, c.name, us.longest_streak;

-- Active Flash Challenges View
CREATE OR REPLACE VIEW active_flash_challenges AS
SELECT 
    fc.*,
    EXTRACT(EPOCH FROM (end_time - NOW())) / 3600 as hours_remaining,
    CASE 
        WHEN NOW() < start_time THEN 'upcoming'
        WHEN NOW() BETWEEN start_time AND end_time THEN 'active'
        ELSE 'ended'
    END as current_status
FROM flash_challenges fc
WHERE fc.status = 'active' OR (NOW() BETWEEN start_time AND end_time);

-- College Leaderboard View
CREATE OR REPLACE VIEW college_leaderboard AS
SELECT 
    c.id,
    c.name as college_name,
    s.name as state_name,
    COUNT(DISTINCT u.id) as total_students,
    SUM(u.xp) as total_xp,
    COUNT(DISTINCT sub.id) as total_submissions,
    COUNT(DISTINCT CASE WHEN sub.status = 'approved' THEN sub.id END) as approved_submissions,
    RANK() OVER (ORDER BY SUM(u.xp) DESC) as college_rank
FROM colleges c
LEFT JOIN states s ON c.state_id = s.id
LEFT JOIN users u ON c.id = u.college_id
LEFT JOIN submissions sub ON u.id = sub.user_id
GROUP BY c.id, c.name, s.name;