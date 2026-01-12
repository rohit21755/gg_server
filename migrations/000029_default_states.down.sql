-- Remove default states by code
DELETE FROM states WHERE code IN (
  'MH', 'KA', 'TN', 'DL', 'GJ', 'RJ', 'WB', 'UP', 'TG', 'KL', 'PB', 'HR', 'MP', 'BR', 'AP'
);
