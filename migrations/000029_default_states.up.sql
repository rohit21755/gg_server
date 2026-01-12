-- Insert default Indian states
-- Using ON CONFLICT DO NOTHING to avoid duplicates if states already exist
INSERT INTO states (name, code) VALUES
('Maharashtra', 'MH'),
('Karnataka', 'KA'),
('Tamil Nadu', 'TN'),
('Delhi', 'DL'),
('Gujarat', 'GJ'),
('Rajasthan', 'RJ'),
('West Bengal', 'WB'),
('Uttar Pradesh', 'UP'),
('Telangana', 'TG'),
('Kerala', 'KL'),
('Punjab', 'PB'),
('Haryana', 'HR'),
('Madhya Pradesh', 'MP'),
('Bihar', 'BR'),
('Andhra Pradesh', 'AP')
ON CONFLICT (code) DO NOTHING;
