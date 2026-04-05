CREATE TABLE daily_feeds (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  active_pond_id BIGINT NOT NULL,
  pond_id BIGINT NOT NULL,
  feed_collection_id BIGINT NOT NULL,
  morning_amount NUMERIC NOT NULL DEFAULT 0,
  evening_amount NUMERIC NOT NULL DEFAULT 0,
  feed_date DATE NOT NULL,
  deleted_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT (now()),
  created_by VARCHAR NOT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT (now()),
  updated_by VARCHAR NOT NULL
);

CREATE INDEX daily_feeds_active_pond_id_idx ON daily_feeds (active_pond_id);
CREATE INDEX daily_feeds_feed_collection_id_idx ON daily_feeds (feed_collection_id);
CREATE INDEX daily_feeds_pond_id_idx ON daily_feeds (pond_id);

ALTER TABLE daily_feeds ADD FOREIGN KEY (active_pond_id) REFERENCES active_ponds (id);
ALTER TABLE daily_feeds ADD FOREIGN KEY (pond_id) REFERENCES ponds (id);
ALTER TABLE daily_feeds ADD FOREIGN KEY (feed_collection_id) REFERENCES feed_collections (id);

INSERT INTO daily_feeds (
  active_pond_id,
  pond_id,
  feed_collection_id,
  morning_amount,
  evening_amount,
  feed_date,
  created_at,
  created_by,
  updated_at,
  updated_by,
  deleted_at
)
SELECT
  dl.active_pond_id,
  ap.pond_id,
  fc.id,
  dl.fresh_morning,
  dl.fresh_evening,
  dl.feed_date,
  dl.created_at,
  dl.created_by,
  dl.updated_at,
  dl.updated_by,
  dl.deleted_at
FROM daily_logs dl
INNER JOIN active_ponds ap ON ap.id = dl.active_pond_id
INNER JOIN feed_collections fc ON fc.id = dl.fresh_feed_collection_id AND fc.deleted_at IS NULL
WHERE dl.fresh_feed_collection_id IS NOT NULL
   AND (dl.fresh_morning <> 0 OR dl.fresh_evening <> 0 OR dl.deleted_at IS NOT NULL);

INSERT INTO daily_feeds (
  active_pond_id,
  pond_id,
  feed_collection_id,
  morning_amount,
  evening_amount,
  feed_date,
  created_at,
  created_by,
  updated_at,
  updated_by,
  deleted_at
)
SELECT
  dl.active_pond_id,
  ap.pond_id,
  fc.id,
  dl.pellet_morning,
  dl.pellet_evening,
  dl.feed_date,
  dl.created_at,
  dl.created_by,
  dl.updated_at,
  dl.updated_by,
  dl.deleted_at
FROM daily_logs dl
INNER JOIN active_ponds ap ON ap.id = dl.active_pond_id
INNER JOIN feed_collections fc ON fc.id = dl.pellet_feed_collection_id AND fc.deleted_at IS NULL
WHERE dl.pellet_feed_collection_id IS NOT NULL
   AND (dl.pellet_morning <> 0 OR dl.pellet_evening <> 0 OR dl.deleted_at IS NOT NULL);

DROP TABLE daily_logs;
