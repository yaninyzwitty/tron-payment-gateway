-- Enable UUID functions if not already enabled (CockroachDB usually has them available)
-- CREATE EXTENSION IF NOT EXISTS pgcrypto;  -- (not needed in CockroachDB >= v21)

-- clients table
CREATE TABLE clients (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name STRING NOT NULL,
  api_key STRING UNIQUE NOT NULL,
  is_active BOOL DEFAULT TRUE,
  created_at TIMESTAMPTZ DEFAULT now()
);

