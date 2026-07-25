-- Scope merchants to a client, and make the contact number unique per client.
-- Adds client_id (backfilling any pre-existing global rows to the lowest client
-- so NOT NULL + FK can be enforced) and a partial unique index so a blank
-- contact number is still allowed to repeat, but a real number cannot within
-- the same client. Soft-deleted rows are excluded so a number frees up on delete.

ALTER TABLE merchants ADD COLUMN client_id bigint;

UPDATE merchants
SET client_id = (SELECT id FROM clients ORDER BY id LIMIT 1)
WHERE client_id IS NULL;

ALTER TABLE merchants ALTER COLUMN client_id SET NOT NULL;

ALTER TABLE merchants
  ADD CONSTRAINT merchants_client_id_fkey FOREIGN KEY (client_id) REFERENCES clients (id);

CREATE INDEX merchants_client_id_idx ON merchants (client_id);

CREATE UNIQUE INDEX merchants_client_contact_uniq
  ON merchants (client_id, contact_number)
  WHERE contact_number <> '' AND deleted_at IS NULL;
