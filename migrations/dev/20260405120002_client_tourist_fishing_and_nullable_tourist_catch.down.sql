UPDATE daily_logs
SET tourist_catch_count = 0
WHERE tourist_catch_count IS NULL;

ALTER TABLE daily_logs
  ALTER COLUMN tourist_catch_count SET NOT NULL;

ALTER TABLE clients
  DROP COLUMN IF EXISTS is_tourist_fishing_enabled;
