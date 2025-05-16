FROM golang:1.24-alpine AS builder

# Build arguments
ARG VERSION=0.0.0-dev

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc g++ musl-dev bash coreutils git tar

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and make build script executable
COPY . /app/
RUN chmod +x /app/build.sh

RUN apk add --no-cache file && cd /app && ./build.sh -v ${VERSION}

# Use a smaller image for the final container
FROM alpine:latest

# Build arguments (needed in this stage too)
ARG VERSION=0.0.0-dev

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the binary from builder
COPY --from=builder /app/build/linux/${ARCH}/${VERSION#v}/usqlmcp /app/usqlmcp

# Create a volume for the database
VOLUME ["/data"]

# Set the entrypoint
ENTRYPOINT ["/app/usqlmcp"]

# Default command (can be overridden)
CMD ["--dsn", "sqlite3:///data/database.db"] 