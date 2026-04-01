package usecase

import (
	"errors"
	"testing-demo/domain"
)

type orderUsecase struct {
	orderRepo domain.OrderRepository
}

// NewOrderUsecase creates a new order usecase
func NewOrderUsecase(repo domain.OrderRepository) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo: repo,
	}
}

func (u *orderUsecase) CreateOrder(customerName, item string) (*domain.Order, error) {
	if customerName == "" || item == "" {
		return nil, errors.New("customer name and item cannot be empty")
	}

	order := &domain.Order{
		CustomerName: customerName,
		Item:         item,
		Status:       "Pending",
	}

	err := u.orderRepo.Create(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *orderUsecase) GetOrder(id int) (*domain.Order, error) {
	return u.orderRepo.GetByID(id)
}

// DeliverPendingOrders is business logic called by Cronjob
func (u *orderUsecase) DeliverPendingOrders() error {
	pendingOrders, err := u.orderRepo.GetPendingOrders()
	if err != nil {
		return err
	}

	for _, order := range pendingOrders {
		_ = u.orderRepo.UpdateStatus(order.ID, "Delivered")
	}

	return nil
}
