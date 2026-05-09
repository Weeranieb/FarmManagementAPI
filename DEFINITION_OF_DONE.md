# Definition of Done — Farm OS API

Before marking any Linear issue as **Done**, all applicable items must be checked.

This is the discipline a BA / SA / QA would normally enforce. As a solo dev, this is your guardrail against shipping half-done work.

---

## Code

- [ ] Branch follows naming: `niebfeelgood/far-<n>-<slug>` (Linear auto-generates this)
- [ ] PR title contains `FAR-<n>` (so Linear auto-links and updates status)
- [ ] PR description has 2-line summary + link to Linear issue
- [ ] `go vet ./...` clean
- [ ] `go test ./...` passes
- [ ] No `fmt.Println` debug statements left in
- [ ] No commented-out blocks of code
- [ ] Errors are wrapped with context (`fmt.Errorf("foo: %w", err)`), not swallowed

## Functional testing

- [ ] All acceptance criteria from the Linear issue are verified via curl / Postman
- [ ] Happy path tested
- [ ] At least 2 error paths tested (e.g., invalid input, unauthorized, not found)
- [ ] Tested with both owner role and staff role (if endpoint is role-scoped)

## Database

- [ ] If schema changed: migration file added under `migrations/dev/`
- [ ] Migration has both `.up.sql` and `.down.sql`
- [ ] Migration ran clean on a fresh DB (test with `make migrate-reset` or equivalent)
- [ ] No SQL injection risk (use parameterized queries)
- [ ] Indexes added for any new query pattern that filters/sorts

## API contract

- [ ] Swagger / OpenAPI spec updated to reflect new endpoint or field
- [ ] If breaking change: bumped major version, called out in CHANGELOG
- [ ] Web team aware (link in Linear or comment)
- [ ] Mobile team aware (link in Linear or comment)

## Security

- [ ] Auth required on all endpoints except explicit public ones
- [ ] Owner-scoped: user cannot read/write resources of other owners
- [ ] No secrets / DB credentials in code (use env vars)
- [ ] Passwords hashed with bcrypt (cost ≥ 10)

## Docs

- [ ] `CHANGELOG.md` updated under `[Unreleased]` with `(FAR-<n>)` tag
- [ ] If new migration: listed under `### Migrations` in CHANGELOG
- [ ] README updated if new env var or setup step added
- [ ] All acceptance criteria checkboxes ticked in the Linear issue description

---

## Definition of Released

Once a release ships:

- [ ] Issue moved to milestone "6. Released"
- [ ] Version bump committed (VERSION file or `go.mod` build tag)
- [ ] CHANGELOG `[Unreleased]` rolled into a new version section with date
- [ ] Migrations applied to dev environment
- [ ] Linear issue status set to **Done**

---

## Skipping items

If a check doesn't apply, write `N/A` next to it in the PR description with a one-line reason. Don't silently skip.
