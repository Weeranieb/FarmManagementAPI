ALTER TABLE daily_logs ADD COLUMN fresh_feed_collection_id BIGINT;
ALTER TABLE daily_logs ADD COLUMN pellet_feed_collection_id BIGINT;

CREATE INDEX IF NOT EXISTS daily_logs_fresh_fc_idx ON daily_logs (fresh_feed_collection_id);
CREATE INDEX IF NOT EXISTS daily_logs_pellet_fc_idx ON daily_logs (pellet_feed_collection_id);

ALTER TABLE daily_logs
  ADD CONSTRAINT daily_logs_fresh_feed_collection_id_fkey
  FOREIGN KEY (fresh_feed_collection_id) REFERENCES feed_collections (id);
ALTER TABLE daily_logs
  ADD CONSTRAINT daily_logs_pellet_feed_collection_id_fkey
  FOREIGN KEY (pellet_feed_collection_id) REFERENCES feed_collections (id);

ALTER TABLE active_ponds RENAME COLUMN fresh_feed_collection_id TO default_fresh_feed_collection_id;
ALTER TABLE active_ponds RENAME COLUMN pellet_feed_collection_id TO default_pellet_feed_collection_id;
