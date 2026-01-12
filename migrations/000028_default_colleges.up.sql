-- Insert default colleges
-- Using ON CONFLICT DO NOTHING to avoid duplicates if colleges already exist
-- Note: This will skip insertion if a college with the same code already exists
INSERT INTO colleges (name, code, is_active) VALUES
('Harvard University', 'HARVARD', true),
('Stanford University', 'STANFORD', true),
('Massachusetts Institute of Technology', 'MIT', true),
('University of California, Berkeley', 'UCB', true),
('Yale University', 'YALE', true),
('Princeton University', 'PRINCETON', true),
('Columbia University', 'COLUMBIA', true),
('University of Chicago', 'UCHICAGO', true),
('University of Pennsylvania', 'UPENN', true),
('California Institute of Technology', 'CALTECH', true),
('Duke University', 'DUKE', true),
('Northwestern University', 'NORTHWESTERN', true),
('Johns Hopkins University', 'JHU', true),
('Dartmouth College', 'DARTMOUTH', true),
('Brown University', 'BROWN', true)
ON CONFLICT (code) DO NOTHING;
