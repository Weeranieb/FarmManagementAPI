# User flow

## Purpose

User account CRUD and listing within a client. Users belong to a client; creation requires `clientId` and `userLevel`. Used for managing who can log in and with what scope (normal user, client admin, super admin).

## Actors / authorization

- **Super admin**: Can create users for any client, get/update users across clients (within API design).
- **Client admin**: Typically can manage users for their own client only.
- **Normal user**: May only get/update their own profile depending on handler logic.

All user endpoints are behind JWT; access is client-scoped via `utils.CanAccessClient(c.UserContext(), targetClientId)` where applicable.

## Endpoints

| Method | Path                | Description                       |
| ------ | ------------------- | --------------------------------- |
| POST   | `/api/v1/user`      | Add a new user                    |
| GET    | `/api/v1/user`      | Get current user (from JWT)       |
| PUT    | `/api/v1/user`      | Update user                       |
| GET    | `/api/v1/user/list` | Get user list (e.g. for dropdown) |

Full request/response schemas: [../openapi.yaml](../openapi.yaml).

## Request / response

- **Create**: Body `CreateUserRequest` — `username`, `password`, `firstName` (required); `lastName`, `userLevel`, `contactNumber`, `clientId` (optional). Response: success with created user.
- **Get current**: No body. Response: `data` as `UserResponse` for the authenticated user.
- **Update**: Body `UpdateUserRequest` — `username`, `firstName`, `lastName`, `userLevel`, `contactNumber`. Response: success with updated user.
- **List**: No body. Response: `data` as list of users (e.g. for dropdown).

## Errors

| HTTP | Code (example) | Meaning                             |
| ---- | -------------- | ----------------------------------- |
| 400  | 500031         | Invalid user input                  |
| 403  | 500024         | Permission denied                   |
| 500  | 500030, 500032 | User not found, user already exists |

Error response shape: `{ "code": "<string>", "message": "<string>" }`. See `internal/errors` for user codes (500030–500039).
