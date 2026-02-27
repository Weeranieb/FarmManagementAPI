# Deploying Farm API (Go) on Vercel

The app is already set up for Vercel: it runs as a **Go serverless function** via `api/index.go` (package `handler`, exposes `Handler`) and `vercel.json`. The entrypoint only imports `src/app`, which in turn uses `internal` packages, so Vercel’s build does not hit Go’s “internal package not allowed” restriction.

## 0. Set up dev environment first

Use Vercel’s **Preview** and **Development** environments with your dev config so you can test before enabling Production.

1. In the dashboard go to **Settings → Environment Variables**.
2. Add the variables below. When adding each variable, select **Preview** and **Development** only (leave **Production** unchecked for now).
3. Use the same values you use locally (e.g. from `configuration/config.env` or `configuration/config.dev.yaml`). You can use **Import .env** and paste your dev env contents, then set the import to apply only to **Preview** and **Development**.
4. Deploy from a non-production branch (e.g. `develop` or `dev`) so the deployment uses the Preview environment and your dev env vars.
5. After everything works, add Production env vars (e.g. a separate DB or stronger `AUTHENTICATION_JWT_SECRET`) and enable **Production** for those variables.

**Variables to set for Dev/Preview:**

| Variable                    | Example (dev)       |
| --------------------------- | ------------------- |
| `DATABASE_HOST`             | your Neon host      |
| `DATABASE_PORT`             | `5432`              |
| `DATABASE_NAME`             | e.g. `neondb`       |
| `DATABASE_USER`             | your Neon user      |
| `DATABASE_PASSWORD`         | your Neon password  |
| `DATABASE_SSL_MODE`         | `require`           |
| `APP_ENVIRONMENT`           | `development`       |
| `APP_LOG_LEVEL`             | `info`              |
| `APP_DEBUG`                 | `true`              |
| `AUTHENTICATION_JWT_SECRET` | your dev JWT secret |
| `AUTHENTICATION_JWT_EXPIRY` | `24h`               |

**Security:** Do not commit real secrets. Keep `configuration/config.env` (and any file with passwords) out of the repo; use Vercel’s env vars (or **Import .env** in the UI without committing that file).

---

## 1. Vercel dashboard settings

In your Vercel project **Settings → General** (and **Build & Output**):

| Setting              | Value           | Why                                                                 |
| -------------------- | --------------- | ------------------------------------------------------------------- |
| **Root Directory**   | `backend`       | So Vercel uses the folder that contains `go.mod` and `vercel.json`. |
| **Framework Preset** | Other           | Not a Node/Next app.                                                |
| **Build Command**    | _(leave empty)_ | Build is driven by `vercel.json` and the `@vercel/go` builder.      |
| **Output Directory** | _(leave empty)_ | Go runs as serverless functions; there is no static output.         |
| **Install Command**  | _(leave empty)_ | Go dependencies are installed from `go.mod` by the Go runtime.      |

Override if Vercel shows defaults: clear **Build Command**, **Output Directory**, and **Install Command** so the UI does not run `npm run build` or `npm install`.

## 2. Environment variables

Your config uses Viper with `AutomaticEnv()` and `SetEnvKeyReplacer` (`.` → `_`). Set these in **Settings → Environment Variables** (at least for Production):

**Database (e.g. Neon):**

- `DATABASE_HOST` — DB host
- `DATABASE_PORT` — e.g. `5432`
- `DATABASE_NAME` — DB name
- `DATABASE_USER` — DB user
- `DATABASE_PASSWORD` — DB password
- `DATABASE_SSL_MODE` — e.g. `require` for Neon

**App & auth:**

- `APP_ENVIRONMENT` — e.g. `production`
- `AUTHENTICATION_JWT_SECRET` — strong secret for JWT signing
- `AUTHENTICATION_JWT_EXPIRY` — e.g. `24h`

**Optional (have defaults):**

- `SERVER_PORT` — usually not needed (Vercel invokes the function).
- `APP_LOG_LEVEL` — e.g. `info`
- `APP_DEBUG` — `false` in production

No config file is required on Vercel; env vars override defaults.

## 3. Deploy

- Connect the repo and set **Root Directory** to `backend`, then deploy, or
- From the repo root: `cd backend && vercel` (Vercel CLI will use `backend/vercel.json` and `go.mod`).

After deploy, all routes are rewritten to the Go handler (see `vercel.json`), so your API is available at the project URL (e.g. `https://<project>.vercel.app/api/v1/...`).

## 4. Viewing logs

Runtime logs (stdout/stderr from your Go function) appear in the Vercel dashboard:

1. Open your **project** on [vercel.com](https://vercel.com).
2. Go to the **Logs** tab (top navigation).
3. Trigger a request (e.g. open `https://<project>.vercel.app/health` or any API route).
4. In Logs, use **Live** or the time range filter to see recent output.

You should see:

- **`[Farm API] serverless cold start – building app`** and **`[Farm API] app ready`** on the first request (cold start).
- **Fiber request lines** (method, path, status, latency) from the logger middleware on each request.

Logs are kept for about 3 days. For longer retention, use [Log Drains](https://vercel.com/docs/drains).
