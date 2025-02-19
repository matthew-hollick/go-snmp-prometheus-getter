FROM golang:1.22-alpine

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/bin/snmp-prometheus-getter ./cmd/snmp-prometheus-getter

# Use a smaller base image for the final container
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=0 /app/bin/snmp-prometheus-getter /app/bin/snmp-prometheus-getter

# Create a non-root user
RUN adduser -D -u 10001 hedgehog_app && \
    chown -R hedgehog_app:hedgehog_app /app

# Switch to non-root user
USER hedgehog_app

# Run the application
CMD ["/app/bin/snmp-prometheus-getter", "-config", "/app/config.toml"]
