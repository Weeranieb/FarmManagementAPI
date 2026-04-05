ALTER TABLE feed_collections
  ADD COLUMN feed_type VARCHAR(20) NOT NULL DEFAULT 'pellet',
  ADD COLUMN fcr NUMERIC(12,4) NULL;

ALTER TABLE feed_collections
  ADD CONSTRAINT feed_collections_feed_type_check
  CHECK (feed_type IN ('fresh', 'pellet'));
