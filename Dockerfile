# Stage 1: Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy only the necessary files for building
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o domain-redirect main.go

# Stage 2: Runtime stage
FROM scratch

# Set working directory
WORKDIR /app

# Copy only the built binary from builder stage
COPY --from=builder /app/domain-redirect .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./domain-redirect"]