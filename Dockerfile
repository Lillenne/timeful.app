# Multi-stage Dockerfile for Timeful.app
# Builds both frontend and backend in a single image

# Stage 1: Build the frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files and install dependencies
COPY frontend/package*.json ./
RUN npm ci --prefer-offline --no-audit

# Copy frontend source and build
COPY frontend/ ./
ENV NODE_ENV=production
RUN npm run build

# Stage 2: Build the backend
FROM golang:1.20-alpine AS backend-builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app/server

# Copy Go module files and download dependencies
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy server source code
COPY server/ ./

# Build the Go server binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -o server .

# Stage 3: Create final runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk update && apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the server binary
COPY --from=backend-builder /app/server/server ./server

# Copy the built frontend
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Create directory for logs
RUN mkdir -p /app && chmod 755 /app

# Expose the port the server runs on
EXPOSE 3002

# Set environment variables
ENV GIN_MODE=release

# Run the server
CMD ["./server", "-release=true"]
