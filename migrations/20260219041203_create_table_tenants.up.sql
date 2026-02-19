CREATE TYPE tenant_status AS ENUM ('active', 'suspended', 'cancelled');

CREATE TABLE tenants
(
  id            UUID PRIMARY KEY,
  name          TEXT              NOT NULL,
  slug          TEXT              NOT NULL UNIQUE,
  email         TEXT              NOT NULL UNIQUE,
  phone         TEXT,
  address       TEXT,
  logo_url      TEXT,
  status        tenant_status     NOT NULL DEFAULT 'active',
  trial_ends_at TIMESTAMPTZ       NOT NULL,
  created_at    TIMESTAMPTZ       NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ       NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tenants_trial_ends_at
ON tenants (trial_ends_at);
