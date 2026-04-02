# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Sao chép go.mod và go.sum trước để tận dụng Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Sao chép toàn bộ mã nguồn
COPY . .

# Biên dịch ứng dụng (Tắt CGO để có bản binary tĩnh, chạy tốt trên Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Runtime image siêu nhẹ
FROM alpine:latest

WORKDIR /root/

# Cài đặt curl để Render/Cloud có thể thực hiện health check nếu cần
RUN apk --no-cache add curl

# Sao chép bản binary từ bước builder
COPY --from=builder /app/main .

# Expose port (Render sẽ tự động dùng biến môi trường PORT)
EXPOSE 8080

# Chạy ứng dụng
CMD ["./main"]
