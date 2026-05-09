-- Add email column to users and enforce uniqueness among non-deleted rows.

ALTER TABLE users ADD COLUMN IF NOT EXISTS email varchar;

CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique_idx
  ON users (email)
  WHERE deleted_at IS NULL AND email IS NOT NULL;
