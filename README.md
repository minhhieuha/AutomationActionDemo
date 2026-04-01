# Go Testing Demo: Order & Delivery System

Dự án này là một demo hoàn chỉnh về Clean Architecture trong Golang (sử dụng Gin làm API và Cronjob chạy ngầm), kết hợp với các khái niệm và kỹ thuật Software Testing cơ bản.

## 1. Các Khái Niệm Về Testing

### 1.1. Unit Test (Kiểm thử mức đơn vị)

- **Định nghĩa**: Là quá trình kiểm thử phần nhỏ nhất của phần mềm. Trong Golang, phần nhỏ nhất thường là một hàm (function) hoặc một phương thức (method).
- **Mục đích**: Đảm bảo từng thành phần nhỏ của code hoạt động đúng một cách độc lập trước khi ráp nối chúng lại với nhau. Thường sử dụng mock/stub để cô lập phần code được test với các thành phần khác như Database hoặc API bên ngoài.
- **Trong bài này**: Cụ thể là viết test cho `usecase` (logic nghiệp vụ tạo Order) và `delivery/http` (xử lý logic API độc lập với Service thực). (Xem file `order_usecase_test.go`).

### 1.2. Test Coverage (Độ phủ mã nguồn)

- **Định nghĩa**: Một thước đo cho biết bao nhiêu phần trăm (%) dòng code (statements/branches) của bạn đã được thực thi (chạy qua) khi chạy các bài test.
- **Mục đích**: Giúp developer nhận biết các đoạn code bị sót (chưa được test) từ đó tăng cường test case phòng tránh rủi ro tiềm ẩn. Coverage cao chưa chắc đã là hệ thống tốt (nếu assert sai), nhưng coverage thấp chắc chắn là ít được kiểm định.
- **Trong Golang**:
  - Chạy `go test ./... -cover` để xem tỷ lệ % test coverage.
  - Chạy `go test ./... -coverprofile=coverage.out` và `go tool cover -html=coverage.out` để xem vùng code nào màu xanh (đã test), màu đỏ (chưa test).

### 1.3. Smoke Test (Kiểm thử khối khói / Kiểm thử cơ bản)

- **Định nghĩa**: Loại bài test rất nhanh, nhẹ, thường chạy ngay sau khi build/deploy.
- **Mục đích**: Xác nhận xem các luồng quan trọng nhất (critical paths) hoặc "nhịp đập" (health) của chương trình có hoạt động không. Nếu ứng dụng "xì khói" (nghĩa là sập ngay lập tức), thì ngưng luôn, không cần chạy các test tốn thời gian khác.
- **Ví dụ Smoke Test bằng CUrl**:

  ```bash
  # Chạy server ở terminal 1: go run main.go
  # Tại terminal 2, chạy API cơ bản nhất (health check):
  curl -i http://localhost:8080/ping

  # Hoặc thử tạo 1 order:
  curl -X POST -H "Content-Type: application/json" -d '{"customer_name": "Nguyen Van A", "item": "Laptop"}' http://localhost:8080/orders
  ```

### 1.4. Stress Test (Kiểm thử khả năng chịu tải cục đọng)

- **Định nghĩa**: Kiểm thử xem hệ thống có thể chịu được lượng truy cập cao đến mức độ nào (hoặc quá tải) trước khi nó "sập" hoàn toàn hoặc phản hồi cực chậm.
- **Mục đích**: Tìm ra nghẽn cổ chai (bottleneck) của server (CPU, RAM, DB max connections), độ trễ (latency), và xem hệ thống recovery thế nào sau khi quá tải.
- **Demo Stress Test với k6 hoặc hey**: https://github.com/rakyll/hey
  Dùng công cụ `hey` (nếu có cài sẵn: `go install github.com/rakyll/hey@latest`):
  ```bash
  # Bắn 1000 request đồng thời (concurrent: 100), tổng cộng 5000 request vào API tạo order
  hey -n 5000 -c 100 -m POST -T "application/json" -d '{"customer_name": "Nguyen Van A", "item": "Laptop"}' http://localhost:8080/orders
  ```

## 2. Kiến Trúc Dự Án (Clean Architecture)

Dự án cấu trúc theo Clean Architecture để đảm bảo Code Testable (dễ dàng viết Unit test nhờ Interface/Mock):

- `domain`: Chứa các Model cốt lõi (`Order`) và Interface (các "khế ước" về Repository, Usecase).
- `repository`: Truy xuất dữ liệu. Ở đây dùng in-memory repo để giữ mọi thứ đơn giản và độc lập với Database thật.
- `usecase`: Code logic nghiệp vụ (ví dụ: tạo order -> set trạng thái là Pending).
- `delivery`: Nơi gọi vào Usecase. Ở đây có 2 loại delivery:
  - `http`: Giao tiếp của dự án bằng API REST qua framework `Gin`.
  - `cronjob`: Chạy nền một Scheduler. Cứ mỗi 10 giây hệ thống sẽ "điểm danh" order Pending và tự chuyển thành Delivered.

## 3. Cách Khởi Chạy

```bash
# 1. Tải các thư viện
go mod tidy

# 2. Cài đặt các file mock/thư viện bên thứ ba (đã cài chung về `go.mod`)
# 3. Chạy App
go run main.go
```

Khi app sống dậy, API chạy ở port `8080` và cronjob cứ 10s cất tiếng log một lần.
