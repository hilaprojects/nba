# Build stage
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum, and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Set the working directory to the 'cmd' directory
WORKDIR /app/cmd

# Build the Go binary for ARM64 statically (no cgo dependencies)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-extldflags=-static" -o main .

# Runtime stage (use a minimal base image for statically linked binary)
FROM scratch

# Copy the statically compiled binary from the builder stage
COPY --from=builder /app/cmd/main /main

# Set the entrypoint to the Go binary
CMD ["/main"]
