-- Enable UUID functions if not already enabled (CockroachDB usually has them available)
-- CREATE EXTENSION IF NOT EXISTS pgcrypto;  -- (not needed in CockroachDB >= v21)



-- accounts table
CREATE TABLE accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  client_id UUID NOT NULL,
  name STRING NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  CONSTRAINT fk_accounts_client FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
);
