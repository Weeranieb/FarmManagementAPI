-- Feed collections are canonical on active_ponds; remove redundant per-row columns from daily_logs.

ALTER TABLE active_ponds RENAME COLUMN default_fresh_feed_collection_id TO fresh_feed_collection_id;
ALTER TABLE active_ponds RENAME COLUMN default_pellet_feed_collection_id TO pellet_feed_collection_id;

ALTER TABLE daily_logs DROP COLUMN IF EXISTS fresh_feed_collection_id CASCADE;
ALTER TABLE daily_logs DROP COLUMN IF EXISTS pellet_feed_collection_id CASCADE;
