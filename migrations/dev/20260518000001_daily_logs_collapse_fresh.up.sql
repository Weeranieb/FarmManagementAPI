-- Collapse fresh_morning + fresh_evening into a single `fresh` column.
-- Fresh feed is given once per day in real operation; pellet stays split (morning/evening).
-- Existing data is preserved by summing the two old columns into the new one.

ALTER TABLE daily_logs ADD COLUMN fresh NUMERIC NOT NULL DEFAULT 0;

UPDATE daily_logs SET fresh = fresh_morning + fresh_evening;

ALTER TABLE daily_logs DROP COLUMN fresh_morning;
ALTER TABLE daily_logs DROP COLUMN fresh_evening;
