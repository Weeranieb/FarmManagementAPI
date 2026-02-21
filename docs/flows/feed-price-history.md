# Feed price history flow

## Purpose

Feed price history records price and date for a feed collection over time. Each record is tied to a `feedCollectionId`. This flow covers CRUD and listing for price history; used together with [feed-collection](feed-collection.md) to manage feed catalog and pricing.

## Actors / authorization

- All feed-price-history endpoints require JWT. Access is client-scoped (via feed collection’s client). Super admin can access any client’s data.

## Endpoints

| Method | Path                              | Description                                     |
| ------ | --------------------------------- | ----------------------------------------------- |
| POST   | `/api/v1/feed-price-history`      | Add a new feed price history record             |
| GET    | `/api/v1/feed-price-history`      | Get all feed price history (for current client) |
| GET    | `/api/v1/feed-price-history/{id}` | Get feed price history by ID                    |
| PUT    | `/api/v1/feed-price-history`      | Update a feed price history record              |

Full request/response schemas: [../openapi.yaml](../openapi.yaml).

## Request / response

- **Create**: Body `CreateFeedPriceHistoryRequest` — `feedCollectionId`, `price`, `priceUpdatedDate` (required). Response: success with created record.
- **List**: No body. Response: `data` as array of `FeedPriceHistoryResponse`.
- **Get by ID**: Path `id`. Response: `data` as `FeedPriceHistoryResponse`.
- **Update**: Body `UpdateFeedPriceHistoryRequest` — `id` (required); `feedCollectionId`, `price`, `priceUpdatedDate`. Response: success with updated record.

## Errors

| HTTP | Code (example) | Meaning                                          |
| ---- | -------------- | ------------------------------------------------ |
| 404  | 500100         | Feed price history not found                     |
| 500  | 500101, 500102 | Feed price history already exists, invalid input |

Error response shape: `{ "code": "<string>", "message": "<string>" }`. See `internal/errors` for feed price history codes (500100–500109).
