-- Add nullable feed_cost snapshot to active_ponds.
-- NULL means the cycle's feed cost has not been snapshotted (legacy rows closed
-- before this feature, or still-active cycles whose feed cost is derived on read).
-- A non-NULL value is the frozen total feed cost captured at close.
ALTER TABLE active_ponds
  ADD COLUMN feed_cost numeric(20,4);
