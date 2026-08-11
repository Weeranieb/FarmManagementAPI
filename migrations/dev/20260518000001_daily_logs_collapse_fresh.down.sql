-- Restore split fresh_morning / fresh_evening columns.
-- The split is lossy: original morning/evening amounts cannot be recovered from the sum,
-- so the entire previous total is placed in fresh_morning and fresh_evening is set to 0.

ALTER TABLE daily_logs ADD COLUMN fresh_morning NUMERIC NOT NULL DEFAULT 0;
ALTER TABLE daily_logs ADD COLUMN fresh_evening NUMERIC NOT NULL DEFAULT 0;

UPDATE daily_logs SET fresh_morning = fresh, fresh_evening = 0;

ALTER TABLE daily_logs DROP COLUMN fresh;
