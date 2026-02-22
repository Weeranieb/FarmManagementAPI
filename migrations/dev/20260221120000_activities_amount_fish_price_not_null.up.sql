-- Backfill NULLs so we can set NOT NULL
UPDATE activities SET amount = COALESCE(amount, 0) WHERE amount IS NULL;
UPDATE activities SET fish_type = COALESCE(fish_type, '') WHERE fish_type IS NULL;
UPDATE activities SET fish_weight = COALESCE(fish_weight, 0) WHERE fish_weight IS NULL;
UPDATE activities SET fish_unit = COALESCE(fish_unit, '') WHERE fish_unit IS NULL;
UPDATE activities SET price_per_unit = COALESCE(price_per_unit, 0) WHERE price_per_unit IS NULL;

-- Make columns NOT NULL with defaults
ALTER TABLE activities
  ALTER COLUMN amount SET DEFAULT 0,
  ALTER COLUMN amount SET NOT NULL;
ALTER TABLE activities
  ALTER COLUMN fish_type SET DEFAULT '',
  ALTER COLUMN fish_type SET NOT NULL;
ALTER TABLE activities
  ALTER COLUMN fish_weight SET DEFAULT 0,
  ALTER COLUMN fish_weight SET NOT NULL;
ALTER TABLE activities
  ALTER COLUMN fish_unit SET DEFAULT '',
  ALTER COLUMN fish_unit SET NOT NULL;
ALTER TABLE activities
  ALTER COLUMN price_per_unit SET DEFAULT 0,
  ALTER COLUMN price_per_unit SET NOT NULL;
