# Stage 1: Build the Go application
FROM golang:1.24-bullseye AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y gcc musl-dev libsqlite3-dev ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy dependency manifests and download
COPY go.mod go.sum ./
RUN go mod download

# Copy source tree
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Build the production web binary
RUN go build -ldflags="-s -w" -o /app/under100-web ./cmd/web

# Stage 2: Lightweight runtime image
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/under100-web /app/under100-web

# Copy templates, static files, and configuration
COPY templates/ /app/templates/
COPY static/ /app/static/
COPY config/ /app/config/

EXPOSE 8000
EXPOSE 10000

ENV PAGODA_HTTP_PORT=10000
ENV PAGODA_HTTP_HOSTNAME=0.0.0.0

ENTRYPOINT ["/app/under100-web"]