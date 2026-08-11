#!/usr/bin/env bash
# ─────────────────────────────────────────────────────────────────────────────
# Daily logical backup of the Boonma Farm production database.
#
#   1. pg_dump (custom format, compressed) from the running db container
#   2. upload to DigitalOcean Spaces (S3-compatible) in the SGP1 region
#   3. prune LOCAL dumps older than $RETENTION_DAYS (remote retention: Spaces
#      lifecycle rule — see DEPLOY.md)
#
# Install cron as the deploy user (02:15 Asia/Bangkok daily):
#   15 2 * * * /opt/boonmafarm/deploy/scripts/backup.sh >> /var/log/boonmafarm-backup.log 2>&1
#
# Config comes from deploy/.env — add these keys to it:
#   SPACES_BUCKET      e.g. boonmafarm-backups
#   SPACES_ENDPOINT    e.g. https://sgp1.digitaloceanspaces.com
#   AWS_ACCESS_KEY_ID  / AWS_SECRET_ACCESS_KEY   (a Spaces access key pair)
#
# Requires the AWS CLI on the droplet:  sudo apt-get install -y awscli
# ─────────────────────────────────────────────────────────────────────────────
set -euo pipefail

DEPLOY_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$DEPLOY_DIR"

# Load config (POSTGRES_*, SPACES_*, AWS_*) from the single .env.
set -a
[ -f .env ] && . ./.env
set +a

: "${POSTGRES_USER:=postgres}"
: "${POSTGRES_DB:=boonmafarm}"
RETENTION_DAYS="${RETENTION_DAYS:-7}"
BACKUP_DIR="${BACKUP_DIR:-$DEPLOY_DIR/backups}"

for var in SPACES_BUCKET SPACES_ENDPOINT AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY; do
  if [ -z "${!var:-}" ]; then
    echo "FATAL: $var is not set (add it to .env or .env.backup)" >&2
    exit 1
  fi
done

mkdir -p "$BACKUP_DIR"
STAMP="$(date +%Y%m%d-%H%M%S)"
FILE="$BACKUP_DIR/boonmafarm-${STAMP}.dump"

echo "[$(date -Is)] dumping ${POSTGRES_DB} -> ${FILE}"
# -Fc = custom format (compressed, restorable with pg_restore, supports parallel).
docker compose exec -T db \
  pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" -Fc > "$FILE"

# Fail loudly if the dump is suspiciously tiny (e.g. auth failure wrote nothing).
if [ "$(stat -c%s "$FILE" 2>/dev/null || stat -f%z "$FILE")" -lt 1024 ]; then
  echo "FATAL: dump is < 1 KiB — treating as failed" >&2
  rm -f "$FILE"
  exit 1
fi

echo "[$(date -Is)] uploading to s3://${SPACES_BUCKET}/db/$(basename "$FILE")"
aws s3 cp "$FILE" "s3://${SPACES_BUCKET}/db/$(basename "$FILE")" \
  --endpoint-url "$SPACES_ENDPOINT"

echo "[$(date -Is)] pruning local dumps older than ${RETENTION_DAYS} days"
find "$BACKUP_DIR" -name 'boonmafarm-*.dump' -type f -mtime "+${RETENTION_DAYS}" -delete

echo "[$(date -Is)] backup OK"
