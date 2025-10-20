-- Enable UUID functions if not already enabled (CockroachDB usually has them available)
-- CREATE EXTENSION IF NOT EXISTS pgcrypto;  -- (not needed in CockroachDB >= v21)

-- Attempts Table (Address regeneration attempts)
CREATE TABLE payment_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    attempt_number INT NOT NULL CHECK (attempt_number > 0),
    generated_wallet STRING NOT NULL,
    generated_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(payment_id, attempt_number)
);


CREATE INDEX idx_payment_attempts_payment_id ON payment_attempts(payment_id);