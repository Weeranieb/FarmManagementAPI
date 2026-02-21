# Merchant flow

## Purpose

Merchants are entities the client does business with (e.g. feed or equipment suppliers). The merchant list is **global** (shared across all clients). This flow covers merchant CRUD and listing.

## Actors / authorization

- All merchant endpoints require JWT.
- **Add (create)**: **Only super admin** can add merchants. This keeps the global list managed centrally.
- List, get by ID, and update: available to authenticated users (so they can select merchants e.g. in activities). Update may be restricted to super admin depending on implementation.

## Endpoints

| Method | Path                    | Description                               |
| ------ | ----------------------- | ----------------------------------------- |
| POST   | `/api/v1/merchant`      | Add a new merchant (**super admin only**) |
| GET    | `/api/v1/merchant`      | Get list of merchants                     |
| GET    | `/api/v1/merchant/{id}` | Get merchant by ID                        |
| PUT    | `/api/v1/merchant`      | Update a merchant                         |

Full request/response schemas: [../openapi.yaml](../openapi.yaml).

## Request / response

- **Create**: Body `CreateMerchantRequest` — `name` (required); `contactNumber`, `location` (optional). Response: success with created merchant.
- **List**: No body. Response: `data` as array of `MerchantResponse`.
- **Get by ID**: Path `id`. Response: `data` as `MerchantResponse`.
- **Update**: Body `UpdateMerchantRequest` — `id` (required); `name`, `contactNumber`, `location`. Response: success with updated merchant.

## Errors

| HTTP | Code (example) | Meaning                                                                |
| ---- | -------------- | ---------------------------------------------------------------------- |
| 403  | 500024         | Permission denied (e.g. non–super admin calling POST to add merchant). |
| 404  | 500060         | Merchant not found                                                     |
| 500  | 500061, 500062 | Merchant already exists, invalid input                                 |

Error response shape: `{ "code": "<string>", "message": "<string>" }`. See `internal/errors` for merchant codes (500060–500069).
