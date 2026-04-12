ALTER TABLE active_ponds ADD COLUMN default_fresh_feed_collection_id BIGINT;
ALTER TABLE active_ponds ADD COLUMN default_pellet_feed_collection_id BIGINT;

ALTER TABLE active_ponds
  ADD CONSTRAINT active_ponds_default_fresh_feed_collection_id_fkey
    FOREIGN KEY (default_fresh_feed_collection_id) REFERENCES feed_collections (id);

ALTER TABLE active_ponds
  ADD CONSTRAINT active_ponds_default_pellet_feed_collection_id_fkey
    FOREIGN KEY (default_pellet_feed_collection_id) REFERENCES feed_collections (id);
