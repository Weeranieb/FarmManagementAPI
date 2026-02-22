-- Revert activities columns to nullable (no default)
ALTER TABLE activities
  ALTER COLUMN amount DROP NOT NULL,
  ALTER COLUMN amount DROP DEFAULT;
ALTER TABLE activities
  ALTER COLUMN fish_type DROP NOT NULL,
  ALTER COLUMN fish_type DROP DEFAULT;
ALTER TABLE activities
  ALTER COLUMN fish_weight DROP NOT NULL,
  ALTER COLUMN fish_weight DROP DEFAULT;
ALTER TABLE activities
  ALTER COLUMN fish_unit DROP NOT NULL,
  ALTER COLUMN fish_unit DROP DEFAULT;
ALTER TABLE activities
  ALTER COLUMN price_per_unit DROP NOT NULL,
  ALTER COLUMN price_per_unit DROP DEFAULT;
