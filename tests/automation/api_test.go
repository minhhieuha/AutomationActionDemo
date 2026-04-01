//go:build automation

package automation

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// TestOrderFlow thực hiện kịch bản test từ đầu đến cuối
func TestOrderFlow(t *testing.T) {
	// Khởi tạo httpexpect gắn với đối tượng t (testing.T)
	e := httpexpect.Default(t, "http://127.0.0.1:8080")

	// 1. Test Ping
	t.Log("🚀 Case 0: Ping check...")
	e.GET("/ping").
		Expect().
		Status(http.StatusOK).
		JSON().Object().Value("message").IsEqual("pong")

	// 2. Test Create Order
	t.Log("🚀 Case 1: Create Order...")
	resp := e.POST("/orders").
		WithJSON(map[string]interface{}{
			"customer_name": "Professional QA",
			"item":          "MacBook Air M3",
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	// Lấy ID từ Response
	orderID := resp.Value("data").Object().Value("id").Number().Raw()

	// 3. Test Get Order
	t.Logf("🚀 Case 2: Get Order ID %.0f...", orderID)
	e.GET("/orders/{id}").WithPath("id", orderID).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Path("$.data.customer_name").IsEqual("Professional QA")
}

// Bạn cũng có thể viết thêm các case test lỗi (Negative Cases)
func TestCreateOrder_InvalidData(t *testing.T) {
	e := httpexpect.Default(t, "http://127.0.0.1:8080")

	t.Log("🚀 Case 3: Create Order with missing data...")
	e.POST("/orders").
		WithJSON(map[string]interface{}{
			"customer_name": "Missing Item",
			// "item" bị thiếu
		}).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Value("message").String().Contains("Invalid request")
}
