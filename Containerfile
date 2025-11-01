# Multi-stage build for minimal final image

# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary
# -ldflags="-s -w" to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o wingspan-goals .

# Stage 2: Create minimal runtime image
FROM registry.access.redhat.com/ubi10-minimal:10.0

# Set working directory
WORKDIR /app
RUN chown 1000:0 /app

# Set user to 1000
USER 1000

# Copy the binary from builder
COPY --from=builder /app/wingspan-goals .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./wingspan-goals"]
