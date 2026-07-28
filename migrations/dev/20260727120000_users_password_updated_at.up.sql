-- Track when a user's current password was set.
--
-- users.updated_at moves on any profile edit (name, email, contact number), so
-- it cannot answer "when did I last change my password". This column is written
-- only by the three code paths that actually write users.password: account
-- creation, self-service change, and super-admin reset.
--
-- Deliberately left NULL for existing rows: the password of an account created
-- before this migration was set at an unknown time, and back-filling created_at
-- would assert "never changed since signup", which we don't know. Clients must
-- treat NULL as "unknown" rather than rendering a date.

-- timestamptz, not the bare `timestamp` used by created_at/updated_at: those
-- store a local wall clock and hand it back tagged as UTC, so a value written
-- at 23:00 Bangkok reads back as 23:00Z and renders as the *next* day for
-- anyone east of Greenwich. This column is shown to users as a date, so the
-- seven-hour drift would be visible. New column, no back-compat to preserve.

ALTER TABLE users ADD COLUMN IF NOT EXISTS password_updated_at timestamptz;
