# Stage 1: Build the binary
FROM golang:1.24.0-alpine AS builder

WORKDIR /app

# copy the dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
# We use CGO_ENABLED=0 so the binary is "static" and runs anywhere
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app ./cmd

# Stage 2: Run the binary
FROM alpine:latest

WORKDIR /app

# Copy only the binary and your config folder from the builder
COPY --from=builder /app/app /app/app
COPY --from=builder /app/config /app/config

EXPOSE 8080

ENTRYPOINT ["/app/app"]