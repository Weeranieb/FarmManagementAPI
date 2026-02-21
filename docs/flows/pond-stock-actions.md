# Pond stock actions (overview)

## Purpose

The pond list page has an **Add** button that opens a modal with three actions: **fill** (add fish / ลงปลา), **move** (transfer fish / ย้ายปลา), and **sell** (ขายปลา). These actions are bound to the concept of **active pond (pond cycle)**. This document gives an overview; each action is described in its own flow.

- [Fill](pond-stock-fill.md) – Add fish to a pond; creates active pond if pond is in maintenance.
- [Move](pond-stock-move.md) – Transfer fish from one pond to another; destination may become active.
- [Sell](pond-stock-sell.md) – Record a sell; optionally close the cycle and return pond to maintenance.

## Concept: Active pond (pond cycle)

- **Active pond** = one “cycle” of a pond. Stored in `active_ponds`: `pond_id`, `start_date`, `end_date`, `is_active`. One pond can have many cycles over time; at most one row per pond has `is_active = true`.
- **Activities** (fill, move, sell) are tied to an **active pond** (`active_pond_id`); move also uses `to_active_pond_id` for the destination cycle.
- **Constraint**: Only one active cycle per pond at a time. To start a new cycle (e.g. first fill after maintenance), the current active cycle for that pond must be closed first (e.g. by sell with “mark to close”, or an explicit close step).
- **API design**: Paths use **pondId** (e.g. `POST /api/v1/pond/{pondId}/fill`) so the frontend only sends what it has; the backend resolves or creates the active pond as needed.

## When is an active pond created?

- **First fill** on a pond in **maintenance**: `POST /api/v1/pond/{pondId}/fill` → backend creates a new active pond for that pond and records the fill activity.
- **Move** into a pond that is in **maintenance**: the **destination** pond gets an active pond created (and becomes active); the move activity is recorded with source and destination active ponds.

## Sell and return to maintenance

When sell is performed with **markToClose**, after the transaction the active pond is closed and the pond returns to **maintenance**.

## Business errors (summary)

- **Fill**: `fishType` must be in the allowed list (e.g. fish type constants). See [pond-stock-fill.md](pond-stock-fill.md#errors).
- **Move**: Source pond must already have an active cycle (cannot move from a pond in maintenance). See [pond-stock-move.md](pond-stock-move.md#errors).
- **Sell**: Pond must have an active cycle (cannot sell from a pond in maintenance). See [pond-stock-sell.md](pond-stock-sell.md#errors).

## Implementation status

These endpoints are **not yet implemented** in the backend; the flow docs and [../openapi.yaml](../openapi.yaml) describe the intended API design. The frontend pond list Add button and StockActionModal already support the three modes.

## See also

- [pond.md](pond.md) – Basic pond CRUD and status (active/maintenance).
- [../openapi.yaml](../openapi.yaml) – Paths under `pond-stock-actions` tag.
