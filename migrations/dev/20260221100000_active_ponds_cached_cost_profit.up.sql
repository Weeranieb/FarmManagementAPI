-- Add cached cost/profit and total fish columns to active_ponds for faster dashboard reads
ALTER TABLE active_ponds
  ADD COLUMN total_cost numeric(20,4) NOT NULL DEFAULT 0,
  ADD COLUMN total_profit numeric(20,4) NOT NULL DEFAULT 0,
  ADD COLUMN net_result numeric(20,4) NOT NULL DEFAULT 0,
  ADD COLUMN total_fish integer NOT NULL DEFAULT 0;
