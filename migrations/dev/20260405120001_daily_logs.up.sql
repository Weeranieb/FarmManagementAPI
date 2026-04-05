-- Legacy init DB used `amount` only; align with application columns before pivot
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.tables
    WHERE table_schema = 'public' AND table_name = 'daily_feeds'
  ) THEN
    IF NOT EXISTS (
      SELECT 1 FROM information_schema.columns
      WHERE table_schema = 'public' AND table_name = 'daily_feeds' AND column_name = 'morning_amount'
    ) THEN
      ALTER TABLE daily_feeds ADD COLUMN morning_amount DOUBLE PRECISION NOT NULL DEFAULT 0;
      ALTER TABLE daily_feeds ADD COLUMN evening_amount DOUBLE PRECISION NOT NULL DEFAULT 0;
      IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'daily_feeds' AND column_name = 'amount'
      ) THEN
        UPDATE daily_feeds SET morning_amount = COALESCE(amount::double precision, 0), evening_amount = 0;
      END IF;
    END IF;
  END IF;
END $$;

CREATE TABLE daily_logs (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  active_pond_id BIGINT NOT NULL,
  feed_date DATE NOT NULL,
  fresh_feed_collection_id BIGINT,
  pellet_feed_collection_id BIGINT,
  fresh_morning NUMERIC NOT NULL DEFAULT 0,
  fresh_evening NUMERIC NOT NULL DEFAULT 0,
  pellet_morning NUMERIC NOT NULL DEFAULT 0,
  pellet_evening NUMERIC NOT NULL DEFAULT 0,
  death_fish_count INTEGER NOT NULL DEFAULT 0,
  tourist_catch_count INTEGER NOT NULL DEFAULT 0,
  deleted_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT (now()),
  created_by VARCHAR NOT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT (now()),
  updated_by VARCHAR NOT NULL
);

CREATE UNIQUE INDEX daily_logs_active_pond_feed_date_uidx
  ON daily_logs (active_pond_id, feed_date)
  WHERE deleted_at IS NULL;

CREATE INDEX daily_logs_active_pond_id_idx ON daily_logs (active_pond_id);
CREATE INDEX daily_logs_fresh_fc_idx ON daily_logs (fresh_feed_collection_id);
CREATE INDEX daily_logs_pellet_fc_idx ON daily_logs (pellet_feed_collection_id);

ALTER TABLE daily_logs ADD FOREIGN KEY (active_pond_id) REFERENCES active_ponds (id);
ALTER TABLE daily_logs ADD FOREIGN KEY (fresh_feed_collection_id) REFERENCES feed_collections (id);
ALTER TABLE daily_logs ADD FOREIGN KEY (pellet_feed_collection_id) REFERENCES feed_collections (id);

-- Migrate from daily_feeds (expects morning_amount / evening_amount columns as in application model)
WITH base AS (
  SELECT df.active_pond_id,
         df.feed_date,
         fc.feed_type,
         df.feed_collection_id,
         df.morning_amount::numeric AS morning_amount,
         df.evening_amount::numeric AS evening_amount
  FROM daily_feeds df
  INNER JOIN feed_collections fc ON fc.id = df.feed_collection_id AND fc.deleted_at IS NULL
  WHERE df.deleted_at IS NULL
),
by_type AS (
  SELECT active_pond_id,
         feed_date,
         feed_type,
         MAX(feed_collection_id) AS feed_collection_id,
         SUM(morning_amount) AS morning_sum,
         SUM(evening_amount) AS evening_sum
  FROM base
  GROUP BY active_pond_id, feed_date, feed_type
),
fresh AS (SELECT * FROM by_type WHERE feed_type = 'fresh'),
pellet AS (SELECT * FROM by_type WHERE feed_type = 'pellet')
INSERT INTO daily_logs (
  active_pond_id,
  feed_date,
  fresh_feed_collection_id,
  pellet_feed_collection_id,
  fresh_morning,
  fresh_evening,
  pellet_morning,
  pellet_evening,
  death_fish_count,
  tourist_catch_count,
  created_at,
  created_by,
  updated_at,
  updated_by,
  deleted_at
)
SELECT
  COALESCE(fresh.active_pond_id, pellet.active_pond_id),
  COALESCE(fresh.feed_date, pellet.feed_date),
  fresh.feed_collection_id,
  pellet.feed_collection_id,
  COALESCE(fresh.morning_sum, 0),
  COALESCE(fresh.evening_sum, 0),
  COALESCE(pellet.morning_sum, 0),
  COALESCE(pellet.evening_sum, 0),
  0,
  0,
  NOW(),
  'migration',
  NOW(),
  'migration',
  NULL
FROM fresh
FULL OUTER JOIN pellet
  ON fresh.active_pond_id = pellet.active_pond_id
 AND fresh.feed_date = pellet.feed_date;

DROP TABLE daily_feeds;
