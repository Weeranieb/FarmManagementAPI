-- Remove code field from farms table
ALTER TABLE farms DROP COLUMN code;

-- Add status to farms table
ALTER TABLE farms ADD COLUMN status varchar NOT NULL DEFAULT 'active';

-- Remove code field from farm_groups table
ALTER TABLE farm_groups DROP COLUMN code;

-- Remove code field from ponds table
ALTER TABLE ponds DROP COLUMN code;

-- Add status to ponds table
ALTER TABLE ponds ADD COLUMN status varchar NOT NULL DEFAULT 'active';

-- Remove code field from feed_collections table
ALTER TABLE feed_collections DROP COLUMN code;
