# Stage 1: Build Stage
FROM golang:1.23.2 AS build-env

WORKDIR /app

# Copy the Go source code and .env file into the container
COPY . .

# Install dependencies and build the application
RUN go mod download
RUN go build -o kasirku-service .

# Stage 2: Production Stage
FROM ubuntu:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build-env /app/kasirku-service /app/

# Copy the .env file into the container (ensure the path is correct)
COPY .env /app/.env

# Expose the port your app will run on
EXPOSE 8080

# Run the application
CMD ["./kasirku-service"]