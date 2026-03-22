CREATE UNIQUE INDEX farm_on_farm_group_farm_group_farm_active_unique
ON farm_on_farm_group (farm_group_id, farm_id)
WHERE deleted_at IS NULL;
