-- Certificates
CREATE TABLE certificates (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    certificate_type VARCHAR(50) CHECK (certificate_type IN ('achievement', 'completion', 'winner', 'participation')),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    issuing_authority VARCHAR(200) DEFAULT 'Grove Growth',
    issue_date DATE NOT NULL,
    expiry_date DATE,
    certificate_url TEXT NOT NULL,
    template_id INTEGER,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign key
ALTER TABLE certificates 
ADD CONSTRAINT fk_certificates_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;