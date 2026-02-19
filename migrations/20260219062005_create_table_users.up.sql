CREATE TYPE user_role AS ENUM (
  'admin',
  'teacher',
  'student',
  'parent',
  'staff'
);

CREATE TABLE users
(
  id             UUID PRIMARY KEY,
  tenant_id      UUID NOT NULL,
  name           TEXT NOT NULL,
  email          TEXT NOT NULL,
  password_hash  TEXT NOT NULL,
  role           user_role NOT NULL,
  is_active      BOOLEAN NOT NULL DEFAULT TRUE,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_user_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE CASCADE,

  CONSTRAINT unique_user_email_per_tenant
    UNIQUE (tenant_id, email)
);
