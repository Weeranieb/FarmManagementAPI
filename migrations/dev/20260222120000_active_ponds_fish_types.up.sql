-- Add fish_types JSONB array to active_ponds (e.g. ["nil", "kaphong"])
ALTER TABLE active_ponds
  ADD COLUMN fish_types jsonb NOT NULL DEFAULT '[]';
