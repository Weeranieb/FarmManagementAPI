DROP INDEX IF EXISTS merchants_client_contact_uniq;
DROP INDEX IF EXISTS merchants_client_id_idx;
ALTER TABLE merchants DROP CONSTRAINT IF EXISTS merchants_client_id_fkey;
ALTER TABLE merchants DROP COLUMN IF EXISTS client_id;
