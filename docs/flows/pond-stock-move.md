# Pond stock action: Move (transfer fish)

## Purpose

Transfer fish from one pond (source) to another (destination). The path uses the **source** `pondId`. If the **destination** pond is in **maintenance**, the backend creates an active pond for it and records the move with both source and destination active ponds.

## Actors / authorization

- All pond stock action endpoints require JWT. Access is client-scoped. Super admin can access any client’s data.

## Endpoints

| Method | Path                         | Description                                                                          |
| ------ | ---------------------------- | ------------------------------------------------------------------------------------ |
| POST   | `/api/v1/pond/{pondId}/move` | Move fish from this pond to another. Path = source pondId; body includes `toPondId`. |

Full request/response schemas: [../openapi.yaml](../openapi.yaml) (tag `pond-stock-actions`, schema `PondMoveRequest`).

## Request / response

- **Path**: `pondId` = source pond.
- **Body** `PondMoveRequest`: `toPondId`, `fishType`, `amount`, `activityDate` (required); `fishWeight` (optional). Backend resolves or creates destination active_pond.
- **Response**: Success with created move activity. Standard `{ "result": true, "data": ... }`.

## Behavior

- Resolve source pond’s active cycle (`active_pond_id`). If none (source in maintenance), return 400/404 as appropriate.
- Resolve destination pond’s active cycle. If destination has no active cycle (maintenance), create a new `active_ponds` row for the destination and optionally set destination pond status to `active`.
- Create activity with `mode = move`, `active_pond_id` = source, `to_active_pond_id` = destination (and other fields from body).

## Errors

| HTTP | Meaning                                                                                                                            |
| ---- | ---------------------------------------------------------------------------------------------------------------------------------- |
| 400  | Validation failed. **Business**: source pond not yet active (maintenance) — move requires the source pond to have an active cycle. |
| 404  | Pond not found (source or destination).                                                                                            |
| 500  | Internal/server error.                                                                                                             |

Error response shape: `{ "code": "<string>", "message": "<string>" }`.

## See also

- [pond-stock-actions.md](pond-stock-actions.md) – Overview and active pond concept.
- [pond-stock-fill.md](pond-stock-fill.md), [pond-stock-sell.md](pond-stock-sell.md) – Other modes.
