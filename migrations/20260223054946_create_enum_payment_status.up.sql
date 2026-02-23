CREATE TYPE payment_status AS ENUM (
  'pending',
  'settlement',
  'expire',
  'cancel',
  'deny'
);