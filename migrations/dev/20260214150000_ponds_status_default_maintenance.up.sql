-- Set default status for new farms and ponds to 'maintenance'
ALTER TABLE farms ALTER COLUMN status SET DEFAULT 'maintenance';
ALTER TABLE ponds ALTER COLUMN status SET DEFAULT 'maintenance';
