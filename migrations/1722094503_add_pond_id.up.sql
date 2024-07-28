-- Up Migration: Modify DailyFeeds table

-- Step 1: Make ActivePondId nullable
ALTER TABLE "DailyFeeds" 
ALTER COLUMN "ActivePondId" DROP NOT NULL;

-- Step 2: Add a new column PondId with NOT NULL constraint
ALTER TABLE "DailyFeeds"
ADD COLUMN "PondId" bigint NOT NULL;

-- If needed, create an index for PondId
CREATE INDEX ON "DailyFeeds" ("PondId");
