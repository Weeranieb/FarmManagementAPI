# Boonma Farm — Production Deployment Runbook

Deploys the Go/Fiber API + self-hosted PostgreSQL to the DigitalOcean droplet,
behind Caddy with automatic TLS. Low traffic (10–20 req/morning), single $6/mo
droplet, no managed services.

```
                         Internet
                            │  :443 (TLS)
              ┌─────────────▼─────────────┐   droplet 168.144.99.22 (SGP1)
              │  caddy   api.boonmafarm.com│   /opt/boonmafarm/deploy
              │   └─ reverse_proxy ──┐     │
              │  app  :8080 ◄────────┘     │   image pulled from GHCR
              │   └─ DATABASE_USER=app_user│
              │  db   :5432 (127.0.0.1 only)│  ← never public; SSH tunnel only
              └────────────────────────────┘
   app.boonmafarm.com (prod React)  → Vercel      ─┐ CORS: app is allowed to
   dev.boonmafarm.com (dev React)   → Vercel      ─┘ call api.boonmafarm.com
```

**Roles / responsibilities**
| Layer | Who does it | Where |
|---|---|---|
| Build image, push to GHCR | GitHub Actions (`deploy-prod.yml`) | on push to `main` |
| Restart `app` on the droplet | GitHub Actions over SSH | `db`/`caddy` untouched → no DB downtime |
| Schema migrations | GitHub Actions (`migrate-prod.yml`), **manual + gated** | over an SSH tunnel, as `sys_admin` |
| TLS certs | Caddy, automatically | Let's Encrypt |

> Files live in the `Weeranieb/FarmManagementAPI` repo under `deploy/` + `Dockerfile`
> + `.github/workflows/`. The droplet runs from a checkout at `/opt/boonmafarm`.

---

## 1. One-time droplet setup

SSH in as `deploy` (root login is already disabled):

```bash
ssh deploy@168.144.99.22
```

Clone the repo and create the runtime dir:

```bash
sudo mkdir -p /opt/boonmafarm && sudo chown "$USER":"$USER" /opt/boonmafarm
git clone https://github.com/Weeranieb/FarmManagementAPI.git /opt/boonmafarm
cd /opt/boonmafarm/deploy
```

Create the secrets file from the template and fill in **every** required value:

```bash
cp .env.prod.example .env
chmod 600 .env
# generate each secret: openssl rand -base64 32 | tr -d '/+=' | cut -c1-40
nano .env
```

Log the droplet in to GHCR so it can pull the (private) image. Create a GitHub
**classic PAT** with only `read:packages`, then:

```bash
echo 'YOUR_GHCR_READ_PAT' | docker login ghcr.io -u Weeranieb --password-stdin
```

*(If you make the GHCR package public instead, you can skip this login.)*

---

## 2. DNS (GoDaddy)

In GoDaddy DNS for `boonmafarm.com`, add:

| Type | Name | Value | TTL |
|---|---|---|---|
| A | `api` | `168.144.99.22` | 600 |

That's all the droplet needs. `app` / `dev` / root stay pointed at Vercel.
Confirm it resolves before starting Caddy:

```bash
dig +short api.boonmafarm.com   # → 168.144.99.22
```

`ufw` already allows 80/443, which is what Caddy's ACME challenge needs.

---

## 3. First boot

Bring up the database and proxy first, let Postgres initialize the roles, then
migrate, then start the app:

```bash
cd /opt/boonmafarm/deploy

# 1) db + caddy. On the very first run, the initdb script creates
#    sys_admin / app_user / readonly_user from the passwords in .env.
docker compose up -d db caddy
docker compose logs -f db      # watch for "✅ Roles ready: ..."; Ctrl-C when done

# 2) TLS: Caddy auto-issues the cert once DNS + ports are good.
docker compose logs caddy | grep -i "certificate obtained"   # may take ~30s
curl -I https://api.boonmafarm.com/health                     # → HTTP/2 200

# 3) Run the schema migration (see §5 for the CI-driven way; on first boot you
#    can also run it locally against the loopback DB as sys_admin):
docker run --rm -v "$PWD/../migrations:/m" --network host migrate/migrate \
  -path=/m/dev \
  -database "postgres://sys_admin:SYS_ADMIN_PW@127.0.0.1:5432/boonmafarm?sslmode=disable" up

# 4) Start the app. (No manual grants needed: migrations run as sys_admin, and
#    ALTER DEFAULT PRIVILEGES already grants app_user/readonly_user on its tables.)
docker compose up -d app
docker compose ps
curl https://api.boonmafarm.com/ready    # → 200 once DB is reachable with tables
```

---

## 4. GitHub secrets (for CI/CD)

Create a **`prod` environment** in the repo (Settings → Environments) and add
required reviewers so deploys/migrations need an approval click. Add these
secrets to the `prod` environment:

| Secret | Value |
|---|---|
| `PROD_SSH_HOST` | `168.144.99.22` |
| `PROD_SSH_USER` | `deploy` |
| `PROD_SSH_KEY`  | private key (PEM) whose public key is in `deploy`'s `~/.ssh/authorized_keys` |
| `PROD_SSH_PORT` | `22` (optional) |
| `PROD_MIGRATE_DATABASE_URL` | `postgres://sys_admin:SYS_ADMIN_PW@127.0.0.1:5432/boonmafarm?sslmode=disable` |

Generate a dedicated deploy key (don't reuse a personal key):

```bash
ssh-keygen -t ed25519 -C "gh-actions-deploy" -f ./gh_deploy -N ""
# put gh_deploy.pub into the droplet's deploy user:
ssh deploy@168.144.99.22 'cat >> ~/.ssh/authorized_keys' < gh_deploy.pub
# paste the PRIVATE key (gh_deploy) into the PROD_SSH_KEY secret, then delete both files
```

`GITHUB_TOKEN` (auto) handles the GHCR push — no PAT needed for that side.

---

## 5. Everyday operations

**Deploy** (automatic): merge/push to `main` → `deploy-prod.yml` builds the image,
pushes `:latest` + `:sha-<commit>` to GHCR, SSHes in, and runs
`docker compose pull app && up -d app`. Only `app` restarts.

**Roll back** to a previous image without a rebuild:

```bash
ssh deploy@168.144.99.22
cd /opt/boonmafarm/deploy
APP_IMAGE=ghcr.io/weeranieb/farmmanagementapi:sha-<old_commit> docker compose up -d app
# (or edit APP_IMAGE in .env to pin it, then `docker compose up -d app`)
```

**Migrate** (manual, gated): Actions → *Migrate Prod Database* → Run workflow →
pick `up` (or `down` to roll back N steps). It opens the SSH tunnel and runs
golang-migrate as `sys_admin`. Run this **before** the deploy that needs the new
columns, or right after — never bundle it into the deploy. If a migration errors and
leaves the version *dirty*, fix it manually over the SSH tunnel with
`migrate -path migrations/dev -database "$URL" force <version>` (kept off the UI on purpose).

**Update droplet config** (compose/Caddyfile/init script changed):

```bash
ssh deploy@168.144.99.22 'cd /opt/boonmafarm && git pull && cd deploy && docker compose up -d'
```

---

## 6. Reaching the database (no public port)

Port 5432 is bound to `127.0.0.1` on the droplet. To connect from your laptop,
tunnel over SSH — read-only work should use `readonly_user`:

```bash
ssh -N -L 5432:127.0.0.1:5432 deploy@168.144.99.22 &
psql "postgres://readonly_user:READONLY_PW@127.0.0.1:5432/boonmafarm"
```

Use `sys_admin` for DDL/index work, `postgres` only for break-glass maintenance.
Never open 5432 in `ufw`.

---

## 7. Auth / password hardening (already done in code)

- App passwords are hashed with **bcrypt cost 12** (`src/internal/service`,
  constant `bcryptCost`). Never stored in plaintext; the column is `password`
  holding the hash. Verify on login with `bcrypt.CompareHashAndPassword`.
- JWT is HS256 signed with `AUTHENTICATION_JWT_SECRET` — keep it 32+ random chars;
  the app refuses to boot if it's empty.
- Legacy cost-10 hashes still verify fine (bcrypt embeds the cost in each hash);
  only new/changed passwords use cost 12. No migration or password reset needed.

**Verify Postgres uses SCRAM (not md5/trust):**

```bash
docker compose exec db psql -U postgres -d boonmafarm -c "SHOW password_encryption;"       # scram-sha-256
docker compose exec db psql -U postgres -d boonmafarm -c "SELECT rolname FROM pg_authid WHERE rolpassword LIKE 'md5%';"  # 0 rows
docker compose exec db grep -E '^(host|local)' /var/lib/postgresql/data/pg_hba.conf         # methods = scram-sha-256
```

**Rotate a role's password (zero downtime):** `ALTER ROLE` only affects *new*
connections, so change it in Postgres first, then roll the credential into the app:

```bash
NEW=$(openssl rand -base64 32 | tr -d '/+=' | cut -c1-40)
docker compose exec -T db psql -U postgres -d boonmafarm \
  -c "SET password_encryption='scram-sha-256'; ALTER ROLE app_user PASSWORD '$NEW';"
# put $NEW in .env (APP_DB_PASSWORD), then: docker compose up -d app   # db keeps running
```

For `sys_admin`, also update the `PROD_MIGRATE_DATABASE_URL` GitHub secret.

---

## 8. Backups

Create a DO Space (SGP1) + a Spaces access key, add its keys to `.env`
(`SPACES_BUCKET`, `SPACES_ENDPOINT`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`),
install the AWS CLI, and schedule the dump:

```bash
sudo apt-get install -y awscli
crontab -e
# 02:15 Asia/Bangkok daily:
15 2 * * * /opt/boonmafarm/deploy/scripts/backup.sh >> /var/log/boonmafarm-backup.log 2>&1
```

`backup.sh` runs `pg_dump -Fc`, uploads to `s3://<bucket>/db/`, and prunes local
copies older than `RETENTION_DAYS`. Set a Spaces **lifecycle rule** to expire
remote objects (e.g. 30 days). **Restore:**

```bash
aws s3 cp s3://<bucket>/db/boonmafarm-YYYYMMDD-HHMMSS.dump ./restore.dump --endpoint-url https://sgp1.digitaloceanspaces.com
cat restore.dump | docker compose exec -T db pg_restore -U postgres -d boonmafarm --clean --if-exists
# Re-assert grants: restored objects are owned by postgres, so the sys_admin default
# privileges don't apply. Run once after a restore:
docker compose exec -T db psql -U postgres -d boonmafarm <<'SQL'
ALTER SCHEMA public OWNER TO sys_admin;
GRANT USAGE ON SCHEMA public TO app_user, readonly_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES   IN SCHEMA public TO app_user;
GRANT USAGE, SELECT                 ON ALL SEQUENCES IN SCHEMA public TO app_user;
GRANT SELECT                        ON ALL TABLES    IN SCHEMA public TO readonly_user;
SQL
```

---

## 9. Monitoring

Point a free **UptimeRobot** monitor at `https://api.boonmafarm.com/health`
(HTTP 200, liveness — no DB dependency). Add a second monitor on `/ready` if you
want DB-connectivity alerting. Caddy also renews certs automatically; no cron for TLS.

---

## 10. Frontend wiring (Vercel)

The prod React app on `app.boonmafarm.com` must call the droplet API. In Vercel
**Production** env:

```
VITE_API_BASE_URL=https://api.boonmafarm.com
```

*(It's a Vite SPA — the var is `VITE_API_BASE_URL`, not `NEXT_PUBLIC_*`. In dev the
client uses the Vite proxy `/api/v1`, so this var only matters for prod builds.)*

**Auth is cookie-based** — `src/lib/api-client.ts` sends `credentials: 'include'`
and the API sets an HttpOnly `jwt_token` cookie (`SameSite=Strict`, `Secure`). This
works across `app.*` ↔ `api.*` **because they share the registrable domain
`boonmafarm.com` (same-site)** — no `SameSite=None` and no `Domain` attribute needed.
Two requirements, both already satisfied:

- `CORS_ALLOWED_ORIGINS=https://app.boonmafarm.com` (specific origin, not `*`, so the
  API returns `Access-Control-Allow-Credentials: true`). Add the dev origin
  comma-separated if `dev.boonmafarm.com` must call prod.
- Frontend served from a `*.boonmafarm.com` host (prod `app.*`, dev `dev.*`).

> ⚠️ The only place cookie auth breaks is a **raw Vercel preview URL**
> (`something.vercel.app`) — that's a different registrable domain (cross-site), so
> the `Strict` cookie isn't sent and third-party cookies are being phased out. Auth
> only works on the mapped `*.boonmafarm.com` domains. If you need login to work on
> ad-hoc preview deploys, that's the one case to add `Authorization: Bearer` (the
> login response already returns `AccessToken` in its body for exactly this).

---

## Security checklist

- [ ] 5432 bound to `127.0.0.1` only; not in `ufw`
- [ ] `.env` is `chmod 600`, never committed (gitignored)
- [ ] App connects as `app_user` (DML only) — not `postgres`, not `sys_admin`
- [ ] `password_encryption = scram-sha-256`; no `md5`/`trust` in `pg_hba.conf`
- [ ] `AUTHENTICATION_JWT_SECRET` is 32+ random chars
- [ ] `CORS_ALLOWED_ORIGINS` is an explicit origin, not `*`
- [ ] GHCR PAT on the droplet is `read:packages` only
- [ ] GitHub `prod` environment has required reviewers
- [ ] Daily backup cron verified (check the log the next morning)
- [ ] TLS cert issued (`curl -I https://api.boonmafarm.com/health` → 200)
```
