# Pond stock action: Sell

## Purpose

Record a sell transaction from a pond. If **markToClose** is true, after the transaction the **active pond** is closed (`is_active = false`, `end_date` set) and the **pond** returns to status **maintenance**. This enforces “only one active cycle per pond at a time”: closing the cycle frees the pond for a new cycle later (e.g. after next fill).

## Actors / authorization

- All pond stock action endpoints require JWT. Access is client-scoped. Super admin can access any client’s data.

## Endpoints

| Method | Path                         | Description                                                            |
| ------ | ---------------------------- | ---------------------------------------------------------------------- |
| POST   | `/api/v1/pond/{pondId}/sell` | Record sell; optionally close active pond and set pond to maintenance. |

Full request/response schemas: [../openapi.yaml](../openapi.yaml) (tag `pond-stock-actions`, schema `PondSellRequest`, `PondSellDetailItem`).

## Request / response

- **Path**: `pondId` = pond to sell from.
- **Body** `PondSellRequest`: `activityDate` (required); `details` (array of per-species lines: fishType, size, amount, fishUnit, pricePerUnit), `buyer`, `markToClose` (optional). If `markToClose` is true, close the active cycle and set pond to maintenance after the transaction.
- **Response**: Success with created sell activity (and sell_details). Standard `{ "result": true, "data": ... }`.

## Behavior

- Resolve pond’s active cycle. If none, return 400/404.
- Create activity with `mode = sell` and related `sell_details` rows from `details`.
- If `markToClose`: update the active_pond row to `is_active = false`, set `end_date`; update pond status to `maintenance`.

## Errors

| HTTP | Meaning                                                                                                              |
| ---- | -------------------------------------------------------------------------------------------------------------------- |
| 400  | Validation failed. **Business**: pond not yet active (maintenance) — sell requires the pond to have an active cycle. |
| 404  | Pond not found.                                                                                                      |
| 500  | Internal/server error.                                                                                               |

Error response shape: `{ "code": "<string>", "message": "<string>" }`.

## See also

- [pond-stock-actions.md](pond-stock-actions.md) – Overview and active pond concept.
- [pond-stock-fill.md](pond-stock-fill.md), [pond-stock-move.md](pond-stock-move.md) – Other modes.
