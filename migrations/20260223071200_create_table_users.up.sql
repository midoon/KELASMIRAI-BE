

CREATE TABLE users
(
  id             UUID PRIMARY KEY,
  tenant_id      UUID NOT NULL,
  name           TEXT NOT NULL,
  email          TEXT NOT NULL,
  password_hash  TEXT NOT NULL,
  role           user_role NOT NULL,
  activated_at   TIMESTAMPTZ,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_user_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE CASCADE,

  CONSTRAINT unique_user_email_per_tenant
    UNIQUE (tenant_id, email)
);
