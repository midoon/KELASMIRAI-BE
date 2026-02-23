CREATE TYPE subscription_status AS ENUM (
  'trial',
  'active',
  'past_due',
  'cancelled'
);