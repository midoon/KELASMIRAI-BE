CREATE TABLE webhook_logs
(
  id               UUID PRIMARY KEY,
  provider         TEXT NOT NULL,
  external_id      TEXT, -- midtrans_order_id
  payload_json     JSONB NOT NULL,
  signature_valid  BOOLEAN NOT NULL DEFAULT FALSE,
  processed        BOOLEAN NOT NULL DEFAULT FALSE,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT unique_provider_external UNIQUE (provider, external_id)
);
