-- Remove default badges by name
DELETE FROM badges WHERE name IN (
  'First Steps',
  'Referral Master',
  'Campaign Champion',
  'Streak Warrior',
  'Hidden Gem',
  'Early Bird'
);
