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

COPY --from=build-env /app/kasirku-service /app
COPY .env /app
# Copy binary dari tahap builder
COPY --from=builder /app/app .

# Jalankan aplikasi
CMD ["./app"]