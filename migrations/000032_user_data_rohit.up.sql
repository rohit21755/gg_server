-- Seed comprehensive data for user ID 3 (Rohit Kumar, r@exmail.com)
-- This migration assumes:
-- - User with ID 3 exists
-- - Badges from migration 000030 exist (IDs 1-6)
-- - At least one campaign exists (we'll create one if needed)
-- - At least one spin wheel exists (we'll use the first one)

-- Step 1: Award badges to user (assuming badges 1-6 exist from migration 000030)
-- First Steps, Referral Master, Campaign Champion, Streak Warrior, Hidden Gem, Early Bird
INSERT INTO user_badges (user_id, badge_id, earned_at)
SELECT 3, 1, NOW() - INTERVAL '30 days'
WHERE NOT EXISTS (SELECT 1 FROM user_badges WHERE user_id = 3 AND badge_id = 1);

INSERT INTO user_badges (user_id, badge_id, earned_at)
SELECT 3, 2, NOW() - INTERVAL '20 days'
WHERE NOT EXISTS (SELECT 1 FROM user_badges WHERE user_id = 3 AND badge_id = 2);

INSERT INTO user_badges (user_id, badge_id, earned_at)
SELECT 3, 3, NOW() - INTERVAL '15 days'
WHERE NOT EXISTS (SELECT 1 FROM user_badges WHERE user_id = 3 AND badge_id = 3);

INSERT INTO user_badges (user_id, badge_id, earned_at)
SELECT 3, 4, NOW() - INTERVAL '10 days'
WHERE NOT EXISTS (SELECT 1 FROM user_badges WHERE user_id = 3 AND badge_id = 4);

-- Step 2: Create a test campaign if none exists and get campaign ID
DO $$
DECLARE
    campaign_id_var INTEGER;
    task1_id INTEGER;
    task2_id INTEGER;
    task3_id INTEGER;
    task4_id INTEGER;
    task5_id INTEGER;
    spin_wheel_id_var INTEGER;
    spin_item_id_var INTEGER;
BEGIN
    -- Get or create campaign
    SELECT id INTO campaign_id_var FROM campaigns WHERE title = 'Summer Engagement Campaign 2024' LIMIT 1;
    
    IF campaign_id_var IS NULL THEN
        INSERT INTO campaigns (uuid, title, description, campaign_type, category, start_date, end_date, status, priority, is_gg_led, created_at)
        VALUES (
            gen_random_uuid()::text,
            'Summer Engagement Campaign 2024',
            'Comprehensive summer campaign for campus ambassadors',
            'thematic',
            'online',
            NOW() - INTERVAL '30 days',
            NOW() + INTERVAL '30 days',
            'active',
            'high',
            true,
            NOW() - INTERVAL '30 days'
        )
        RETURNING id INTO campaign_id_var;
    END IF;

    -- Step 3: Create tasks (mix of completed and ongoing)
    -- Task 1: Completed task
    INSERT INTO tasks (uuid, campaign_id, title, description, task_type, proof_type, xp_reward, coin_reward, duration_hours, priority, assignment_type, is_active, created_at)
    SELECT gen_random_uuid()::text, campaign_id_var, 'Create Instagram Post', 'Create an engaging Instagram post about our brand', 'online', 'screenshot', 500, 100, 72, 'medium', 'individual', true, NOW() - INTERVAL '25 days'
    WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Create Instagram Post' AND campaign_id = campaign_id_var)
    RETURNING id INTO task1_id;

    -- Get task1_id if it already exists
    IF task1_id IS NULL THEN
        SELECT id INTO task1_id FROM tasks WHERE title = 'Create Instagram Post' AND campaign_id = campaign_id_var LIMIT 1;
    END IF;

    -- Task 2: Completed task
    INSERT INTO tasks (uuid, campaign_id, title, description, task_type, proof_type, xp_reward, coin_reward, duration_hours, priority, assignment_type, is_active, created_at)
    SELECT gen_random_uuid()::text, campaign_id_var, 'Share on LinkedIn', 'Share campaign announcement on LinkedIn', 'online', 'url', 300, 50, 48, 'low', 'individual', true, NOW() - INTERVAL '20 days'
    WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Share on LinkedIn' AND campaign_id = campaign_id_var)
    RETURNING id INTO task2_id;

    IF task2_id IS NULL THEN
        SELECT id INTO task2_id FROM tasks WHERE title = 'Share on LinkedIn' AND campaign_id = campaign_id_var LIMIT 1;
    END IF;

    -- Task 3: Ongoing task (pending submission)
    INSERT INTO tasks (uuid, campaign_id, title, description, task_type, proof_type, xp_reward, coin_reward, duration_hours, priority, assignment_type, is_active, created_at)
    SELECT gen_random_uuid()::text, campaign_id_var, 'Create YouTube Video', 'Create a 2-5 minute YouTube video review', 'online', 'video', 1500, 300, 168, 'high', 'individual', true, NOW() - INTERVAL '5 days'
    WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Create YouTube Video' AND campaign_id = campaign_id_var)
    RETURNING id INTO task3_id;

    IF task3_id IS NULL THEN
        SELECT id INTO task3_id FROM tasks WHERE title = 'Create YouTube Video' AND campaign_id = campaign_id_var LIMIT 1;
    END IF;

    -- Task 4: Ongoing task (under review)
    INSERT INTO tasks (uuid, campaign_id, title, description, task_type, proof_type, xp_reward, coin_reward, duration_hours, priority, assignment_type, is_active, created_at)
    SELECT gen_random_uuid()::text, campaign_id_var, 'Campus Event Organization', 'Organize a campus event with at least 50 participants', 'offline', 'pdf', 2000, 500, NULL, 'high', 'college', true, NOW() - INTERVAL '10 days'
    WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Campus Event Organization' AND campaign_id = campaign_id_var)
    RETURNING id INTO task4_id;

    IF task4_id IS NULL THEN
        SELECT id INTO task4_id FROM tasks WHERE title = 'Campus Event Organization' AND campaign_id = campaign_id_var LIMIT 1;
    END IF;

    -- Task 5: Available task (not started)
    INSERT INTO tasks (uuid, campaign_id, title, description, task_type, proof_type, xp_reward, coin_reward, duration_hours, priority, assignment_type, is_active, created_at)
    SELECT gen_random_uuid()::text, campaign_id_var, 'Social Media Story Series', 'Create a 7-day story series on Instagram', 'online', 'screenshot', 800, 150, 168, 'medium', 'individual', true, NOW() - INTERVAL '3 days'
    WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Social Media Story Series' AND campaign_id = campaign_id_var)
    RETURNING id INTO task5_id;

    IF task5_id IS NULL THEN
        SELECT id INTO task5_id FROM tasks WHERE title = 'Social Media Story Series' AND campaign_id = campaign_id_var LIMIT 1;
    END IF;

    -- Step 4: Create submissions for completed tasks (only if they don't exist)
    -- Submission 1: Approved (completed)
    INSERT INTO submissions (uuid, task_id, user_id, campaign_id, proof_type, proof_url, proof_text, status, submitted_at, reviewed_at, xp_awarded, coins_awarded, is_winner)
    SELECT 
        gen_random_uuid()::text,
        task1_id,
        3,
        campaign_id_var,
        'screenshot',
        'https://example.com/submissions/instagram-post-1.jpg',
        'Posted on Instagram with hashtags #campusambassador #brand',
        'approved',
        NOW() - INTERVAL '24 days',
        NOW() - INTERVAL '23 days',
        500,
        100,
        false
    WHERE NOT EXISTS (SELECT 1 FROM submissions WHERE task_id = task1_id AND user_id = 3);

    -- Submission 2: Approved (completed)
    INSERT INTO submissions (uuid, task_id, user_id, campaign_id, proof_type, proof_url, proof_text, status, submitted_at, reviewed_at, xp_awarded, coins_awarded, is_winner)
    SELECT 
        gen_random_uuid()::text,
        task2_id,
        3,
        campaign_id_var,
        'url',
        'https://www.linkedin.com/posts/rohit-kumar_campaign-activity-1234567890',
        'Shared campaign announcement on LinkedIn',
        'approved',
        NOW() - INTERVAL '19 days',
        NOW() - INTERVAL '18 days',
        300,
        50,
        false
    WHERE NOT EXISTS (SELECT 1 FROM submissions WHERE task_id = task2_id AND user_id = 3);

    -- Submission 3: Pending (ongoing)
    INSERT INTO submissions (uuid, task_id, user_id, campaign_id, proof_type, proof_url, proof_text, status, submitted_at)
    SELECT 
        gen_random_uuid()::text,
        task3_id,
        3,
        campaign_id_var,
        'video',
        'https://www.youtube.com/watch?v=example123',
        'Created YouTube video review of the product',
        'pending',
        NOW() - INTERVAL '2 days'
    WHERE NOT EXISTS (SELECT 1 FROM submissions WHERE task_id = task3_id AND user_id = 3);

    -- Submission 4: Under Review (ongoing)
    INSERT INTO submissions (uuid, task_id, user_id, campaign_id, proof_type, proof_url, proof_text, status, submitted_at)
    SELECT 
        gen_random_uuid()::text,
        task4_id,
        3,
        campaign_id_var,
        'pdf',
        'https://example.com/submissions/event-report.pdf',
        'Organized campus event with 75 participants',
        'under_review',
        NOW() - INTERVAL '1 day'
    WHERE NOT EXISTS (SELECT 1 FROM submissions WHERE task_id = task4_id AND user_id = 3);

    -- Step 5: Create task assignments for ongoing/available tasks
    INSERT INTO task_assignments (task_id, assignee_type, assignee_id, status, assigned_at)
    SELECT task3_id, 'user', 3, 'accepted', NOW() - INTERVAL '5 days'
    WHERE NOT EXISTS (SELECT 1 FROM task_assignments WHERE task_id = task3_id AND assignee_id = 3 AND assignee_type = 'user');

    INSERT INTO task_assignments (task_id, assignee_type, assignee_id, status, assigned_at)
    SELECT task4_id, 'user', 3, 'accepted', NOW() - INTERVAL '10 days'
    WHERE NOT EXISTS (SELECT 1 FROM task_assignments WHERE task_id = task4_id AND assignee_id = 3 AND assignee_type = 'user');

    INSERT INTO task_assignments (task_id, assignee_type, assignee_id, status, assigned_at)
    SELECT task5_id, 'user', 3, 'assigned', NOW() - INTERVAL '3 days'
    WHERE NOT EXISTS (SELECT 1 FROM task_assignments WHERE task_id = task5_id AND assignee_id = 3 AND assignee_type = 'user');

    -- Step 6: Create spin wheel data
    -- Get first active spin wheel
    SELECT id INTO spin_wheel_id_var FROM spin_wheels WHERE is_active = true LIMIT 1;
    
    IF spin_wheel_id_var IS NOT NULL THEN
        -- Get a spin wheel item (preferably XP reward)
        SELECT id INTO spin_item_id_var FROM spin_wheel_items 
        WHERE spin_wheel_id = spin_wheel_id_var AND item_type = 'xp' AND is_active = true 
        LIMIT 1;
        
        IF spin_item_id_var IS NOT NULL THEN
            -- Create multiple spin records (user has spun the wheel multiple times)
            INSERT INTO user_spins (user_id, spin_wheel_id, spin_wheel_item_id, earned_value, spun_at)
            SELECT 3, spin_wheel_id_var, spin_item_id_var, 100, NOW() - INTERVAL '7 days'
            WHERE NOT EXISTS (SELECT 1 FROM user_spins WHERE user_id = 3 AND spun_at::date = (NOW() - INTERVAL '7 days')::date);

            INSERT INTO user_spins (user_id, spin_wheel_id, spin_wheel_item_id, earned_value, spun_at)
            SELECT 3, spin_wheel_id_var, spin_item_id_var, 250, NOW() - INTERVAL '5 days'
            WHERE NOT EXISTS (SELECT 1 FROM user_spins WHERE user_id = 3 AND spun_at::date = (NOW() - INTERVAL '5 days')::date);

            INSERT INTO user_spins (user_id, spin_wheel_id, spin_wheel_item_id, earned_value, spun_at)
            SELECT 3, spin_wheel_id_var, spin_item_id_var, 500, NOW() - INTERVAL '3 days'
            WHERE NOT EXISTS (SELECT 1 FROM user_spins WHERE user_id = 3 AND spun_at::date = (NOW() - INTERVAL '3 days')::date);
        END IF;
    END IF;

    -- Step 7: Create XP transactions for completed tasks
    INSERT INTO xp_transactions (user_id, transaction_type, amount, balance_after, source_id, source_type, description)
    SELECT 3, 'task_completion', 500, 600, task1_id, 'task', 'Completed: Create Instagram Post'
    WHERE NOT EXISTS (SELECT 1 FROM xp_transactions WHERE user_id = 3 AND source_id = task1_id AND source_type = 'task');

    INSERT INTO xp_transactions (user_id, transaction_type, amount, balance_after, source_id, source_type, description)
    SELECT 3, 'task_completion', 300, 900, task2_id, 'task', 'Completed: Share on LinkedIn'
    WHERE NOT EXISTS (SELECT 1 FROM xp_transactions WHERE user_id = 3 AND source_id = task2_id AND source_type = 'task');

END $$;
