# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

# Install build dependencies for CGO (required by go-sqlite3)
RUN apk add --no-cache gcc musl-dev

WORKDIR /src

# Copy dependency files and pre-fetch modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the server binary with CGO enabled
# Use -ldflags to strip symbols and debug info for a smaller binary size
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/quazaar \
    ./cmd/server

# Stage 2: Runtime image
FROM alpine:3.19

# Install runtime dependencies (like ca-certificates for Spotify API calls)
RUN apk add --no-cache ca-certificates tzdata sqlite-libs

# Create a non-root user to run the application securely
RUN addgroup -S quazaar && adduser -S quazaar -G quazaar

USER quazaar
WORKDIR /home/quazaar

# Copy the compiled binary from the builder stage
COPY --from=builder --chown=quazaar:quazaar /app/quazaar /home/quazaar/quazaar

# Copy the .env.example as a default template if needed
COPY --chown=quazaar:quazaar .env.example /home/quazaar/.env.example

# Expose the default Quazaar port
EXPOSE 8765

# Create a mount point for SQLite database persistence
# The database will be created in /home/quazaar/.quazaar/quazaar.db
VOLUME ["/home/quazaar/.quazaar"]

# Start the server
CMD ["/home/quazaar/quazaar"]
