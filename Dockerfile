# Multi-stage Dockerfile for optimized image size
# Build stage
FROM golang:1.25.7-alpine AS builder

# Use a specific German mirror for Alpine packages
RUN echo "http://mirror1.hs-esslingen.de/pub/Mirrors/alpine/v3.22/main" > /etc/apk/repositories && \
    echo "http://mirror1.hs-esslingen.de/pub/Mirrors/alpine/v3.22/community" >> /etc/apk/repositories

# Install build dependencies only
RUN apk update && apk add --no-cache make git librdkafka-dev build-base

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=1 go build -ldflags="-w -s" -tags static,musl -o /app/main ./

# Runtime stage - minimal Alpine image
FROM alpine:3.22

# Use a specific German mirror for Alpine packages
RUN echo "http://mirror1.hs-esslingen.de/pub/Mirrors/alpine/v3.22/main" > /etc/apk/repositories && \
    echo "http://mirror1.hs-esslingen.de/pub/Mirrors/alpine/v3.22/community" >> /etc/apk/repositories

# Install only runtime dependencies
RUN apk update && apk add --no-cache ca-certificates curl && \
    rm -rf /var/cache/apk/*

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown appuser:appgroup /app/main

# Switch to non-root user
USER appuser

EXPOSE 8080

# Set the default command to run the application
ENTRYPOINT ["./main", "serve"]
