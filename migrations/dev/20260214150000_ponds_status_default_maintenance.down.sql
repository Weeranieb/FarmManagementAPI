-- Revert farms and ponds status default to 'active'
ALTER TABLE farms ALTER COLUMN status SET DEFAULT 'active';
ALTER TABLE ponds ALTER COLUMN status SET DEFAULT 'active';
