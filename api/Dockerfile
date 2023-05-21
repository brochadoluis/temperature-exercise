# Use the official Golang base image
FROM golang:1.17-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod and sum files
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the source code into the container
COPY ./api .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o api .

# Use a lightweight base image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app .

# Expose the desired port
EXPOSE 8080

# Define the command to run the Go application
CMD ["./api"]
