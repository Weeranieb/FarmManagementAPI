ALTER TABLE feed_collections DROP CONSTRAINT IF EXISTS feed_collections_feed_type_check;
ALTER TABLE feed_collections DROP COLUMN IF EXISTS fcr;
ALTER TABLE feed_collections DROP COLUMN IF EXISTS feed_type;
