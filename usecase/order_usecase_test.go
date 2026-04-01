package usecase_test

import (
	"errors"
	"testing"
	"testing-demo/domain"
	"testing-demo/mocks"
	"testing-demo/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// Unit Test cho Order Usecase
// Mỗi test case cô lập hoàn toàn khỏi repository thật
// nhờ sử dụng MockOrderRepository (testify/mock)
// ============================================================

func TestCreateOrder_Success(t *testing.T) {
	// Arrange: tạo mock repo
	mockRepo := new(mocks.MockOrderRepository)
	uc := usecase.NewOrderUsecase(mockRepo)

	// Kỳ vọng: khi gọi Create với bất kỳ Order nào -> trả về nil (không lỗi)
	mockRepo.On("Create", mock.AnythingOfType("*domain.Order")).Return(nil)

	// Act: gọi usecase
	order, err := uc.CreateOrder("Nguyen Van A", "Laptop")

	// Assert: không lỗi, order được tạo đúng
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "Nguyen Van A", order.CustomerName)
	assert.Equal(t, "Laptop", order.Item)
	assert.Equal(t, "Pending", order.Status)

	// Verify: đảm bảo repo.Create đã được gọi đúng 1 lần
	mockRepo.AssertExpectations(t)
}

func TestCreateOrder_EmptyCustomerName(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	uc := usecase.NewOrderUsecase(mockRepo)

	// Act: customer_name rỗng
	order, err := uc.CreateOrder("", "Laptop")

	// Assert: phải trả về lỗi, order nil
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, "customer name and item cannot be empty", err.Error())

	// Verify: repo.Create KHÔNG được gọi (vì validate fail trước)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestCreateOrder_EmptyItem(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	uc := usecase.NewOrderUsecase(mockRepo)

	order, err := uc.CreateOrder("Nguyen Van A", "")

	assert.Error(t, err)
	assert.Nil(t, order)
}

func TestCreateOrder_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	uc := usecase.NewOrderUsecase(mockRepo)

	// Giả lập repo trả lỗi khi lưu
	mockRepo.On("Create", mock.AnythingOfType("*domain.Order")).
		Return(errors.New("database connection lost"))

	order, err := uc.CreateOrder("Nguyen Van A", "Laptop")

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, "database connection lost", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetOrder_Success(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	uc := usecase.NewOrderUsecase(mockRepo)

	expectedOrder := &domain.Order{
		ID:           1,
		CustomerName: "Nguyen Van A",
		Item:         "Laptop",
		Status:       "Pending",
	}
	mockRepo.On("GetByID", 1).Return(expectedOrder, nil)

	order, err := uc.GetOrder(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
	mockRepo.AssertExpectations(t)
}

func TestGetOrder_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	uc := usecase.NewOrderUsecase(mockRepo)

	mockRepo.On("GetByID", 999).Return(nil, errors.New("order not found"))

	order, err := uc.GetOrder(999)

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, "order not found", err.Error())
}

func TestDeliverPendingOrders_Success(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	uc := usecase.NewOrderUsecase(mockRepo)

	pendingOrders := []*domain.Order{
		{ID: 1, CustomerName: "A", Item: "Phone", Status: "Pending"},
		{ID: 2, CustomerName: "B", Item: "Tablet", Status: "Pending"},
	}

	mockRepo.On("GetPendingOrders").Return(pendingOrders, nil)
	mockRepo.On("UpdateStatus", 1, "Delivered").Return(nil)
	mockRepo.On("UpdateStatus", 2, "Delivered").Return(nil)

	err := uc.DeliverPendingOrders()

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeliverPendingOrders_NoPending(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	uc := usecase.NewOrderUsecase(mockRepo)

	mockRepo.On("GetPendingOrders").Return([]*domain.Order{}, nil)

	err := uc.DeliverPendingOrders()

	assert.NoError(t, err)
	// UpdateStatus should NOT be called when there are no pending orders
	mockRepo.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything)
}
