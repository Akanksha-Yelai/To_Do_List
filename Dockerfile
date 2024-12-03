# Stage 1: Build
#FROM golang:1.22 AS builder

# Set the working directory
#WORKDIR /app

# Copy go mod and sum files
#COPY go.mod go.sum ./

# Download the Go modules dependencies
#RUN go mod tidy

# Copy the rest of your application code
#COPY . .

# for ARM systems
#ENV GOOS=linux GOARCH=amd64

# Build the Go application for linux
#RUN go build -o todo-app .

# Stage 2: Final
#FROM debian:bullseye-slim

# Set the working directory in the final image
#WORKDIR /app

# Copy the compiled binary from the builder stage
#COPY --from=builder /app/todo-app /app/todo-app

# Ensure it's executable
#RUN chmod +x /app/todo-app

# Install necessary dependencies
#RUN apt-get update && apt-get install -y \
    #libc6 \
    #libstdc++6 \
    #&& rm -rf /var/lib/apt/lists/*

# Copy the .env file into the container
#COPY .env /app/.env

# Copy the static files
#COPY static/ ./static/

# Expose port 8080 to the outside world
#EXPOSE 8080 5432

# Command to run the application
#ENTRYPOINT ["/app/todo-app"]

# Stage 1: Build
FROM golang:1.22 AS builder

# Set the working directory to /app
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Build the Go application statically linked for Linux
RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o todo-app ./main.go

# Step 1: Use Alpine Linux as the base image and add the mini root filesystem
FROM alpine:3.20.3

# Step 3: Set the working directory to /app
WORKDIR /app

# Copy the compiled binary from the builder stage into the container
COPY --from=builder /app/todo-app /app/todo-app

# Step 5: Copy the '.env' file to the /app directory
COPY .env /app/.env

# Step 6: Copy the static directory to serve static files (like HTML, CSS, JS)
COPY static /app/static

# Step 7: Expose port 8080 to allow communication
EXPOSE 8080

# Step 8: Set the entrypoint to run the 'todo-app' application
ENTRYPOINT ["/app/todo-app"]