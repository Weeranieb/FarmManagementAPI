# Builder stage
FROM golang:1.20.4-alpine as builder

WORKDIR /
# Make necessary directories
RUN mkdir -p api-server

WORKDIR /api-server
# Copy necessary files and directories
COPY . .

# Build go sources
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/server ./cmd/server

# Final stage
FROM ubuntu:focal

# Copy the API server binary and templates
COPY --from=builder /api-server/bin/server .

EXPOSE 4000
CMD sleep 5 && ./server