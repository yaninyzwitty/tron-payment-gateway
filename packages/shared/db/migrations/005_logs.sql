-- Enable UUID functions if not already enabled (CockroachDB usually has them available)
-- CREATE EXTENSION IF NOT EXISTS pgcrypto;  -- (not needed in CockroachDB >= v21)

-- Logs Table (Blockchain events, webhook status, errors, etc.)
CREATE TABLE logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID REFERENCES payments(id) ON DELETE CASCADE,
    event_type STRING NOT NULL, -- e.g., 'ADDRESS_GENERATED', 'TX_CONFIRMED', 'WEBHOOK_SENT'
    message STRING,
    raw_data JSONB,
    created_at TIMESTAMPTZ DEFAULT now()
);