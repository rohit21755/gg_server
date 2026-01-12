-- Insert default rewards
-- Using WHERE NOT EXISTS to avoid duplicates if rewards already exist
INSERT INTO rewards_store (name, description, reward_type, category, image_url, xp_cost, coin_cost, cash_cost, quantity_available, quantity_sold, is_featured, is_active, validity_days)
SELECT 'Amazon Gift Card ₹500', 'Redeemable Amazon gift card worth ₹500', 'gift_card', 'digital', 'https://example.com/rewards/amazon-500.jpg', 5000, 0, NULL, 100, 0, true, true, 365
WHERE NOT EXISTS (SELECT 1 FROM rewards_store WHERE name = 'Amazon Gift Card ₹500');

INSERT INTO rewards_store (name, description, reward_type, category, image_url, xp_cost, coin_cost, cash_cost, quantity_available, quantity_sold, is_featured, is_active, validity_days)
SELECT 'Branded T-Shirt', 'Premium quality branded t-shirt', 'physical', 'merchandise', 'https://example.com/rewards/tshirt.jpg', 3000, 500, NULL, 50, 0, true, true, NULL
WHERE NOT EXISTS (SELECT 1 FROM rewards_store WHERE name = 'Branded T-Shirt');

INSERT INTO rewards_store (name, description, reward_type, category, image_url, xp_cost, coin_cost, cash_cost, quantity_available, quantity_sold, is_featured, is_active, validity_days)
SELECT 'XP Boost 2x (7 days)', 'Double your XP earnings for 7 days', 'xp_boost', 'digital', 'https://example.com/rewards/xp-boost.jpg', 2000, 0, NULL, NULL, 0, false, true, 7
WHERE NOT EXISTS (SELECT 1 FROM rewards_store WHERE name = 'XP Boost 2x (7 days)');

INSERT INTO rewards_store (name, description, reward_type, category, image_url, xp_cost, coin_cost, cash_cost, quantity_available, quantity_sold, is_featured, is_active, validity_days)
SELECT 'Premium Profile Skin', 'Exclusive profile skin design', 'profile_skin', 'digital', 'https://example.com/rewards/skin-premium.jpg', 1500, 200, NULL, NULL, 0, false, true, NULL
WHERE NOT EXISTS (SELECT 1 FROM rewards_store WHERE name = 'Premium Profile Skin');

INSERT INTO rewards_store (name, description, reward_type, category, image_url, xp_cost, coin_cost, cash_cost, quantity_available, quantity_sold, is_featured, is_active, validity_days)
SELECT 'Certificate of Excellence', 'Digital certificate for outstanding performance', 'certificate', 'digital', 'https://example.com/rewards/certificate.jpg', 1000, 0, NULL, NULL, 0, false, true, NULL
WHERE NOT EXISTS (SELECT 1 FROM rewards_store WHERE name = 'Certificate of Excellence');
