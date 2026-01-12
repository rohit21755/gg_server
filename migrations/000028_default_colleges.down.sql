-- Remove default colleges by code
DELETE FROM colleges WHERE code IN (
  'HARVARD',
  'STANFORD',
  'MIT',
  'UCB',
  'YALE',
  'PRINCETON',
  'COLUMBIA',
  'UCHICAGO',
  'UPENN',
  'CALTECH',
  'DUKE',
  'NORTHWESTERN',
  'JHU',
  'DARTMOUTH',
  'BROWN'
);
