# 🚀 Go CI/CD Blueprint: Smart Workflow & Cloud Deployment

Tài liệu này hướng dẫn chi tiết quy trình CI/CD "thông minh" và chuyên nghiệp được xây dựng cho dự án Go sử dụng GitHub Actions và Render.com.

---

## 🛠 1. Kiến trúc Quy trình (Architecture)

Quy trình được thiết kế theo nguyên tắc: **"Build một lần, Test mọi quy mô, Deploy an toàn."**

- **Ngôn ngữ:** Go 1.25.6 (Gin Framework)
- **Đóng gói:** Docker (Multi-stage build)
- **CI Tool:** GitHub Actions + `gotestsum` + `go-test-report`
- **Nền tảng Cloud:** Render.com (via Deploy Hooks)

---

## 🚦 2. Chiến lược Kiểm thử Thông minh (CI Pipeline)

Hệ thống CI tự động phân chia các bài kiểm tra dựa trên loại sự kiện để tối ưu thời gian:

| Sự kiện | Loại kiểm tra | Mục tiêu |
| :--- | :--- | :--- |
| **Push lên bất kỳ nhánh nào** | **Build & Unit Test** | Phản hồi tức thì cho Developer về logic cơ bản và lỗi cú pháp. |
| **Tạo Pull Request (PR)** | **Full Suite (Unit + Automation)** | Đảm bảo tính ổn định của toàn hệ thống trước khi gộp code. |
| **Push vào `main`/`develop`** | **Full Suite + CD** | Kiểm tra lần cuối và kích hoạt triển khai lên Cloud. |

> [!TIP]
> **Cơ chế Fallback Reporting:** Ngay cả khi code gặp lỗi **Panic**, hệ thống vẫn sẽ cố gắng tạo ra file `report.html` để bạn có thể xem chi tiết lỗi một cách trực quan nhất.

---

## 📦 3. Quy trình CD (Continuous Deployment)

Ứng dụng được đóng gói qua Docker giúp đảm bảo môi trường chạy nhất quán.

### Triển khai lên Render.com:
1. **Dockerfile:** Đóng gói ứng dụng thành file binary siêu nhẹ (Alpine Linux).
2. **Dynamic Port:** Ứng dụng tự động lắng nghe cổng `$PORT` do Cloud cấp.
3. **Deploy Hook:** GitHub Actions sẽ "bắn" tín hiệu qua API của Render để cập nhật bản mới nhất ngay sau khi tất cả các bài Test đạt kết quả **XANH (PASS)**.

---

## 📝 4. Hướng dẫn vận hành cho Developer

### Bước 1: Phát triển tính năng mới
Tạo một nhánh mới từ `develop`:
```bash
git checkout -b feature/ten-tinh-nang
```

### Bước 2: Kiểm tra khi đang code
Mỗi khi bạn `git push`, CI sẽ chạy **Unit Test**. Bạn có thể vào tab **Actions** trên GitHub để xem Dashboard HTML báo cáo kết quả.

### Bước 3: Gộp code vào hệ thống
Tạo một **Pull Request (PR)** từ nhánh của bạn vào `develop`. Lúc này bộ **Automation Test** sẽ chạy để kiểm tra các API Call và Database logic.

### Bước 4: Triển khai tự động
Khi PR được Merge vào nhánh chính, quy trình CD sẽ tự động đẩy ứng dụng lên Render.

---

## 🔑 5. Các cấu hình bắt buộc (Prerequisites)

Để quy trình chạy đúng, hãy đảm bảo bạn đã cấu hình các Secret sau trên GitHub:
- `RENDER_DEPLOY_HOOK_URL`: URL API từ trang Settings của Render Web Service.

---
*Tài liệu này được soạn thảo để phục vụ việc chuẩn hóa quy trình DevOps cho dự án.*
