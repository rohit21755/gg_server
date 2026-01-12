-- Remove default rewards by name
DELETE FROM rewards_store WHERE name IN (
  'Amazon Gift Card ₹500',
  'Branded T-Shirt',
  'XP Boost 2x (7 days)',
  'Premium Profile Skin',
  'Certificate of Excellence'
);
