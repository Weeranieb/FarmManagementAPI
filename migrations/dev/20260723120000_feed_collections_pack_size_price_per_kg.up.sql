-- Pack size is required going forward: every feed is bought in a pack, and the
-- per-กก. price derivation depends on it. Add nullable, backfill existing rows
-- with a sensible default per feed type (fresh crate ≈ 30 กก., pellet bag ≈ 20
-- กก.), then enforce NOT NULL. New rows get their default from the service
-- (constants.DefaultPackSizeKg) since the default is type-dependent.
ALTER TABLE feed_collections ADD COLUMN pack_size_kg NUMERIC(12,4);

UPDATE feed_collections
   SET pack_size_kg = CASE WHEN feed_type = 'fresh' THEN 30 ELSE 20 END
 WHERE pack_size_kg IS NULL;

ALTER TABLE feed_collections ALTER COLUMN pack_size_kg SET NOT NULL;

-- Price-per-กก. is snapshotted per price-history entry; stays nullable for rows
-- recorded before it existed (no lossless backfill without the pack size at the
-- time of entry).
ALTER TABLE feed_price_histories ADD COLUMN price_per_kg NUMERIC(20,4) NULL;
