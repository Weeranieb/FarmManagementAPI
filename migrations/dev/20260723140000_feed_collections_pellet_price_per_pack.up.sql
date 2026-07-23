-- Pellet's canonical price unit becomes the bag (ถุง), matching how it's bought
-- (fresh already stores per ลัง). For each pellet price-history row, keep the
-- per-กก. figure in price_per_kg (that's what the weight-fed cost calc uses) and
-- rewrite price to the per-ถุง headline = per-กก. × pack size. Then relabel the
-- collection's unit. Existing pond feed costs are unchanged: cost stays
-- kg × price_per_kg. Fresh rows are untouched.
UPDATE feed_price_histories fph
   SET price_per_kg = COALESCE(fph.price_per_kg, fph.price),
       price        = COALESCE(fph.price_per_kg, fph.price) * fc.pack_size_kg
  FROM feed_collections fc
 WHERE fph.feed_collection_id = fc.id
   AND fc.feed_type = 'pellet';

UPDATE feed_collections
   SET unit = 'ถุง'
 WHERE feed_type = 'pellet';
