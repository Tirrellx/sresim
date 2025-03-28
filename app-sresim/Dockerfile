# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with additional security flags
RUN CGO_ENABLED=0 GOOS=linux \
    go build -a -installsuffix cgo \
    -ldflags="-w -s -extldflags '-static' \
    -X main.Version=$(git describe --tags --always --dirty) \
    -X main.BuildTime=$(date -u +'%Y-%m-%d_%H:%M:%S')" \
    -trimpath \
    -o sresim ./cmd/main.go

# Run stage
FROM gcr.io/distroless/static-debian12:nonroot

# Copy binary from builder
COPY --from=builder /app/sresim /sresim

# Copy config files if needed
COPY --from=builder /app/k8s/configmap.yaml /config.yaml

# Set security-related environment variables
ENV GIN_MODE=release \
    TZ=UTC \
    GOMAXPROCS=1 \
    GODEBUG=netdns=go \
    GOTRACEBACK=single

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/sresim", "health"]

# Run the application with security flags
ENTRYPOINT ["/sresim"]
