# Stage 1: Build the Go application
FROM golang:1.23.2 AS build

WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the Go source code and build the application
COPY . ./
RUN go build -o kasirku-service .

# Stage 2: Create the final production image
FROM ubuntu:latest

WORKDIR /app

# Copy the Go binary from the build stage
COPY --from=build /app/kasirku-service /app/

# Copy the .env file
COPY .env /app/

# Expose the port the app will run on
EXPOSE 8080

# Run the Go application
CMD ["./kasirku-service"]