# Build stage
FROM --platform=$BUILDPLATFORM golang:1.26.6-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o redifire .

# Runtime stage (distroless)
FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/redifire /app/redifire

# Copy example config (optional, users should mount their own)
COPY --from=builder /app/config.example.yaml /app/config.example.yaml

USER nonroot:nonroot

ENTRYPOINT ["/app/redifire"]
