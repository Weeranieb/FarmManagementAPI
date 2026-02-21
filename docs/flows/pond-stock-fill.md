# Pond stock action: Fill (add fish)

## Purpose

Add fish to a pond. If the pond is in **maintenance** (no current active cycle), the backend **creates** a new active pond for that pond and records the fill activity. The frontend only sends `pondId` in the path; no need to send `activePondId`.

## Actors / authorization

- All pond stock action endpoints require JWT. Access is client-scoped (ponds under the user’s client). Super admin can access any client’s data.

## Endpoints

| Method | Path                         | Description                                                                        |
| ------ | ---------------------------- | ---------------------------------------------------------------------------------- |
| POST   | `/api/v1/pond/{pondId}/fill` | Add fish to the pond; backend resolves or creates active_pond and stamps activity. |

Full request/response schemas: [../openapi.yaml](../openapi.yaml) (tag `pond-stock-actions`, schema `PondFillRequest`).

## Request / response

- **Body** `PondFillRequest`: `fishType`, `amount`, `activityDate` (required); `fishWeight`, `fishUnit`, `pricePerUnit` (optional).
- **Response**: Success with created activity (and active_pond if one was created). Standard `{ "result": true, "data": ... }`.

## Behavior

- If pond has an active cycle: use that `active_pond_id` and create an activity with `mode = fill`.
- If pond is in maintenance (no active cycle): create a new row in `active_ponds` for this pond (`is_active = true`, `start_date` from activity or today), then create the fill activity. The pond may be updated to status `active` depending on business rules.

## Errors

| HTTP | Meaning                                                                                                                                      |
| ---- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| 400  | Validation failed (e.g. invalid amount). **Business**: `fishType` not in the allowed list (e.g. not one of the defined fish type constants). |
| 404  | Pond not found.                                                                                                                              |
| 500  | Internal/server error.                                                                                                                       |

Error response shape: `{ "code": "<string>", "message": "<string>" }`.

## See also

- [pond-stock-actions.md](pond-stock-actions.md) – Overview and active pond concept.
- [pond-stock-move.md](pond-stock-move.md), [pond-stock-sell.md](pond-stock-sell.md) – Other modes.
