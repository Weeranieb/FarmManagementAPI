-- Re-add price_per_kg and reconstruct it as price / pack_size_kg — the inverse of
-- the original per-pack → per-กก. derivation. This restores the column that the
-- older migrations reference so a full `down` sequence stays valid.
ALTER TABLE feed_price_histories ADD COLUMN price_per_kg NUMERIC(20,4) NULL;

UPDATE feed_price_histories fph
   SET price_per_kg = fph.price / fc.pack_size_kg
  FROM feed_collections fc
 WHERE fph.feed_collection_id = fc.id
   AND fc.pack_size_kg > 0;
