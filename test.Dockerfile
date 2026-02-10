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

RUN make tools
RUN make check

RUN make test