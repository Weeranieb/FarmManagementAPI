#!/bin/bash
# ─────────────────────────────────────────────────────────────────────────────
# Postgres role bootstrap — runs ONCE, on first cluster init, via the official
# image's /docker-entrypoint-initdb.d hook. Executes as the bootstrap superuser
# ($POSTGRES_USER) against the application database ($POSTGRES_DB).
#
# Least-privilege model (NOT to be confused with application login users):
#
#   postgres       superuser        maintenance / break-glass only — app NEVER uses it
#   sys_admin      schema owner      runs migrations (CREATE/ALTER/DROP in public); NOT superuser
#   app_user       DML only          the API connects as this (SELECT/INSERT/UPDATE/DELETE)
#   readonly_user  SELECT only       reporting / analytics / dashboards over the SSH tunnel
#
# Because migrations run AS sys_admin, we set ALTER DEFAULT PRIVILEGES FOR ROLE
# sys_admin so every future table/sequence/view it creates auto-grants the right
# access to app_user and readonly_user — no manual GRANT after each migration.
#
# Re-running this on an existing cluster is a no-op (initdb hooks only fire on an
# empty data dir). To re-assert grants on a live DB (e.g. after a restore), see the
# "Backups → Restore" SQL snippet in DEPLOY.md.
# ─────────────────────────────────────────────────────────────────────────────
set -euo pipefail

: "${APP_DB_USER:=app_user}"
: "${READONLY_DB_USER:=readonly_user}"
: "${SYS_ADMIN_DB_USER:=sys_admin}"

for var in APP_DB_PASSWORD READONLY_DB_PASSWORD SYS_ADMIN_DB_PASSWORD; do
  if [ -z "${!var:-}" ]; then
    echo "FATAL: $var must be set to initialize database roles" >&2
    exit 1
  fi
done

psql -v ON_ERROR_STOP=1 \
     --username "$POSTGRES_USER" \
     --dbname "$POSTGRES_DB" \
     --set app_user="$APP_DB_USER" \
     --set readonly_user="$READONLY_DB_USER" \
     --set sys_admin="$SYS_ADMIN_DB_USER" \
     --set app_password="$APP_DB_PASSWORD" \
     --set readonly_password="$READONLY_DB_PASSWORD" \
     --set sys_admin_password="$SYS_ADMIN_DB_PASSWORD" <<'SQL'
-- Store the new roles' passwords as SCRAM, never md5.
SET password_encryption = 'scram-sha-256';

-- ── Login roles ─────────────────────────────────────────────────────────────
CREATE ROLE :"sys_admin"     LOGIN PASSWORD :'sys_admin_password';
CREATE ROLE :"app_user"      LOGIN PASSWORD :'app_password';
CREATE ROLE :"readonly_user" LOGIN PASSWORD :'readonly_password';

-- ── Schema ownership ────────────────────────────────────────────────────────
-- Give the public schema to sys_admin so migrations need no superuser.
ALTER SCHEMA public OWNER TO :"sys_admin";

-- PG15 already removes CREATE-on-public from PUBLIC, but be explicit.
REVOKE CREATE ON SCHEMA public FROM PUBLIC;

-- Every login role needs USAGE to reference objects in the schema.
GRANT USAGE ON SCHEMA public TO :"app_user", :"readonly_user";

-- ── Privileges on objects that already exist (none on first init; safe) ─────
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES    IN SCHEMA public TO :"app_user";
GRANT USAGE, SELECT                  ON ALL SEQUENCES  IN SCHEMA public TO :"app_user";
GRANT SELECT                         ON ALL TABLES     IN SCHEMA public TO :"readonly_user";

-- ── Default privileges for objects sys_admin creates LATER (migrations) ─────
ALTER DEFAULT PRIVILEGES FOR ROLE :"sys_admin" IN SCHEMA public
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO :"app_user";
ALTER DEFAULT PRIVILEGES FOR ROLE :"sys_admin" IN SCHEMA public
  GRANT USAGE, SELECT ON SEQUENCES TO :"app_user";
ALTER DEFAULT PRIVILEGES FOR ROLE :"sys_admin" IN SCHEMA public
  GRANT SELECT ON TABLES TO :"readonly_user";

-- app_user / readonly_user deliberately get NO CREATE/ALTER/DROP: they don't own
-- objects and lack CREATE on the schema, so DDL is impossible for them.
SQL

echo "✅ Roles ready: ${SYS_ADMIN_DB_USER} (schema owner), ${APP_DB_USER} (DML), ${READONLY_DB_USER} (read-only)"
