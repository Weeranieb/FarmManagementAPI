-- Restore old columns on sell_details
ALTER TABLE sell_details
  ADD COLUMN IF NOT EXISTS size VARCHAR(100),
  ADD COLUMN IF NOT EXISTS fish_type VARCHAR(100),
  ADD COLUMN IF NOT EXISTS amount NUMERIC,
  ADD COLUMN IF NOT EXISTS fish_unit VARCHAR(20);

-- Drop new columns from sell_details
ALTER TABLE sell_details
  DROP COLUMN IF EXISTS fish_size_grade_id,
  DROP COLUMN IF EXISTS weight;

-- Drop fish_size_grades table
DROP TABLE IF EXISTS fish_size_grades;
