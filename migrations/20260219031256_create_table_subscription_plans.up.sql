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

INSERT INTO subscription_plans
(id, name, price_monthly, price_yearly, duration_days, max_students, max_teachers, features_json, is_active)
VALUES

-- FREE
(
  '11111111-1111-1111-1111-111111111111',
  'Free',
  0,
  0,
  30,
  100,
  20,
  '{"modules": ["academic", "attendance"], "support": "community"}',
  TRUE
),

-- STANDARD
(
  '22222222-2222-2222-2222-222222222222',
  'Standard',
  299000,
  2990000,
  30,
  400,
  60,
  '{"modules": ["academic", "attendance", "finance", "report"], "support": "email"}',
  TRUE
),

-- PREMIUM
(
  '33333333-3333-3333-3333-333333333333',
  'Premium',
  599000,
  5990000,
  30,
  900,
  120,
  '{"modules": ["all"], "support": "priority_email"}',
  TRUE
),

-- EXCELLENT
(
  '44444444-4444-4444-4444-444444444444',
  'Excellent',
  999000,
  9990000,
  30,
  NULL, -- unlimited
  NULL, -- unlimited
  '{"modules": ["all"], "support": "priority_whatsapp", "dedicated_account_manager": true}',
  TRUE
);
