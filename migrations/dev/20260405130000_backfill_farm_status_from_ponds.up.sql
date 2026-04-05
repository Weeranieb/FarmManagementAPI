-- Derive farm status from ponds: active if any non-deleted pond is active; else maintenance.
UPDATE farms f
SET status = CASE
  WHEN EXISTS (
    SELECT 1
    FROM ponds p
    WHERE p.farm_id = f.id
      AND p.deleted_at IS NULL
      AND p.status = 'active'
  ) THEN 'active'
  ELSE 'maintenance'
END
WHERE f.deleted_at IS NULL;
