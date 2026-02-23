CREATE TABLE subscription_plans
(
  id             UUID PRIMARY KEY,
  name           TEXT           NOT NULL,
  price_monthly  NUMERIC(15,2)  NOT NULL DEFAULT 0,
  price_yearly   NUMERIC(15,2)  NOT NULL DEFAULT 0,
  duration_days  INTEGER        NOT NULL DEFAULT 30,
  max_students   INTEGER,
  max_teachers   INTEGER,
  features_json  JSONB,
  is_active      BOOLEAN        NOT NULL DEFAULT TRUE,
  created_at     TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

ALTER TABLE subscription_plans
ADD CONSTRAINT subscription_plans_name_unique UNIQUE (name);
