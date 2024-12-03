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