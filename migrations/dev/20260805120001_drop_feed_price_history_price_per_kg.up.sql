-- Pellet is now priced per ถุง directly — the headline `price` is the per-pack
-- price for both feed types, and daily pellet feeding is logged in ถุง too, so
-- cost = ถุง × price. The derived per-กก. snapshot (price_per_kg) is no longer
-- read by any cost math and is dropped.
ALTER TABLE feed_price_histories DROP COLUMN IF EXISTS price_per_kg;
