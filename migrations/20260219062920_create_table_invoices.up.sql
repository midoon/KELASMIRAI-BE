CREATE TYPE invoice_status AS ENUM (
  'pending',
  'paid',
  'expired',
  'cancelled'
);

CREATE TABLE invoices
(
  id                UUID PRIMARY KEY,
  tenant_id         UUID NOT NULL,
  subscription_id   UUID NOT NULL,
  code              TEXT NOT NULL,
  amount            NUMERIC(15,2) NOT NULL CHECK (amount >= 0),
  status            invoice_status NOT NULL DEFAULT 'pending',
  due_date          TIMESTAMPTZ NOT NULL,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_invoice_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE CASCADE,

  CONSTRAINT fk_invoice_subscription
    FOREIGN KEY (subscription_id)
    REFERENCES tenant_subscriptions(id)
    ON DELETE CASCADE,

  CONSTRAINT unique_invoice_code UNIQUE (code)
);

CREATE INDEX idx_invoices_tenant_id
ON invoices (tenant_id);

CREATE INDEX idx_invoices_status
ON invoices (status);

CREATE INDEX idx_invoices_due_date
ON invoices (due_date);
