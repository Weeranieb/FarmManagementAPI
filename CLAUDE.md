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
