# Build stage
FROM golang:1.22.0 AS builder
WORKDIR /app

# Copy go.mod và go.sum trước để cache dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ source code
COPY . .

# Build từng service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/grpc-server ./cmd/grpc_server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/api-gateway ./cmd/api_gateway

# Runtime stage
FROM alpine:latest
WORKDIR /root/

# Copy binary từ builder stage
COPY --from=builder /app/grpc-server /root/grpc-server
COPY --from=builder /app/api-gateway /root/api-gateway
COPY --from=builder /app/templates /root/templates

# Cấp quyền thực thi
RUN chmod +x /root/grpc-server /root/api-gateway

# Chạy theo tham số command từ docker-compose
CMD ["./grpc-server"]
