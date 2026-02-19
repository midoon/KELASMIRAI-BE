CREATE TYPE subscription_status AS ENUM (
  'trial',
  'active',
  'past_due',
  'cancelled'
);

CREATE TABLE tenant_subscriptions
(
  id                        UUID PRIMARY KEY,
  tenant_id                 UUID NOT NULL,
  plan_id                   UUID NOT NULL,
  billing_cycle             TEXT NOT NULL CHECK (billing_cycle IN ('monthly', 'yearly')),
  price                     NUMERIC(15,2) NOT NULL CHECK (price >= 0),
  status                    subscription_status NOT NULL DEFAULT 'trial',
  started_at                TIMESTAMPTZ NOT NULL,
  ended_at                  TIMESTAMPTZ,
  next_billing_at           TIMESTAMPTZ NOT NULL,
  midtrans_subscription_id  TEXT,
  created_at                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE CASCADE,

  CONSTRAINT fk_plan
    FOREIGN KEY (plan_id)
    REFERENCES subscription_plans(id)
);

CREATE INDEX idx_tenant_subscriptions_tenant_id
ON tenant_subscriptions (tenant_id);

CREATE INDEX idx_tenant_subscriptions_status
ON tenant_subscriptions (status);

CREATE INDEX idx_tenant_subscriptions_next_billing
ON tenant_subscriptions (next_billing_at);
