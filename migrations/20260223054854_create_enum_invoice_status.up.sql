CREATE TYPE invoice_status AS ENUM (
  'pending',
  'paid',
  'expired',
  'cancelled'
);