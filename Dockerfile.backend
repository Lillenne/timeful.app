# Stage 1: Build the Go server
FROM golang:1.20-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy the server source code
COPY server/ ./

# Build the Go server binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -o server .

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk update && apk --no-cache add ca-certificates tzdata

# Create a non-root user to run the application
# Use a fixed UID/GID for better compatibility with volume mounts
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy the server binary from the builder stage
COPY --from=builder /app/server ./server

# Create directory for logs and set ownership
RUN mkdir -p /app/logs && \
    chown -R appuser:appuser /app && \
    chmod 755 /app/logs

# Switch to non-root user
USER appuser

# Expose the port the server runs on
EXPOSE 3002

# Set environment variables with defaults
ENV GIN_MODE=release

# Run the server
CMD ["./server", "-release=true"]
