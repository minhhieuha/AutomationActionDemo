package mocks

import (
	"testing-demo/domain"

	"github.com/stretchr/testify/mock"
)

// MockOrderUsecase is a mock implementation of domain.OrderUsecase
type MockOrderUsecase struct {
	mock.Mock
}

func (m *MockOrderUsecase) CreateOrder(customerName, item string) (*domain.Order, error) {
	args := m.Called(customerName, item)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderUsecase) GetOrder(id int) (*domain.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderUsecase) DeliverPendingOrders() error {
	args := m.Called()
	return args.Error(0)
}
