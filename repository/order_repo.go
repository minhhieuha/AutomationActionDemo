package repository

import (
	"errors"
	"sync"
	"testing-demo/domain"
	"time"
)

type inMemoryOrderRepo struct {
	sync.RWMutex
	orders map[int]*domain.Order
	nextID int
}

// NewInMemoryOrderRepo creates a new in-memory order repository
func NewInMemoryOrderRepo() domain.OrderRepository {
	return &inMemoryOrderRepo{
		orders: make(map[int]*domain.Order),
		nextID: 1,
	}
}

func (r *inMemoryOrderRepo) Create(order *domain.Order) error {
	r.Lock()
	defer r.Unlock()

	order.ID = r.nextID
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}
	r.orders[order.ID] = order
	r.nextID++
	return nil
}

func (r *inMemoryOrderRepo) GetByID(id int) (*domain.Order, error) {
	r.RLock()
	defer r.RUnlock()

	if order, exists := r.orders[id]; exists {
		return order, nil
	}
	return nil, errors.New("order not found")
}

func (r *inMemoryOrderRepo) GetPendingOrders() ([]*domain.Order, error) {
	r.RLock()
	defer r.RUnlock()

	var pending []*domain.Order
	for _, order := range r.orders {
		if order.Status == "Pending" {
			pending = append(pending, order)
		}
	}
	return pending, nil
}

func (r *inMemoryOrderRepo) UpdateStatus(id int, status string) error {
	r.Lock()
	defer r.Unlock()

	if order, exists := r.orders[id]; exists {
		order.Status = status
		return nil
	}
	return errors.New("order not found")
}
