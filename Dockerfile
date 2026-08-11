# syntax=docker/dockerfile:1
#
# Production image for the Boonma Farm API (Go/Fiber).
# Built in CI and pushed to GHCR; the 1 GB droplet only ever pulls it.
# Multi-stage: fat toolchain in `build`, tiny static binary in `runtime`.

# ---- Build stage ------------------------------------------------------------
FROM golang:1.25-alpine AS build

WORKDIR /src

# Download modules first so this layer caches unless go.mod/go.sum change.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Now the source.
COPY . .

# CGO is off: the *server* uses the pure-Go pgx/gorm-postgres path (only the
# test suite needs CGO for SQLite). A static binary runs on a minimal image.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
      -trimpath -ldflags="-s -w" \
      -o /out/server ./src/cmd/api

# ---- Runtime stage ----------------------------------------------------------
FROM alpine:3.20 AS runtime

# ca-certificates: outbound TLS. tzdata: Asia/Bangkok timestamps.
RUN apk add --no-cache ca-certificates tzdata \
 && addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=build /out/server /app/server

# Never run the app as root.
USER app

EXPOSE 8080

ENTRYPOINT ["/app/server"]
