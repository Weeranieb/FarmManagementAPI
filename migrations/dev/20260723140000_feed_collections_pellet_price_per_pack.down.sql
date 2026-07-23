-- Revert pellet prices to per-กก. (held in price_per_kg) and the unit to กก.
UPDATE feed_price_histories fph
   SET price = fph.price_per_kg
  FROM feed_collections fc
 WHERE fph.feed_collection_id = fc.id
   AND fc.feed_type = 'pellet'
   AND fph.price_per_kg IS NOT NULL;

UPDATE feed_collections
   SET unit = 'กก.'
 WHERE feed_type = 'pellet';
