-- Remove feed_cost snapshot column from active_ponds
ALTER TABLE active_ponds
  DROP COLUMN IF EXISTS feed_cost;
