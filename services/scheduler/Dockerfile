#
# Dockerfile for the Transfa Scheduler Microservice
#
# This Dockerfile uses a multi-stage build to create a minimal and secure container image.
# Stage 1: Builds the Go application into a static binary.
# Stage 2: Creates the final image by copying the binary into a scratch image.
#

# Stage 1: Build the application
# Use the official Golang Alpine image for a small and secure base for building.
FROM golang:1.21-alpine AS builder

# Set necessary environment variables for cross-compilation and static linking.
ENV CGO_ENABLED=0
ENV GOOS=linux

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies first.
# This leverages Docker's layer caching for faster subsequent builds.
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy the rest of the application source code into the container.
COPY . .

# Build the application, creating a static binary.
# -ldflags="-w -s" strips debugging information, reducing the final binary size.
RUN go build -ldflags="-w -s" -o /scheduler-service ./cmd/main.go

# Stage 2: Create a minimal final image from scratch
# 'scratch' is the most minimal image possible, containing only our binary.
FROM scratch

# Set the working directory in the final image.
WORKDIR /

# Copy the built binary from the builder stage into the final image.
COPY --from=builder /scheduler-service /scheduler-service

# Expose the port for health checks.
EXPOSE 8080

# Set the entrypoint for the container to run our application.
ENTRYPOINT ["/scheduler-service"]