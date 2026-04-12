# Deploying the backend on Vercel

Use [`.env.example`](.env.example) as a checklist of variable names to copy into Vercel (values differ per environment).

This API is deployed as a **single Go serverless function** (`api/index.go` → `src/app`) with rewrites so all routes hit that function. Follow these steps so builds succeed.

## 1. Root directory (required for monorepos)

If the Git repository root contains both `frontend/` and `backend/`, the Vercel project **must** use this backend folder as its root.

1. Open [Vercel Dashboard](https://vercel.com/dashboard) → your **backend** project.
2. **Settings** → **General** → **Root Directory**.
3. Set to **`backend`** (or the path to this folder from repo root).
4. Save and redeploy.

If Root Directory is wrong, the build often fails at `go mod tidy` because `go.mod` is not found.

## 2. Go version

`go.mod` uses `go 1.24` (no patch) so Vercel’s `@vercel/go` builder can select a supported toolchain. Do not pin an arbitrary patch (e.g. `1.24.9`) unless you have confirmed Vercel can download that exact release.

## 3. Environment variables

Set these under **Settings** → **Environment Variables** for Production (and Preview if needed). Viper maps nested config with `_` (e.g. `database.host` → `DATABASE_HOST`).

| Variable | Notes |
|----------|--------|
| `DATABASE_HOST` | Postgres host |
| `DATABASE_PORT` | e.g. `5432` |
| `DATABASE_NAME` | Database name |
| `DATABASE_USER` | Database user |
| `DATABASE_PASSWORD` | Database password |
| `DATABASE_SSL_MODE` | e.g. `require` for hosted Postgres |
| `AUTHENTICATION_JWT_SECRET` | **Required** — empty value causes startup `log.Fatal` |
| `SERVER_PORT` | Optional; default `8080` |
| `APP_ENVIRONMENT` | e.g. `production` |

Optional app paths (if you use uploads on serverless, prefer external storage):

- `APP_DAILY_LOG_UPLOAD_PATH`
- `APP_DAILY_FEED_UPLOAD_PATH` (legacy alias)

## 4. Redeploy

After changing Root Directory, `go.mod`, or env vars, trigger a new deployment from the Vercel UI or by pushing to your connected branch.
