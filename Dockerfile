FROM golang:1.21-alpine AS builder

# Set the working directory to /app
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main main.go

# Use a smaller Alpine-based image for the final stage
FROM alpine:latest

# Copy the binary to the final image
COPY --from=builder /app/main /app/main

# Set working directory
WORKDIR /app

# Run the binary
CMD ["./main"]