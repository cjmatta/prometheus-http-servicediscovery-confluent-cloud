FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/prometheus-http-servicediscovery-confluent-cloud ./cmd/main.go

# Use a minimal alpine image for the final image
FROM alpine:3.17

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/prometheus-http-servicediscovery-confluent-cloud .

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/prometheus-http-servicediscovery-confluent-cloud"]