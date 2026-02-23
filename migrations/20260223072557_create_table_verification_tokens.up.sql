CREATE TABLE email_verification_tokens
(
  id           UUID PRIMARY KEY,
  user_id      UUID NOT NULL,
  token        TEXT NOT NULL,
  expires_at   TIMESTAMPTZ NOT NULL,
  used_at      TIMESTAMPTZ,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_email_token_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

  CONSTRAINT unique_email_verification_token
    UNIQUE (token)
);

CREATE INDEX idx_email_verification_user_id
ON email_verification_tokens (user_id);

CREATE INDEX idx_email_verification_expires_at
ON email_verification_tokens (expires_at);

CREATE INDEX idx_email_verification_unused
ON email_verification_tokens (used_at)
WHERE used_at IS NULL;
