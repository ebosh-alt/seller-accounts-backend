# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder
WORKDIR /app

# Go deps
COPY go.mod go.sum ./
RUN go mod download

# Source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd/server

FROM alpine:3.20
WORKDIR /app

# Copy binary and assets
COPY --from=builder /app/bin/server /usr/local/bin/server
COPY config ./config

ENV APP_ENV=prod
EXPOSE 8080

CMD ["server"]
