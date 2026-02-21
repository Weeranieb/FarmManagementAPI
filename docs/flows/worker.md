# Worker flow

## Purpose

Workers are associated with a farm group and a client. This flow covers worker CRUD and listing for the current client (e.g. for assignment to tasks or ponds). Fields include `farmGroupId`, `hireDate`, `salary`, and `nationality`.

## Actors / authorization

- All worker endpoints require JWT. Access is client-scoped; users see and manage only workers belonging to their client. Super admin can access any client’s workers.

## Endpoints

| Method | Path                  | Description              |
| ------ | --------------------- | ------------------------ |
| POST   | `/api/v1/worker`      | Add a new worker         |
| GET    | `/api/v1/worker`      | List workers (paginated) |
| GET    | `/api/v1/worker/{id}` | Get worker by ID         |
| PUT    | `/api/v1/worker`      | Update a worker          |

Full request/response schemas: [../openapi.yaml](../openapi.yaml).

## Request / response

- **Create**: Body `CreateWorkerRequest` — `farmGroupId`, `firstName`, `nationality`, `salary` (required); `lastName`, `contactNumber`, `hireDate` (optional). Response: success with created worker.
- **List**: No body (pagination may be via query params if implemented). Response: `data` with items and total (e.g. `PageResponse`).
- **Get by ID**: Path `id`. Response: `data` as `WorkerResponse`.
- **Update**: Body `UpdateWorkerRequest` — `id` (required); `farmGroupId`, `firstName`, `lastName`, `contactNumber`, `nationality`, `salary`, `hireDate`, `isActive`. Response: success with updated worker.

## Errors

| HTTP | Code (example) | Meaning                              |
| ---- | -------------- | ------------------------------------ |
| 404  | 500080         | Worker not found                     |
| 500  | 500081, 500082 | Worker already exists, invalid input |

Error response shape: `{ "code": "<string>", "message": "<string>" }`. See `internal/errors` for worker codes (500080–500089).
