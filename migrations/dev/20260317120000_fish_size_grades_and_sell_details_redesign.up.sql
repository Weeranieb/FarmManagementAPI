-- Create fish_size_grades master data table
CREATE TABLE IF NOT EXISTS fish_size_grades (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  sort_index INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by VARCHAR(100) NOT NULL DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_by VARCHAR(100) NOT NULL DEFAULT '',
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_fish_size_grades_name ON fish_size_grades (name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_fish_size_grades_sort_index ON fish_size_grades (sort_index);
CREATE INDEX IF NOT EXISTS idx_fish_size_grades_deleted_at ON fish_size_grades (deleted_at);

-- Seed initial fish size grades with explicit IDs
INSERT INTO fish_size_grades (id, name, sort_index) VALUES
  ( 1, '8โล',     1),
  ( 2, '7โล',     2),
  ( 3, '6โล',     3),
  ( 4, '5โล',     4),
  ( 5, '4โล',     5),
  ( 6, '3โล',     6),
  ( 7, '2โล',     7),
  ( 8, 'โล',      8),
  ( 9, 'เป็น',     9),
  (10, 'โลเป็น',  10),
  (11, '4ขีด',    11),
  (12, '7ขีด',    12),
  (13, '8ขีด',    13),
  (14, '9ขีด',    14),
  (15, 'แผล(ญ)',  15),
  (16, 'แผล(ล)',  16),
  (17, 'ผอม',     17),
  (18, 'เล็ก',    18)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, sort_index = EXCLUDED.sort_index;

-- Reset the sequence so the next auto-increment value is above the seeded IDs
SELECT setval('fish_size_grades_id_seq', (SELECT COALESCE(MAX(id), 0) FROM fish_size_grades));

-- Redesign sell_details: add new columns
ALTER TABLE sell_details
  ADD COLUMN IF NOT EXISTS fish_size_grade_id INTEGER REFERENCES fish_size_grades(id),
  ADD COLUMN IF NOT EXISTS weight NUMERIC;

-- Drop old columns from sell_details
ALTER TABLE sell_details
  DROP COLUMN IF EXISTS size,
  DROP COLUMN IF EXISTS fish_type,
  DROP COLUMN IF EXISTS amount,
  DROP COLUMN IF EXISTS fish_unit;
