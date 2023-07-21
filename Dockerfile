# Use the base image with Go support
FROM golang:alpine AS builder

# Copy the source code
COPY . /app/

# Set the working directory
WORKDIR /app

# Build the application
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o chat-app ./cmd/main.go

# Create the final image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Copy the executable from the builder image
COPY --from=builder /app/chat-app /app/chat-app

EXPOSE 8080

# Set the entrypoint command
CMD ["/app/chat-app"]

