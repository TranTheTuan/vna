# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency manifests and install
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Run Stage
FROM alpine:latest

WORKDIR /app

# Ensure logs directory exists
RUN mkdir -p /app/logs

# Copy binary from builder
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
