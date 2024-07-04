# Use the official Go image as the base image
FROM golang:1.22-alpine3.19 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project files to the container
COPY . .

# Build the Go project
RUN go build -o app .

## Here create the new container that will server as the runner. Copy in the build output.
FROM alpine:3.19

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder container
COPY --from=builder /app/app .

# Create a user to run the binary
RUN adduser -D appuser
USER appuser

# Set the entry point to run the built binary
CMD ["./app"]
