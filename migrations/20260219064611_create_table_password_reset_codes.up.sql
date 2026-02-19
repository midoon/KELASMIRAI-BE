CREATE TABLE password_reset_codes
(
  id           UUID PRIMARY KEY,
  user_id      UUID NOT NULL,
  code         TEXT NOT NULL,
  expires_at   TIMESTAMPTZ NOT NULL,
  used_at      TIMESTAMPTZ,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_password_reset_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE UNIQUE INDEX uniq_active_password_reset_per_user
ON password_reset_codes (user_id)
WHERE used_at IS NULL;

CREATE INDEX idx_password_reset_expires_at
ON password_reset_codes (expires_at);
