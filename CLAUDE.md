# Project conventions

## Branch naming

Feature branches use **`<type>/FAR-XX-<kebab-description>`** — Conventional-Commit type prefix, slash, Linear ticket, kebab description.

Examples:

- `feat/FAR-76-download-template-for-add-multiple-ponds`
- `fix/FAR-123-pond-area-validation`
- `refactor/FAR-200-extract-auth-middleware`

Non-ticket work omits the FAR-XX: `chore/release-please`, `ci/migration-workflows`.

### Do NOT use

- **Linear's auto-suggested `gitBranchName`** (e.g. `niebfeelgood/far-76-...`) — wrong case, wrong prefix.
- **Bare `FAR-XX-*`** without a type prefix — older branches like `FAR-9-login-screen` predate the rule and are not the current convention.

PR titles follow Conventional Commits (linted by CI).

## Do not auto-commit

Never run `git commit`, `git push`, or open a PR automatically just because a task is implemented and verified. Finish coding, summarize what changed, and stop.

Phrases like "do it", "do what you plan", "proceed", or "go ahead" mean **implement** — not commit. Only commit when the user explicitly says commit / push / ship / open PR.

## Regenerate Swagger after adding/changing an API

When you add a new HTTP endpoint, change a route, or modify request/response DTOs or `@Swagger` annotations, regenerate the OpenAPI spec as a post-edit step:

```
make gen-swag
```

This runs `swag init -g src/cmd/api/main.go -o docs --propertyStrategy snakecase` and updates `docs/`. Commit the regenerated `docs/` files alongside the code change so `/swagger/*` and downstream clients stay in sync.

Skip this only for changes that don't affect the public API surface (internal refactors, repository/service-only changes, comment fixes).

## Go utilities

- Prefer **[`github.com/samber/lo`](https://github.com/samber/lo)** for generic helpers (e.g. `lo.ToPtr`, `lo.FromPtr`, `lo.EmptyableToPtr`, `lo.Map`, `lo.Filter`) instead of rolling your own `toPtr` / `fromPtr` or similar one-off algorithms.
- If a helper is used in **multiple packages** (handler, service, repository, etc.), put it in **`src/internal/utils/`** — not duplicated inline or in a single feature package.
