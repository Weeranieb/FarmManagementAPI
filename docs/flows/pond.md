# Pond flow

## Purpose

Ponds belong to a farm. This flow covers creating multiple ponds at once for a farm, listing ponds by farm, getting/updating/deleting a single pond. Newly created ponds start with status `maintenance`; status can be `active` or `maintenance`.

## Actors / authorization

- All pond endpoints require JWT. Access is client-scoped (user can only operate on ponds under farms of their client). Super admin can access any client’s data.

## Endpoints

| Method | Path                | Description                                |
| ------ | ------------------- | ------------------------------------------ |
| POST   | `/api/v1/pond`      | Create multiple ponds for a farm           |
| GET    | `/api/v1/pond`      | Get list of ponds by farm (query `farmId`) |
| GET    | `/api/v1/pond/{id}` | Get pond by ID                             |
| PUT    | `/api/v1/pond/{id}` | Update pond (farmId, name, status)         |
| DELETE | `/api/v1/pond/{id}` | Delete a pond                              |

Full request/response schemas: [../openapi.yaml](../openapi.yaml).

## Request / response

- **Create**: Body `CreatePondsRequest` — `farmId` (required), `names` (array of strings, min 1). Response: success with created ponds. New ponds have status `maintenance`.
- **List**: Query `farmId` (required). Response: `data` as array of `PondResponse`.
- **Get by ID**: Path `id`. Response: `data` as `PondResponse`.
- **Update**: Path `id`; body `UpdatePondBody` — `farmId`, `name`, `status` (optional; enum `active`, `maintenance`). Response: success with updated pond.
- **Delete**: Path `id`. Response: success without data.

## Errors

| HTTP | Code (example) | Meaning                            |
| ---- | -------------- | ---------------------------------- |
| 404  | 500070         | Pond not found                     |
| 500  | 500071, 500072 | Pond already exists, invalid input |

Error response shape: `{ "code": "<string>", "message": "<string>" }`. See `internal/errors` for pond codes (500070–500079).

## See also

- **Pond stock actions (fill / move / sell)** – [pond-stock-actions.md](pond-stock-actions.md). Fill, move, and sell flows and active pond (pond cycle) lifecycle; intended API design.
