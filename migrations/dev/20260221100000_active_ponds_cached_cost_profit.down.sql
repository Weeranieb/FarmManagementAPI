-- Remove cached cost/profit and total fish columns from active_ponds
ALTER TABLE active_ponds
  DROP COLUMN IF EXISTS total_fish,
  DROP COLUMN IF EXISTS net_result,
  DROP COLUMN IF EXISTS total_profit,
  DROP COLUMN IF EXISTS total_cost;
