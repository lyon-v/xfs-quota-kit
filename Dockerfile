# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN make build

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/build/xfs-quota-kit .
COPY --from=builder /app/configs ./configs

# Create directories
RUN mkdir -p /var/log/xfs-quota-kit /var/backups/xfs-quota-kit

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ./xfs-quota-kit version || exit 1

# Run binary
CMD ["./xfs-quota-kit", "server", "--config", "configs/config.yaml"] 