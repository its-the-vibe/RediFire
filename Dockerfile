# Build stage
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o redifire .

# Runtime stage
FROM scratch

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/redifire /app/redifire

# Copy CA certificates for HTTPS connections
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy example config (optional, users should mount their own)
COPY --from=builder /app/config.example.yaml /app/config.example.yaml

ENTRYPOINT ["/app/redifire"]
