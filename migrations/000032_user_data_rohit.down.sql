-- Remove all data created for user ID 3 (Rohit Kumar)

-- Remove XP transactions
DELETE FROM xp_transactions WHERE user_id = 3;

-- Remove user spins
DELETE FROM user_spins WHERE user_id = 3;

-- Remove task assignments
DELETE FROM task_assignments WHERE assignee_id = 3 AND assignee_type = 'user';

-- Remove submissions
DELETE FROM submissions WHERE user_id = 3;

-- Remove tasks created for this user (if they were only for this user)
-- Note: This is a conservative approach - you may want to keep tasks if they're used by others
DELETE FROM tasks WHERE id IN (
    SELECT t.id FROM tasks t
    LEFT JOIN submissions s ON t.id = s.task_id
    WHERE s.user_id = 3 OR s.user_id IS NULL
) AND title IN (
    'Create Instagram Post',
    'Share on LinkedIn',
    'Create YouTube Video',
    'Campus Event Organization',
    'Social Media Story Series'
);

-- Remove campaign if it was only for this user
DELETE FROM campaigns WHERE title = 'Summer Engagement Campaign 2024' 
AND id NOT IN (SELECT DISTINCT campaign_id FROM submissions WHERE user_id != 3);

-- Remove user badges
DELETE FROM user_badges WHERE user_id = 3;
