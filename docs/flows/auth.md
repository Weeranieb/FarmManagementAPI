# Auth flow

## Purpose

Authentication flow for registering users, logging in with username/password, and logging out. These endpoints are **public** (no JWT required). On successful login, the API returns a JWT and may set it in an HTTP-only cookie for browser clients.

## Actors / authorization

- **Register**: Caller provides `clientId`, `username`, `password`, `firstName`, and optional fields. No auth required.
- **Login**: Any caller with valid credentials. No auth required.
- **Logout**: Typically called by an authenticated user to clear the server-side session/cookie; may be called without a valid token.

## Endpoints

| Method | Path                    | Description                                          |
| ------ | ----------------------- | ---------------------------------------------------- |
| POST   | `/api/v1/auth/register` | Register a new user                                  |
| POST   | `/api/v1/auth/login`    | Login; returns JWT (and sets cookie when applicable) |
| POST   | `/api/v1/auth/logout`   | Logout; clear auth cookie                            |

Full request/response schemas: [../openapi.yaml](../openapi.yaml).

## Request / response

- **Register**: Body `RegisterRequest` — `clientId`, `username`, `password`, `firstName` (required); `lastName`, `userLevel`, `contactNumber` (optional). Response: `{ "result": true, "data": ... }`.
- **Login**: Body `LoginRequest` — `username`, `password` (required); `rememberMe` (optional). Response: `{ "result": true, "data": { "accessToken", "expiredAt", "user" } }` (or similar).
- **Logout**: No body. Response: `{ "result": true }`.

## Errors

| HTTP | Code (example) | Meaning                                  |
| ---- | -------------- | ---------------------------------------- |
| 400  | 500011         | Invalid request body                     |
| 401  | 500021         | Invalid username or password (login)     |
| 500  | 500020         | User already exists (register)           |
| 500  | 500022, 500023 | Invalid or expired token (if token used) |

Error response shape: `{ "code": "<string>", "message": "<string>" }`. See `internal/errors` for all auth error codes (500020–500029).
