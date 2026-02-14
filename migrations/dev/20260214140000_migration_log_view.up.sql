-- Table of known migrations (add a row here when you add a new migration file)
CREATE TABLE migration_info (
  version BIGINT PRIMARY KEY,
  name TEXT NOT NULL
);

INSERT INTO migration_info (version, name) VALUES
  (20251026101634, 'init_db'),
  (20260111164341, 'fix_is_active_bit_to_boolean'),
  (20260204200437, 'remove_code_fields');

-- View: list of applied migrations (query this instead of schema_migrations for a clear log)
CREATE VIEW migration_log AS
SELECT m.version, m.name
FROM migration_info m
WHERE m.version <= (SELECT version FROM schema_migrations LIMIT 1)
ORDER BY m.version;
