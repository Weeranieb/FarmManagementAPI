# Client flow

## Purpose

Client (tenant/organization) CRUD. Clients are top-level entities; users and farms belong to a client. Used for multi-tenant scoping: normal users and client admins are restricted to their own client; super admins can manage all clients.

## Actors / authorization

- **Super admin (userLevel 3)**: Can create clients (`POST /client`), get client list for dropdown (`GET /client/list`), get any client by ID, and update any client.
- **Client admin / normal user**: Can get and update only their own client (enforced via `clientId` from JWT). Cannot create clients or access the global client list.

Access is enforced using `utils.CanAccessClient(c.UserContext(), targetClientId)`; super admin bypasses the check.

## Endpoints

| Method | Path                  | Description                                     |
| ------ | --------------------- | ----------------------------------------------- |
| POST   | `/api/v1/client`      | Add a new client (super admin only)             |
| PUT    | `/api/v1/client`      | Update client (own client or super admin)       |
| GET    | `/api/v1/client/list` | Get client list for dropdown (super admin only) |
| GET    | `/api/v1/client/{id}` | Get client by ID                                |

Full request/response schemas: [../openapi.yaml](../openapi.yaml).

## Request / response

- **Create**: Body `CreateClientRequest` — `name`, `ownerName`, `contactNumber` (required). Response: success with created client data.
- **Update**: Body `UpdateClientRequest` — `id` (required); `name`, `ownerName`, `contactNumber`, `isActive` (optional). Response: success with updated client.
- **List**: No body. Response: `data` as array of `DropdownItem` (key = id, value = name).
- **Get by ID**: Path `id`. Response: `data` as `ClientResponse`.

## Errors

| HTTP | Code (example) | Meaning                                                                                   |
| ---- | -------------- | ----------------------------------------------------------------------------------------- |
| 403  | 500024         | Permission denied (e.g. non–super admin calling create/list, or accessing another client) |
| 404  | 500110         | Client not found                                                                          |
| 500  | 500111, 500112 | Client already exists, invalid input                                                      |

Error response shape: `{ "code": "<string>", "message": "<string>" }`. See `internal/errors` for client codes (500110–500119).
