-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to update user level based on XP
CREATE OR REPLACE FUNCTION update_user_level()
RETURNS TRIGGER AS $$
DECLARE
    new_level_id INTEGER;
BEGIN
    -- Find the appropriate level based on XP
    SELECT id INTO new_level_id
    FROM levels
    WHERE NEW.xp >= min_xp AND (NEW.xp <= max_xp OR max_xp IS NULL)
    ORDER BY rank_order DESC
    LIMIT 1;
    
    IF new_level_id IS NOT NULL AND NEW.level_id != new_level_id THEN
        NEW.level_id = new_level_id;
        
        -- Log level up event
        INSERT INTO activity_logs (user_id, activity_type, activity_data)
        VALUES (NEW.id, 'level_up', jsonb_build_object(
            'old_level', OLD.level_id,
            'new_level', new_level_id,
            'xp', NEW.xp
        ));
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to process XP transactions
CREATE OR REPLACE FUNCTION process_xp_transaction()
RETURNS TRIGGER AS $$
BEGIN
    -- Update user's XP
    UPDATE users 
    SET xp = xp + NEW.amount
    WHERE id = NEW.user_id;
    
    -- Set balance_after
    NEW.balance_after = (SELECT xp FROM users WHERE id = NEW.user_id);
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to update college stats when user XP changes
CREATE OR REPLACE FUNCTION update_college_stats()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.college_id IS DISTINCT FROM NEW.college_id THEN
        -- Remove XP from old college
        UPDATE colleges 
        SET total_xp = total_xp - OLD.xp,
            total_cas = total_cas - 1
        WHERE id = OLD.college_id;
        
        -- Add XP to new college
        UPDATE colleges 
        SET total_xp = total_xp + NEW.xp,
            total_cas = total_cas + 1
        WHERE id = NEW.college_id;
    ELSIF OLD.xp != NEW.xp THEN
        -- Update XP difference for same college
        UPDATE colleges 
        SET total_xp = total_xp + (NEW.xp - OLD.xp)
        WHERE id = NEW.college_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to update campaign participant count
CREATE OR REPLACE FUNCTION update_campaign_participants()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'active' AND OLD.status != 'active' THEN
        UPDATE campaigns 
        SET current_participants = current_participants + 1
        WHERE id = NEW.campaign_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to log activity on user actions
CREATE OR REPLACE FUNCTION log_user_activity()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_TABLE_NAME = 'submissions' THEN
        INSERT INTO activity_logs (user_id, activity_type, activity_data)
        VALUES (NEW.user_id, 'submission_' || NEW.status, jsonb_build_object(
            'submission_id', NEW.id,
            'task_id', NEW.task_id,
            'campaign_id', NEW.campaign_id,
            'xp_awarded', NEW.xp_awarded
        ));
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply triggers to tables
CREATE TRIGGER update_users_updated_at 
BEFORE UPDATE ON users 
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_level_trigger 
BEFORE UPDATE OF xp ON users 
FOR EACH ROW 
EXECUTE FUNCTION update_user_level();

CREATE TRIGGER update_college_stats_trigger 
AFTER UPDATE OF xp, college_id ON users 
FOR EACH ROW 
EXECUTE FUNCTION update_college_stats();

CREATE TRIGGER process_xp_transaction_trigger 
BEFORE INSERT ON xp_transactions 
FOR EACH ROW 
EXECUTE FUNCTION process_xp_transaction();

CREATE TRIGGER update_campaigns_updated_at 
BEFORE UPDATE ON campaigns 
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tasks_updated_at 
BEFORE UPDATE ON tasks 
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_submissions_updated_at 
BEFORE UPDATE ON submissions 
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER log_submission_activity 
AFTER INSERT OR UPDATE ON submissions 
FOR EACH ROW 
EXECUTE FUNCTION log_user_activity();

CREATE TRIGGER update_rewards_store_updated_at 
BEFORE UPDATE ON rewards_store 
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_system_config_updated_at 
BEFORE UPDATE ON system_config 
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at_column();