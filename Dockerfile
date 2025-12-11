# Multi-stage build for Cloud Run deployment
# Frontend is hosted separately on Vercel, so we only build the API server
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (marketing/ and node_modules excluded via .dockerignore, web/ included for embed)
COPY . .

# Remove any node_modules that might have been copied (safety check)
RUN rm -rf web/node_modules marketing/node_modules || true

# Update go.sum to include all dependencies
RUN go mod tidy

# Build Go binary (API server only)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o barracuda-api .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/barracuda-api .

# Expose port (Cloud Run uses PORT env var, default to 8080)
ENV PORT=8080
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Run the API server
# Use shell form to expand PORT environment variable
CMD ./barracuda-api api --port ${PORT:-8080}

