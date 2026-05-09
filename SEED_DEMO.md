# Demo Data Seeding — Farm OS API

Pattern for resetting demo data so you can demo without fear of corrupting it.

> **Why this exists**: solving the original pain — "ไม่กล้า demo เพราะกลัวข้อมูลพัง". With this in place, you can let the demo audience click anything; one command resets state for the next demo.

---

## Goals

1. **Idempotent** — running the seed twice produces the same result
2. **Fast** — under 5 seconds to wipe + reseed
3. **Safe** — refuses to run against production
4. **Self-contained** — no external API calls, all data is local

---

## Recommended structure

```
backend/
├── cmd/
│   └── seed/
│       └── main.go              # entry point: `go run ./cmd/seed -env=demo`
├── internal/
│   └── seed/
│       ├── seed.go              # orchestrator
│       ├── users.go             # demo users (owner + staff)
│       ├── farms.go             # 1 farm + 3 ponds
│       ├── feed_logs.go         # ~14 days of feed logs
│       ├── movements.go         # 1 stocking + 1 transfer per pond
│       └── transactions.go      # ~5 buys + ~3 sells
└── migrations/
    └── dev/                     # existing migrations
```

---

## Safety guard (CRITICAL)

Always check the env before running. Refuse if it's prod.

```go
// cmd/seed/main.go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "strings"
)

func main() {
    env := flag.String("env", "", "Environment: dev | demo")
    flag.Parse()

    if *env == "" {
        log.Fatal("must specify -env=dev or -env=demo")
    }

    if strings.Contains(strings.ToLower(*env), "prod") {
        log.Fatal("REFUSED: never seed against prod")
    }

    // Belt and braces: also check DB host
    dbHost := os.Getenv("DB_HOST")
    if strings.Contains(dbHost, "prod") {
        log.Fatalf("REFUSED: DB_HOST=%q looks like prod", dbHost)
    }

    fmt.Printf("Seeding %s environment...\n", *env)
    // ... call seed.Run(ctx, db)
}
```

---

## Reset + reseed pattern

```go
// internal/seed/seed.go
package seed

import (
    "context"
    "database/sql"
)

// Run wipes all demo-tagged data, then re-inserts a known fixture set.
// Tables are truncated in dependency order (children first).
func Run(ctx context.Context, db *sql.DB) error {
    return withTx(ctx, db, func(tx *sql.Tx) error {
        if err := wipe(ctx, tx); err != nil {
            return err
        }
        if err := seedUsers(ctx, tx); err != nil {
            return err
        }
        if err := seedFarms(ctx, tx); err != nil {
            return err
        }
        if err := seedFeedLogs(ctx, tx); err != nil {
            return err
        }
        if err := seedMovements(ctx, tx); err != nil {
            return err
        }
        return seedTransactions(ctx, tx)
    })
}

func wipe(ctx context.Context, tx *sql.Tx) error {
    // Order matters — children first
    tables := []string{
        "transactions",
        "transfers",
        "stockings",
        "feed_logs",
        "ponds",
        "farms",
        "user_roles",
        "users",
    }
    for _, t := range tables {
        if _, err := tx.ExecContext(ctx, "TRUNCATE TABLE "+t+" RESTART IDENTITY CASCADE"); err != nil {
            return err
        }
    }
    return nil
}
```

---

## Demo fixture data (suggested)

### Users

| Email | Password | Role | Notes |
|---|---|---|---|
| `owner@demo.farm` | `demo1234` | owner | demo account, full access |
| `staff@demo.farm` | `demo1234` | staff | demo account, mobile only |

### Farms / Ponds

- **Farm**: "ฟาร์มสาธิตปลานิล" (Demo Tilapia Farm)
- **Pond A**: 200 m² · current count 800 (after stock + transfer)
- **Pond B**: 150 m² · current count 200 (after transfer in)
- **Pond C**: 300 m² · current count 1500 (untouched after stocking)

### Feed logs

- 14 days of daily logs per pond, ~5 kg/day, varying feed types

### Movements

- Pond A: 1 stocking of 1000 fish + 1 transfer of 200 to Pond B
- Pond B: 1 transfer in of 200 from Pond A
- Pond C: 1 stocking of 1500 fish

### Transactions

- 5 buys: fingerlings × 2, feed × 2, supplies × 1
- 3 sells: small / medium harvest

---

## Makefile targets

```makefile
.PHONY: seed-demo seed-dev migrate-reset

migrate-reset:
	@echo "Dropping and re-running migrations..."
	migrate -path migrations/dev -database "$$DATABASE_URL" drop -f
	migrate -path migrations/dev -database "$$DATABASE_URL" up

seed-dev: migrate-reset
	go run ./cmd/seed -env=dev

seed-demo: migrate-reset
	go run ./cmd/seed -env=demo
	@echo "✅ Demo data ready. Login: owner@demo.farm / demo1234"
```

---

## Pre-demo checklist

10 minutes before a demo:

1. `make seed-demo`
2. Visit `/login` on web with `owner@demo.farm` → confirm dashboard loads
3. Open mobile app, login as `staff@demo.farm` → confirm feed log entry works
4. Take a screenshot of the home dashboard for backup if live demo fails

---

## When things go wrong mid-demo

If demo data gets corrupted by an audience member typing wild input:

1. Apologize briefly, say "let me reset to fresh data"
2. Open terminal: `make seed-demo`
3. Refresh browser / restart mobile app
4. Continue — total time ~30 seconds

This is **why** the seed exists. Use it without shame.
