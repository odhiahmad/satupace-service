# Gunakan base image Golang
FROM golang:1.23.2 as builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Copy source code
WORKDIR /app
COPY . .

# Build aplikasi
RUN go mod tidy
RUN go build -o app .

# Buat stage runtime dengan image lebih ringan
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=build-env /app/app /app/

# Copy the .env file into the container (ensure the path is correct)
COPY .env /app/.env

# Expose the port your app will run on
EXPOSE 8080

# Jalankan aplikasi
CMD ["./app"]