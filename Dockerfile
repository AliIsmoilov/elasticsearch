# Stage 1: Build
FROM golang:1.24.0-alpine AS builder

WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy all project files
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

# Stage 2: Minimal runtime
FROM alpine:3.18

WORKDIR /app

# Copy the binary
COPY --from=builder /app/app .

# Copy the config folder (Casbin model + policy)
COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./app"]
