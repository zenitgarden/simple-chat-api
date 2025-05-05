# Stage 1: Builder
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build your Go app
RUN go build -o main ./cmd

# Stage 2: Runtime
FROM alpine:3.18

# Install curl and unzip to fetch migrate tool
RUN apk add --no-cache curl unzip

# Download golang-migrate CLI
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz \
  | tar xz && \
  mv migrate.linux-amd64 /usr/bin/migrate

# Create working directory
WORKDIR /app

# Copy compiled Go app and migrations
COPY --from=builder /app/main /app/
COPY --from=builder /app/migrations /app/migrations

# Expose app port
EXPOSE 8080

# Run migrations, then start the app
CMD migrate -path ./migrations -database "$DATABASE_URL" up && ./main
