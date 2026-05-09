# Changelog — Farm OS API

All notable changes to the Farm OS Go backend API will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Linear issue references use the `FAR-N` format and link to https://linear.app/farm-os.

---

## [Unreleased]

### Added
- _work in progress_

### Changed
- _none_

### Fixed
- _none_

### Removed
- _none_

### Migrations
- _none_

---

## [0.1.0] — 2026-06-08 (Demo Release)

First public-ish API release. Powers web + mobile demo flow.

### Added
- `POST /auth/signup`, `POST /auth/login`, `POST /auth/logout` (FAR-6)
- `POST /auth/refresh` + auth middleware with role check (owner / staff) (FAR-7)
- `POST/GET/PUT/DELETE /farms` — owner-scoped (FAR-13)
- `POST/GET/PUT /farms/:id/ponds` — returns `current_fish_count` (FAR-14)
- `POST /ponds/:id/feed-logs`, `GET` history with date range, `PUT`/`DELETE` within 24h (FAR-21)
- `POST /ponds/:id/stockings`, `POST /transfers`, `GET` movement history per pond — pond count auto-updates atomically, idempotent on retry (FAR-27)
- `POST /transactions/buy`, `POST /transactions/sell`, `GET /transactions/summary` (FAR-33)

### Migrations

| Migration | Notes |
|---|---|
| `2026XXXXXXXX_add_users_and_roles` | Adds `users`, `user_roles` tables |
| `2026XXXXXXXX_add_farms_ponds` | Adds `farms`, `ponds` tables with `current_fish_count` |
| `2026XXXXXXXX_add_feed_logs` | Adds `feed_logs` table |
| `2026XXXXXXXX_add_movements` | Adds `stockings`, `transfers` tables + trigger to update pond count |
| `2026XXXXXXXX_add_transactions` | Adds `transactions` (buy/sell) table |

### Notes
- Tested against PostgreSQL 15
- Mobile dependency: app must support `Authorization: Bearer <jwt>`

---

## How to update this file

1. While working on a Linear issue, append your change under `[Unreleased]` in the appropriate category.
2. Each entry: `- Short description (FAR-N)`.
3. If your change adds/modifies a DB migration, also list it under `### Migrations` with notes.
4. On release day, move `[Unreleased]` entries under a new version heading with the date.
5. Bump version in `go.mod` or VERSION file accordingly.

## When to bump major / minor / patch

- **Major (X.0.0)** — breaking API changes (path/payload/response shape that breaks web or mobile)
- **Minor (0.X.0)** — backward-compatible new endpoints or fields
- **Patch (0.0.X)** — bug fixes, performance, internal refactor
