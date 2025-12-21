-- XP Transactions Log
CREATE TABLE xp_transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    transaction_type VARCHAR(50) NOT NULL CHECK (transaction_type IN ('task_completion', 'referral', 'streak', 'spin_wheel', 'mystery_box', 'quiz', 'battle_win', 'redemption', 'correction', 'bonus')),
    amount INTEGER NOT NULL,
    balance_after INTEGER NOT NULL,
    source_id INTEGER,
    source_type VARCHAR(50),
    description TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign key
ALTER TABLE xp_transactions 
ADD CONSTRAINT fk_xp_transactions_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;