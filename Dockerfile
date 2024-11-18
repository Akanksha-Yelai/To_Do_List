# Stage 1: Build
FROM golang:1.22 AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy the rest of your application code
COPY . .

# Build the Go application for Windows
RUN go build -o todo-app.exe .

# List the contents of the /app directory to confirm binary is created
RUN ls -l /app

# Stage 2: Final
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/todo-app.exe .

# Copy the .env file into the container
COPY .env /app/.env

# Copy the static files
COPY static/ ./static/

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["/app/todo-app.exe"]