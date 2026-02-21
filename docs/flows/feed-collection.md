# Feed collection flow

## Purpose

Feed collections represent types of feed (e.g. by name and unit) and can have associated price history. This flow covers feed collection CRUD and **paginated listing** with optional ordering and keyword search. Used for managing feed catalog and linking to feed price history.

## Actors / authorization

- All feed-collection endpoints require JWT. Access is client-scoped; users see only feed collections for their client. Super admin can access any client’s data.

## Endpoints

| Method | Path                           | Description                                                            |
| ------ | ------------------------------ | ---------------------------------------------------------------------- |
| POST   | `/api/v1/feed-collection`      | Add a new feed collection (optionally with initial feedPriceHistories) |
| GET    | `/api/v1/feed-collection`      | Get paginated list (query `page`, `pageSize`, `orderBy`, `keyword`)    |
| GET    | `/api/v1/feed-collection/{id}` | Get feed collection by ID                                              |
| PUT    | `/api/v1/feed-collection`      | Update a feed collection                                               |

Full request/response schemas: [../openapi.yaml](../openapi.yaml).

## Request / response

- **Create**: Body `CreateFeedCollectionRequest` — `name`, `unit` (required); `feedPriceHistories` (optional array of `{ price, priceUpdatedDate }`). Response: success with created feed collection and any created price history.
- **List**: Query `page`, `pageSize` (required); `orderBy`, `keyword` (optional). Response: `data` as paginated list (e.g. `FeedCollectionPageResponse` with latest price info).
- **Get by ID**: Path `id`. Response: `data` as `FeedCollectionResponse`.
- **Update**: Body `UpdateFeedCollectionRequest` — `id` (required); `name`, `unit`. Response: success with updated feed collection.

## Errors

| HTTP | Code (example) | Meaning                                       |
| ---- | -------------- | --------------------------------------------- |
| 404  | 500090         | Feed collection not found                     |
| 500  | 500091, 500092 | Feed collection already exists, invalid input |

Error response shape: `{ "code": "<string>", "message": "<string>" }`. See `internal/errors` for feed collection codes (500090–500099).
