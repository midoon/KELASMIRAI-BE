CREATE TYPE payment_status AS ENUM (
  'pending',
  'settlement',
  'expire',
  'cancel',
  'deny'
);

CREATE TABLE payments
(
  id                        UUID PRIMARY KEY,
  invoice_id                UUID NOT NULL,
  midtrans_order_id         TEXT NOT NULL,
  midtrans_transaction_id   TEXT,
  amount                    NUMERIC(15,2) NOT NULL CHECK (amount >= 0),
  status                    payment_status NOT NULL DEFAULT 'pending',
  paid_at                   TIMESTAMPTZ,
  raw_response_json         JSONB,
  created_at                TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_payment_invoice
    FOREIGN KEY (invoice_id)
    REFERENCES invoices(id)
    ON DELETE RESTRICT,

  CONSTRAINT unique_midtrans_order_id UNIQUE (midtrans_order_id)
);

CREATE INDEX idx_payments_invoice_id
ON payments (invoice_id);

CREATE INDEX idx_payments_status
ON payments (status);

CREATE INDEX idx_payments_midtrans_transaction_id
ON payments (midtrans_transaction_id);
