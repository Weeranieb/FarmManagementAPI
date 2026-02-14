-- Add code field back to feed_collections table
ALTER TABLE feed_collections ADD COLUMN code varchar NOT NULL DEFAULT '';

-- Remove status from ponds table
ALTER TABLE ponds DROP COLUMN IF EXISTS status;

-- Add code field back to ponds table
ALTER TABLE ponds ADD COLUMN code varchar NOT NULL DEFAULT '';

-- Add code field back to farm_groups table
ALTER TABLE farm_groups ADD COLUMN code varchar NOT NULL DEFAULT '';

-- Remove status from farms table
ALTER TABLE farms DROP COLUMN IF EXISTS status;

-- Add code field back to farms table
ALTER TABLE farms ADD COLUMN code varchar NOT NULL DEFAULT '';
