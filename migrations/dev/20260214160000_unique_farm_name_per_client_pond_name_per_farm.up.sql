-- Unique farm name per client (soft-delete aware)
CREATE UNIQUE INDEX farms_client_id_name_key ON farms (client_id, name) WHERE deleted_at IS NULL;

-- Unique pond name per farm (soft-delete aware)
CREATE UNIQUE INDEX ponds_farm_id_name_key ON ponds (farm_id, name) WHERE deleted_at IS NULL;
