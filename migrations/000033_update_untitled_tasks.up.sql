-- Update tasks with untitled or empty names and descriptions
-- This migration updates tasks that have "untitled" (case-insensitive), empty, or NULL titles/descriptions

DO $$
DECLARE
    task_record RECORD;
    new_title VARCHAR(200);
    new_description TEXT;
    campaign_title VARCHAR(200);
BEGIN
    -- Loop through all tasks that need updating
    FOR task_record IN 
        SELECT 
            t.id,
            t.title,
            t.description,
            t.task_type,
            t.proof_type,
            t.campaign_id,
            t.priority,
            t.xp_reward,
            t.coin_reward,
            c.title as campaign_title
        FROM tasks t
        LEFT JOIN campaigns c ON t.campaign_id = c.id
        WHERE 
            -- Check for untitled or empty titles
            (t.title IS NULL OR TRIM(t.title) = '' OR LOWER(TRIM(t.title)) = 'untitled' OR LOWER(TRIM(t.title)) LIKE '%untitled%')
            OR
            -- Check for empty or untitled descriptions
            (t.description IS NULL OR TRIM(t.description) = '' OR LOWER(TRIM(t.description)) = 'untitled' OR LOWER(TRIM(t.description)) LIKE '%untitled%')
        ORDER BY t.id
    LOOP
        -- Get campaign title if available
        campaign_title := COALESCE(task_record.campaign_title, 'Campaign');
        
        -- Generate a meaningful title based on task properties
        new_title := CASE
            WHEN task_record.task_type = 'online' AND task_record.proof_type = 'screenshot' THEN
                'Social Media Engagement Task'
            WHEN task_record.task_type = 'online' AND task_record.proof_type = 'url' THEN
                'Online Content Sharing Task'
            WHEN task_record.task_type = 'online' AND task_record.proof_type = 'video' THEN
                'Video Content Creation Task'
            WHEN task_record.task_type = 'online' AND task_record.proof_type = 'pdf' THEN
                'Document Submission Task'
            WHEN task_record.task_type = 'online' AND task_record.proof_type = 'text' THEN
                'Text-Based Submission Task'
            WHEN task_record.task_type = 'offline' AND task_record.proof_type = 'screenshot' THEN
                'Offline Event Participation Task'
            WHEN task_record.task_type = 'offline' AND task_record.proof_type = 'pdf' THEN
                'Offline Event Report Task'
            WHEN task_record.task_type = 'offline' AND task_record.proof_type = 'video' THEN
                'Offline Event Video Task'
            WHEN task_record.task_type = 'solo' THEN
                'Individual Task'
            WHEN task_record.task_type = 'group' THEN
                'Group Collaboration Task'
            ELSE
                'Campaign Task'
        END;
        
        -- Add priority info to title if available and not medium
        IF task_record.priority IS NOT NULL AND task_record.priority != 'medium' THEN
            new_title := new_title || ' (' || INITCAP(task_record.priority) || ' Priority)';
        END IF;
        
        -- Generate a meaningful description
        new_description := CASE
            WHEN task_record.task_type = 'online' AND task_record.proof_type = 'screenshot' THEN
                'Complete this online engagement task by sharing content on social media platforms. Submit a screenshot as proof of completion.'
            WHEN task_record.task_type = 'online' AND task_record.proof_type = 'url' THEN
                'Share content online and provide the URL as proof of your submission.'
            WHEN task_record.task_type = 'online' AND task_record.proof_type = 'video' THEN
                'Create and share a video content piece. Submit the video link as proof of completion.'
            WHEN task_record.task_type = 'online' AND task_record.proof_type = 'pdf' THEN
                'Complete the task and submit a PDF document as proof of your work.'
            WHEN task_record.task_type = 'online' AND task_record.proof_type = 'text' THEN
                'Complete the task and submit your response as text. Follow the submission guidelines carefully.'
            WHEN task_record.task_type = 'offline' AND task_record.proof_type = 'screenshot' THEN
                'Participate in an offline event and submit a screenshot as proof of attendance.'
            WHEN task_record.task_type = 'offline' AND task_record.proof_type = 'pdf' THEN
                'Organize or participate in an offline event and submit a detailed report in PDF format.'
            WHEN task_record.task_type = 'offline' AND task_record.proof_type = 'video' THEN
                'Record and submit a video of the offline event as proof of participation.'
            WHEN task_record.task_type = 'solo' THEN
                'Complete this individual task on your own. Follow the submission guidelines to earn rewards.'
            WHEN task_record.task_type = 'group' THEN
                'Collaborate with your team to complete this group task. Ensure all team members contribute.'
            ELSE
                'Complete this campaign task according to the provided instructions. Submit proof of completion to earn rewards.'
        END;
        
        -- Add reward information to description
        IF task_record.xp_reward > 0 OR task_record.coin_reward > 0 THEN
            new_description := new_description || E'\n\nRewards: ';
            IF task_record.xp_reward > 0 THEN
                new_description := new_description || task_record.xp_reward || ' XP';
            END IF;
            IF task_record.xp_reward > 0 AND task_record.coin_reward > 0 THEN
                new_description := new_description || ' and ';
            END IF;
            IF task_record.coin_reward > 0 THEN
                new_description := new_description || task_record.coin_reward || ' Coins';
            END IF;
        END IF;
        
        -- Update the task
        UPDATE tasks
        SET 
            title = CASE 
                WHEN title IS NULL OR TRIM(title) = '' OR LOWER(TRIM(title)) = 'untitled' OR LOWER(TRIM(title)) LIKE '%untitled%' 
                THEN new_title 
                ELSE title 
            END,
            description = CASE 
                WHEN description IS NULL OR TRIM(description) = '' OR LOWER(TRIM(description)) = 'untitled' OR LOWER(TRIM(description)) LIKE '%untitled%' 
                THEN new_description 
                ELSE description 
            END,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = task_record.id;
        
        RAISE NOTICE 'Updated task ID % with title: %', task_record.id, new_title;
    END LOOP;
    
    RAISE NOTICE 'Migration completed: Updated all untitled tasks';
END $$;
