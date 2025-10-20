-- Enable UUID functions if not already enabled (CockroachDB usually has them available)
-- CREATE EXTENSION IF NOT EXISTS pgcrypto;  -- (not needed in CockroachDB >= v21)

-- Payments Table (Main transactions)
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    amount DECIMAL(18,6) NOT NULL,
    unique_wallet STRING NOT NULL,
    status STRING NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'CONFIRMED', 'EXPIRED')),
    expires_at TIMESTAMPTZ NOT NULL,
    confirmed_at TIMESTAMPTZ,
    attempt_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now()
);



