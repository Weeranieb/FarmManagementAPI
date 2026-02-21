# Farm hierarchy flow

## Purpose

Farms belong to a client and contain ponds. This flow covers farm CRUD, listing farms by client, and retrieving the **farm hierarchy** (farms with nested ponds) for UI tree/table views. Super admin can create/update farms and may pass `clientId` when listing or fetching hierarchy; others are scoped to their client.

## Actors / authorization

- **Super admin**: Can add farms (`POST /farm`), update any farm (`PUT /farm/{id}`), list farms with optional `clientId`, and get hierarchy with optional `clientId`.
- **Client admin / normal user**: Can list farms for their client, get farm by ID (within client), and get hierarchy for their client. Cannot create farms or update farms (unless allowed by handler).

Access is client-scoped; super admin can operate on any client.

## Endpoints

| Method | Path                     | Description                                                   |
| ------ | ------------------------ | ------------------------------------------------------------- |
| POST   | `/api/v1/farm`           | Add a new farm (super admin only)                             |
| GET    | `/api/v1/farm`           | Get list of farms (optional query `clientId`)                 |
| GET    | `/api/v1/farm/hierarchy` | Get farms with nested ponds (optional `clientId`)             |
| GET    | `/api/v1/farm/{id}`      | Get farm by ID (detail + summary + ponds)                     |
| PUT    | `/api/v1/farm/{id}`      | Update farm (super admin only; name only; clientId preserved) |

Full request/response schemas: [../openapi.yaml](../openapi.yaml).

## Request / response

- **Create**: Body `CreateFarmRequest` — `clientId`, `name` (required). Response: success with created farm.
- **List**: Query `clientId` (optional). Response: `data` as list of farms with total counts.
- **Hierarchy**: Query `clientId` (optional). Response: `data` as array of `FarmHierarchyItem` (farm + ponds).
- **Get by ID**: Path `id`. Response: `data` as `FarmDetailResponse` (id, clientId, name, status, summary, ponds).
- **Update**: Path `id`; body `UpdateFarmBody` — `name`. Response: success with updated farm.

## Errors

| HTTP | Code (example)         | Meaning                                       |
| ---- | ---------------------- | --------------------------------------------- |
| 403  | 500024                 | Permission denied                             |
| 500  | 500040, 500041, 500042 | Farm not found, already exists, invalid input |

Error response shape: `{ "code": "<string>", "message": "<string>" }`. See `internal/errors` for farm codes (500040–500049).
