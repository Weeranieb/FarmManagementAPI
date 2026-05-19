# Changelog — Farm OS API

All notable changes to the Farm OS Go backend API will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Linear issue references use the `FAR-N` format and link to https://linear.app/farm-os.

---

## [0.2.4](https://github.com/Weeranieb/FarmManagementAPI/compare/v0.2.3...v0.2.4) (2026-05-19)


### Features

* **daily-log:** collapse feed columns and improve API error responses ([#28](https://github.com/Weeranieb/FarmManagementAPI/issues/28)) ([89693b2](https://github.com/Weeranieb/FarmManagementAPI/commit/89693b29ad36fd98c662fe5e8ca424182b795405))

## [0.2.3](https://github.com/Weeranieb/FarmManagementAPI/compare/v0.2.2...v0.2.3) (2026-05-15)


### Features

* **pond:** add GET /pond/:pondId/activities endpoint ([#26](https://github.com/Weeranieb/FarmManagementAPI/issues/26)) ([5aa90e6](https://github.com/Weeranieb/FarmManagementAPI/commit/5aa90e6601b90ec3785e983171d7862eb38450ed))

## [0.2.2](https://github.com/Weeranieb/FarmManagementAPI/compare/v0.2.1...v0.2.2) (2026-05-14)


### Features

* **pond:** serve bulk-import Excel template from backend (FAR-76) ([#23](https://github.com/Weeranieb/FarmManagementAPI/issues/23)) ([c355378](https://github.com/Weeranieb/FarmManagementAPI/commit/c355378cb851cbd61a8799f6d31334c33a241ea7))

## [0.2.1](https://github.com/Weeranieb/FarmManagementAPI/compare/v0.2.0...v0.2.1) (2026-05-13)


### Features

* **pond:** refactor pond creation and update logic to include area field ([#18](https://github.com/Weeranieb/FarmManagementAPI/issues/18)) ([c42ef90](https://github.com/Weeranieb/FarmManagementAPI/commit/c42ef90310d539d4ab25ef3831393698bfe4752b))

## [0.2.0](https://github.com/Weeranieb/FarmManagementAPI/compare/v0.1.0...v0.2.0) (2026-05-09)


### Features

* add basic farm Group ([0457ea2](https://github.com/Weeranieb/FarmManagementAPI/commit/0457ea2d65de639388d4d08919b00788b2ced81d))
* add daily feed API and wire feed collection updates ([6a3873c](https://github.com/Weeranieb/FarmManagementAPI/commit/6a3873c06f0f622dcfb801c5340122c4d1f85780))
* add deleteDays field to DailyLogBulkUpsertRequest and update documentation ([72135dd](https://github.com/Weeranieb/FarmManagementAPI/commit/72135dd53804fb2ae0025ae2f7c6a5c0e93d8348))
* add endpoint for importing daily logs from Excel template ([7918f4f](https://github.com/Weeranieb/FarmManagementAPI/commit/7918f4f73434c714c8bc0081cb894996355c24e2))
* add farm hierarchy endpoint ([e249ee5](https://github.com/Weeranieb/FarmManagementAPI/commit/e249ee5ac80bbe5c9ffaeb6f1a23195bd831032a))
* add farm hierarchy endpoint ([69d26b9](https://github.com/Weeranieb/FarmManagementAPI/commit/69d26b907294df8e498279050e0a89ae0a4c6dd6))
* add isTouristFishingEnabled field to documentation and update pond service methods ([2a908f1](https://github.com/Weeranieb/FarmManagementAPI/commit/2a908f1fb893759c06d3c13b83e30088b1048bc2))
* add latest activity type to pond management ([7d34d5c](https://github.com/Weeranieb/FarmManagementAPI/commit/7d34d5c7c8d9e845d4c2b40f5aeeb5a1121ff16e))
* add PostgreSQL service and enhance login functionality ([a94c6a9](https://github.com/Weeranieb/FarmManagementAPI/commit/a94c6a99124162c5a4ea44e72fbadbe4a8b46f4f))
* add security middleware — CORS, helmet, rate limiter ([9b6a7cb](https://github.com/Weeranieb/FarmManagementAPI/commit/9b6a7cb28e2e2bbc8d265cd14ea763dae6462a3c))
* add worker ([8ecc90b](https://github.com/Weeranieb/FarmManagementAPI/commit/8ecc90be4d4286750c94d8d33419cae7cc187a93))
* enhance API security and improve request validation ([3d8af99](https://github.com/Weeranieb/FarmManagementAPI/commit/3d8af9979ff56dc55f8be88be1d91f237ab9e803))
* enhance client and daily log DTOs with optional fields ([3cb999c](https://github.com/Weeranieb/FarmManagementAPI/commit/3cb999cc35d7b5609467a3dd82daad550fe30dc7))
* enhance daily feed upload functionality ([5706e3e](https://github.com/Weeranieb/FarmManagementAPI/commit/5706e3e0bf8e0dd4b0f53ec6194694cbd0de86da))
* enhance daily log functionality with new fields and repository methods ([2f3720c](https://github.com/Weeranieb/FarmManagementAPI/commit/2f3720c6feddb9b8a09a2ebfdd7975f0c1f118d2))
* enhance farm management API with updated request structures and permissions ([c59b1cf](https://github.com/Weeranieb/FarmManagementAPI/commit/c59b1cfef68fd6515259f5a6afe1c479d40f6220))
* enhance farm management API with updated request structures and… ([a61074f](https://github.com/Weeranieb/FarmManagementAPI/commit/a61074f84bd9237e503bfe2bdeac0c2925190e3e))
* enhance farm response mapping and repository methods ([c02ac71](https://github.com/Weeranieb/FarmManagementAPI/commit/c02ac71967a3daf79e3ee4ad4dd1bbe0ab2d8d69))
* enhance password validation for user registration and management ([3702719](https://github.com/Weeranieb/FarmManagementAPI/commit/37027197532882615a8bc7a5b5f0cf987af9657f))
* enhance password validation for user registration and management ([342250c](https://github.com/Weeranieb/FarmManagementAPI/commit/342250c3ced3116d13aeb270659dbb790f379527))
* enhance pond management with additional fields and repository methods ([0bf1228](https://github.com/Weeranieb/FarmManagementAPI/commit/0bf12283bea557f78d60e7c4d4a959c9b41baa05))
* enhance pond management with fill and update functionalities ([6a39f59](https://github.com/Weeranieb/FarmManagementAPI/commit/6a39f5981ed5c337a241f47352c88e6932cace54))
* enhance pond management with fill and update functionalities ([823455f](https://github.com/Weeranieb/FarmManagementAPI/commit/823455f692b8519e1038b7576c2b08c526cfaa8d))
* enhance user and farm management with new fields and admin functionalities ([a10e4e4](https://github.com/Weeranieb/FarmManagementAPI/commit/a10e4e4d7f00a3fb0177d32ff1b1705cf040bf86))
* farm group CRUD, pond latest activity for move-in, unique farm-group joins ([bee5aea](https://github.com/Weeranieb/FarmManagementAPI/commit/bee5aea4396c682ef4a78f4c51004abe524815ba))
* implement Bearer authentication for all API endpoints ([c384797](https://github.com/Weeranieb/FarmManagementAPI/commit/c384797f9d2d207313d3a44a6e870c6443fc4079))
* implement bulk pond creation endpoint ([56bdef7](https://github.com/Weeranieb/FarmManagementAPI/commit/56bdef71e68e9e93bd8f1806b33ce94cd21a5b90))
* implement bulk pond creation endpoint ([8b4b6f6](https://github.com/Weeranieb/FarmManagementAPI/commit/8b4b6f6b20de174dab0b8a52e6dad1845ad6505e))
* implement get all client ([185b7d4](https://github.com/Weeranieb/FarmManagementAPI/commit/185b7d4bdf12201b8c3d2924f8310b7ca9418387))
* implement get all client ([2cd9ff6](https://github.com/Weeranieb/FarmManagementAPI/commit/2cd9ff6f48449a1f107663324e5eb7d17425bc4e))
* implement get method ([48ecc42](https://github.com/Weeranieb/FarmManagementAPI/commit/48ecc4260fdfa60ecede67e40b97d1e270bf0ba3))
* implement get method ([f1140d4](https://github.com/Weeranieb/FarmManagementAPI/commit/f1140d49769735406ece7e8484dfcc1994630114))
* implement get-farms-by-client ([9109a7d](https://github.com/Weeranieb/FarmManagementAPI/commit/9109a7dee5cf9dc657e6bd4bd691aebd6a7dbe04))
* implement get-farms-by-client ([b5c3622](https://github.com/Weeranieb/FarmManagementAPI/commit/b5c36225b33e224818f675de62a1f4ebe5af7819))
* implement sell page detail ([0424012](https://github.com/Weeranieb/FarmManagementAPI/commit/04240127ae993473790bd0bae867ea5d4f873c1a))
* implement user profile update functionality ([1f304e1](https://github.com/Weeranieb/FarmManagementAPI/commit/1f304e19188352404d9a2111d4fd70e71b924866))
* init api one line ([4a74ee6](https://github.com/Weeranieb/FarmManagementAPI/commit/4a74ee649761ef16f771495482ab615a7e1500a2))
* **merchant:** super-admin guards, update DTO, soft delete, air for run ([4f71201](https://github.com/Weeranieb/FarmManagementAPI/commit/4f71201bb1eac31966e0a687c3c91237204fa1d0))
* move pond - handler, service, cost utils, batch additional costs ([335447f](https://github.com/Weeranieb/FarmManagementAPI/commit/335447fa85da90137875277497f9f6a9ebc2cd88))
* only super can create user ([b41e52b](https://github.com/Weeranieb/FarmManagementAPI/commit/b41e52bf28b7664df173d78f285d9ec37c3f6734))
* **pond:** add start date to PondResponse and handle maintenance status ([537f961](https://github.com/Weeranieb/FarmManagementAPI/commit/537f9610a446bcfad488f2786d1618132236b4a9))
* **pond:** implement sell functionality for ponds ([eabc2ba](https://github.com/Weeranieb/FarmManagementAPI/commit/eabc2ba4efb9c280bee2154a93c4f4f38ae83ae1))
* refactor daily feed to daily log and enhance upload functionality ([3048467](https://github.com/Weeranieb/FarmManagementAPI/commit/304846744e1850f6dcb2841f4ad4ab688226cc4a))
* Review & Confirm preview APIs and cost refactor ([a15aef0](https://github.com/Weeranieb/FarmManagementAPI/commit/a15aef00b99db95a6c275d7fd33d9aedb869e5c9))
* update both auth and cookie ([5c9c215](https://github.com/Weeranieb/FarmManagementAPI/commit/5c9c215ee98d4e028d4c9afe949ef24d4df9b09b))
* update dependencies and enhance user registration and management ([9a86a09](https://github.com/Weeranieb/FarmManagementAPI/commit/9a86a0908c7066d50e741202af2b1a64cd57f96e))
* update farm and pond management API to use path parameters and refined request structures ([1dac7ce](https://github.com/Weeranieb/FarmManagementAPI/commit/1dac7ce713f120132d03109457638cead18c9b54))
* update middleware and create suer user ([30c6192](https://github.com/Weeranieb/FarmManagementAPI/commit/30c619214c018df383d06088082cca064f5da3dc))
* update readme ([e679864](https://github.com/Weeranieb/FarmManagementAPI/commit/e6798641ef669c3077e883d5dcee6553d9ce1f10))
* **user:** self change-password endpoint PUT /user/password ([1b52189](https://github.com/Weeranieb/FarmManagementAPI/commit/1b52189076c043e53a35c397bbfb7b3906abf3eb))
* **user:** self-update returns the updated user ([f75037a](https://github.com/Weeranieb/FarmManagementAPI/commit/f75037a13875aabcf4f5bd7adc0d89c4b3fac158))


### Bug Fixes

* **build:** add cgo build tag to all test files with mock imports ([ae47e9a](https://github.com/Weeranieb/FarmManagementAPI/commit/ae47e9aef376a6929441fa5e5dd9dc020507ddcc))
* **build:** add cgo build tag to sqlite test files for Vercel deploy ([9cff05e](https://github.com/Weeranieb/FarmManagementAPI/commit/9cff05ef3374dc7c61e6dba8e2ee214cf93bf7ca))
* **build:** commit mock files for Vercel go mod tidy ([56d2abd](https://github.com/Weeranieb/FarmManagementAPI/commit/56d2abd20bea372a07974bbab36bae2f6941e510))
* resolve all golangci-lint errors and make lint blocking ([18bc2ce](https://github.com/Weeranieb/FarmManagementAPI/commit/18bc2ce3d1435391b6fa14266c8afec5d1285db2))
* update farm and pond status to maintenance during creation ([fc0ffc7](https://github.com/Weeranieb/FarmManagementAPI/commit/fc0ffc7a0876da1d6b7afc42dace809bef50eb52))
* update Vercel rewrites and enhance request handling ([e562ff9](https://github.com/Weeranieb/FarmManagementAPI/commit/e562ff959f02584acfb6c05e937cabd947a831a9))
* **vercel:** correct rewrite destination from /api/index to /api ([029af50](https://github.com/Weeranieb/FarmManagementAPI/commit/029af5058b63c68d0515c4b78fd260271c2b650d))

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
