# Build stage
FROM golang:1.25-alpine AS builder

# Install necessary packages
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o bin/gogym-api \
    cmd/api/main.go

# Runtime stage
FROM alpine:3.19

# Install necessary runtime packages
RUN apk add --no-cache ca-certificates tzdata curl

# Create non-root user
RUN addgroup -g 1001 -S gogym && \
    adduser -S gogym -u 1001 -G gogym

# Set timezone to Asia/Tokyo
ENV TZ=Asia/Tokyo
RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/bin/gogym-api ./gogym-api
COPY --from=builder /app/infra ./infra

# Change ownership to non-root user
RUN chown -R gogym:gogym /app

# Switch to non-root user
USER gogym

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./gogym-api"]