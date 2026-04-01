package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing-demo/domain"
	"testing-demo/mocks"
	"time"

	httpDelivery "testing-demo/delivery/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// Unit Test cho HTTP Handler (Delivery layer)
// Sử dụng httptest để giả lập HTTP request + MockOrderUsecase
// ============================================================

func setupRouter(mockUC *mocks.MockOrderUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	httpDelivery.NewOrderHandler(router, mockUC)
	return router
}

func TestPing(t *testing.T) {
	mockUC := new(mocks.MockOrderUsecase)
	router := setupRouter(mockUC)

	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "pong", resp["message"])
}

func TestCreateOrder_Handler_Success(t *testing.T) {
	mockUC := new(mocks.MockOrderUsecase)
	router := setupRouter(mockUC)

	expectedOrder := &domain.Order{
		ID:           1,
		CustomerName: "Nguyen Van A",
		Item:         "Laptop",
		Status:       "Pending",
		CreatedAt:    time.Now(),
	}
	mockUC.On("CreateOrder", "Nguyen Van A", "Laptop").Return(expectedOrder, nil)

	body := `{"customer_name": "Nguyen Van A", "item": "Laptop"}`
	req, _ := http.NewRequest("POST", "/orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, "Order created successfully", resp["message"])
	mockUC.AssertExpectations(t)
}

func TestCreateOrder_Handler_BadRequest(t *testing.T) {
	mockUC := new(mocks.MockOrderUsecase)
	router := setupRouter(mockUC)

	// missing "item" field
	body := `{"customer_name": "Nguyen Van A"}`
	req, _ := http.NewRequest("POST", "/orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUC.AssertNotCalled(t, "CreateOrder", mock.Anything, mock.Anything)
}

func TestCreateOrder_Handler_UsecaseError(t *testing.T) {
	mockUC := new(mocks.MockOrderUsecase)
	router := setupRouter(mockUC)

	mockUC.On("CreateOrder", "Nguyen Van A", "Laptop").
		Return(nil, errors.New("internal error"))

	body := `{"customer_name": "Nguyen Van A", "item": "Laptop"}`
	req, _ := http.NewRequest("POST", "/orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetOrder_Handler_Success(t *testing.T) {
	mockUC := new(mocks.MockOrderUsecase)
	router := setupRouter(mockUC)

	expectedOrder := &domain.Order{
		ID:           1,
		CustomerName: "Nguyen Van A",
		Item:         "Laptop",
		Status:       "Pending",
	}
	mockUC.On("GetOrder", 1).Return(expectedOrder, nil)

	req, _ := http.NewRequest("GET", "/orders/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetOrder_Handler_InvalidID(t *testing.T) {
	mockUC := new(mocks.MockOrderUsecase)
	router := setupRouter(mockUC)

	req, _ := http.NewRequest("GET", "/orders/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetOrder_Handler_NotFound(t *testing.T) {
	mockUC := new(mocks.MockOrderUsecase)
	router := setupRouter(mockUC)

	mockUC.On("GetOrder", 999).Return(nil, errors.New("order not found"))

	req, _ := http.NewRequest("GET", "/orders/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================================
// Table-Driven Unit Tests (Phiên bản chuyên nghiệp hơn)
// ============================================================

func TestCreateOrder_TableDriven(t *testing.T) {
	mockUC := new(mocks.MockOrderUsecase)
	router := setupRouter(mockUC)

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func()
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:        "Success",
			requestBody: `{"customer_name": "Nguyen Van A", "item": "Laptop"}`,
			mockSetup: func() {
				mockUC.On("CreateOrder", "Nguyen Van A", "Laptop").
					Return(&domain.Order{ID: 1, CustomerName: "Nguyen Van A", Item: "Laptop"}, nil).
					Once()
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Bad Request - Missing Item",
			requestBody:    `{"customer_name": "Nguyen Van A"}`,
			mockSetup:      func() {}, // Không gọi usecase
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Usecase Error - Internal Server Error",
			requestBody: `{"customer_name": "Nguyen Van A", "item": "Laptop"}`,
			mockSetup: func() {
				mockUC.On("CreateOrder", "Nguyen Van A", "Laptop").
					Return(nil, errors.New("database connection failed")).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req, _ := http.NewRequest("POST", "/orders", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestGetOrder_TableDriven(t *testing.T) {
	mockUC := new(mocks.MockOrderUsecase)
	router := setupRouter(mockUC)

	tests := []struct {
		name           string
		orderID        string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:    "Success",
			orderID: "1",
			mockSetup: func() {
				mockUC.On("GetOrder", 1).
					Return(&domain.Order{ID: 1, CustomerName: "A", Item: "B"}, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid ID format",
			orderID:        "abc",
			mockSetup:      func() {}, // Không gọi usecase
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Order Not Found",
			orderID: "999",
			mockSetup: func() {
				mockUC.On("GetOrder", 999).
					Return(nil, errors.New("order not found")).
					Once()
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req, _ := http.NewRequest("GET", "/orders/"+tt.orderID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUC.AssertExpectations(t)
		})
	}
}
