ALTER TABLE clients
  ADD COLUMN IF NOT EXISTS is_tourist_fishing_enabled BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE daily_logs
  ALTER COLUMN tourist_catch_count DROP NOT NULL;
