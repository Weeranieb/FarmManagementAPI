-- Down Migration: Revert changes to DailyFeeds table

-- Step 1: Drop the PondId column
ALTER TABLE "DailyFeeds"
DROP COLUMN "PondId";

-- Step 2: Make ActivePondId non-nullable
ALTER TABLE "DailyFeeds" 
ALTER COLUMN "ActivePondId" SET NOT NULL;

-- Optionally, drop the index if it was created
DROP INDEX IF EXISTS "DailyFeeds_PondId_idx";
