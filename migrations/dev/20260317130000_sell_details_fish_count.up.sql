-- Optional fish count per sell detail line (for recording head count when known)
ALTER TABLE sell_details
  ADD COLUMN IF NOT EXISTS fish_count INTEGER;
