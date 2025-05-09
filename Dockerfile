FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments
ARG VERSION=0.0.0-dev
ARG ARCH=amd64

# Build the application
RUN ./build.sh -v ${VERSION} -a ${ARCH}

# Use a smaller image for the final container
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the binary from builder
COPY --from=builder /app/build/linux/${ARCH}/${VERSION}/usqlmcp /app/usqlmcp

# Create a volume for the database
VOLUME ["/data"]

# Set the entrypoint
ENTRYPOINT ["/app/usqlmcp"]

# Default command (can be overridden)
CMD ["--dsn", "sqlite3:///data/database.db"] 