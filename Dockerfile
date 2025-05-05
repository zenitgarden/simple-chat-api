# Stage 1: Builder
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go binary
RUN go build -o main ./cmd

# Stage 2: Runtime
FROM alpine:3.18

# Install required tools
RUN apk add --no-cache curl tar ca-certificates

# Download and install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz \
  -o migrate.tar.gz && \
  tar -xzf migrate.tar.gz && \
  mv migrate /usr/bin/migrate && \
  chmod +x /usr/bin/migrate && \
  rm migrate.tar.gz

# Set working directory
WORKDIR /app

# Copy compiled binary and migrations
COPY --from=builder /app/main /app/
COPY --from=builder /app/migrations /app/migrations

# Expose application port
EXPOSE 8080

# Run migration and start the app
CMD migrate -path ./migrations -database "$DATABASE_URL" up && ./main
